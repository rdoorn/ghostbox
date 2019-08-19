package ghostbox

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rdoorn/gohelper/app"
	"github.com/rdoorn/gohelper/ginhelper"
	"github.com/rdoorn/gohelper/logging"
	"github.com/spf13/viper"
)

type Ghostbox struct {
	app.App

	// log is the default log handler
	logger logging.SimpleLogger

	// router is the gin handler
	router gin.Context

	// server is the http server handler
	server   http.Server
	shutdown chan struct{}
}

func New() *Ghostbox {
	return &Ghostbox{
		shutdown: make(chan struct{}),
	}
}

func (g *Ghostbox) Start(c *Config) error {
	if err := g.setLogging(viper.GetString("log_level"), viper.GetStringSlice("log_output")...); err != nil {
		return err
	}

	router := gin.New()
	router.Use(gin.Recovery(), ginhelper.Logger(g.logger))
	v1 := router.Group("/v1")
	{
		v1.POST("/version", g.apiV1Version)
		v1.GET("/hello/:name", g.apiV1Hello)
		v1.POST("/users", g.apiV1Signup) // new user
		// v1.PUT("/users", g.apiV1UpdateUser) // update existing user
	}

	// start https server
	g.server = http.Server{
		Addr:      c.Addr(),
		Handler:   router,
		TLSConfig: c.TLS.ParseValid(),
	}

	g.logger.Infof("Webservice started", "listener", c.Addr())
	if err := g.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	<-g.shutdown
	return nil
}

func (g *Ghostbox) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	g.server.SetKeepAlivesEnabled(false)
	if err := g.server.Shutdown(ctx); err != nil {
		g.logger.Fatalf("Could not gracefully shutdown the server: %v\n", err)
	}
	close(g.shutdown)
	g.logger.Infof("Webservice shutdown")
}

func (g *Ghostbox) setLogging(level string, output ...string) error {
	g.logger, _ = logging.NewZap(output...)
	logWrapper, err := g.logWrapper("core", level)
	if err != nil {
		return err
	}
	g.WithLogging(logWrapper)
	return nil
}

func (g *Ghostbox) logWrapper(description, level string) (logging.SimpleLogger, error) {
	logLevel, err := logging.ToLevel(level)
	if err != nil {
		return nil, fmt.Errorf("error setting log level: %s\n", err)
	}
	var prefix []interface{}
	prefix = append(prefix, "func")
	prefix = append(prefix, description)
	return &logging.Wrapper{
		Log:    g.logger,
		Level:  logLevel,
		Prefix: prefix,
	}, nil
}
