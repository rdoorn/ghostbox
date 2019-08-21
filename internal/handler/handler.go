package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rdoorn/gohelper/app"
	"github.com/rdoorn/gohelper/ginhelper"
	"github.com/rdoorn/gohelper/logging"
	"github.com/rdoorn/ixxi/internal/models"
	"github.com/spf13/viper"
)

type Handler struct {
	app.App

	// log is the default log handler
	logger logging.SimpleLogger

	// router is the gin handler
	router gin.Context

	// server is the http server handler
	server   http.Server
	shutdown chan struct{}

	users models.UserInterface
	files models.FileInterface
	db    models.DBInterface

	salt string
}

func New() *Handler {
	return &Handler{
		shutdown: make(chan struct{}),
	}
}

func (h *Handler) Start(c *Config) error {
	if err := h.setLogging(viper.GetString("log_level"), viper.GetStringSlice("log_output")...); err != nil {
		return err
	}

	h.salt = c.PasswordSalt

	// user store
	usersInterface, err := c.Users.Setup()
	if err != nil {
		return err
	}
	h.users = usersInterface.(models.UserInterface)

	// file store
	fileInterface, err := c.File.Setup()
	if err != nil {
		return err
	}
	h.files = fileInterface.(models.FileInterface)

	// db store
	dbInterface, err := c.DB.Setup()
	if err != nil {
		return err
	}
	h.db = dbInterface.(models.DBInterface)

	router := gin.New()
	router.Use(gin.Recovery(), ginhelper.Logger(h.logger))
	v1 := router.Group("/v1")
	{
		v1.POST("/version", h.apiV1Version)
		v1.POST("/login", h.apiV1Login)  // login
		v1.POST("/users", h.apiV1Signup) // new user
		// v1.PUT("/users", h.apiV1UpdateUser) // update existing user
	}
	{
		hasAccount := v1.Group("/")
		hasAccount.Use(JWTAuthenticationRequired())
		hasAccount.GET("/users/:username/activate/:token", h.apiV1ActivateAccount) // login
	}
	{
		standardUser := v1.Group("/")
		standardUser.Use(JWTAuthenticationRequired("user"))
		standardUser.GET("/hello/:name", h.apiV1Hello)
	}

	// start https server
	h.server = http.Server{
		Addr:      c.Addr(),
		Handler:   router,
		TLSConfig: c.TLS.ParseValid(),
	}

	h.logger.Infof("Webservice started", "listener", c.Addr())
	if err := h.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	<-h.shutdown
	return nil
}

func (h *Handler) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	h.server.SetKeepAlivesEnabled(false)
	if err := h.server.Shutdown(ctx); err != nil {
		h.logger.Fatalf("Could not gracefully shutdown the server: %v\n", err)
	}
	close(h.shutdown)
	h.logger.Infof("Webservice shutdown")
}

func (h *Handler) setLogging(level string, output ...string) error {
	h.logger, _ = logging.NewZap(output...)
	logWrapper, err := h.logWrapper("core", level)
	if err != nil {
		return err
	}
	h.WithLogging(logWrapper)
	return nil
}

func (h *Handler) logWrapper(description, level string) (logging.SimpleLogger, error) {
	logLevel, err := logging.ToLevel(level)
	if err != nil {
		return nil, fmt.Errorf("error setting log level: %s\n", err)
	}
	var prefix []interface{}
	prefix = append(prefix, "func")
	prefix = append(prefix, description)
	return &logging.Wrapper{
		Log:    h.logger,
		Level:  logLevel,
		Prefix: prefix,
	}, nil
}
