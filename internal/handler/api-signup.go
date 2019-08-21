package handler

import "github.com/gin-gonic/gin"

type signupRequestV1 struct {
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

	user, err := h.users.GetByEmail(signupRequest.Email)
	if user != nil {
		c.JSON(400, gin.H{"error": "email address is already in use"})
		return
	}

	err = h.users.CreateUser(signupRequest.Firstname, signupRequest.Lastname, signupRequest.Email, signupRequest.Password)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"name": signupRequest.Firstname})
	//m.Infof("hello %s", c.Param("name"))
}
