package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rdoorn/gohelper/passwordhash"
)

type signupRequestV1 struct {
	Username  string `json:"username" binding:"required"`
	Firstname string `json:"firstname" binding:"required"`
	Lastname  string `json:"lastname" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

func (h *Handler) apiV1Signup(c *gin.Context) {
	//c.JSON(200, gin.H{"name": c.Param("name")})
	var signupRequest signupRequestV1
	if err := c.BindJSON(&signupRequest); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user, err := h.users.GetByUsername(signupRequest.Username)
	if user != nil {
		c.JSON(400, gin.H{"error": "username is already in use"})
		return
	}

	user, err = h.users.GetByEmail(signupRequest.Email)
	if user != nil {
		c.JSON(400, gin.H{"error": "email address is already in use"})
		return
	}

	hasher, _ := passwordhash.New("sha256")
	passwordSalt := hasher.Salt()
	passwordHash := hasher.Hash(h.salt, passwordSalt, signupRequest.Password)

	activationToken := uuid.New().String()

	err = h.users.CreateUser(signupRequest.Username, signupRequest.Firstname, signupRequest.Lastname, signupRequest.Email, passwordSalt, passwordHash, activationToken)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	h.Infof("new user", "username", signupRequest.Username, "activationToken", activationToken)

	c.JSON(200, gin.H{"name": signupRequest.Firstname})
	//m.Infof("hello %s", c.Param("name"))
}
