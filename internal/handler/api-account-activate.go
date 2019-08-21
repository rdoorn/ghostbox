package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func (h *Handler) apiV1ActivateAccount(c *gin.Context) {

	activationToken := c.Param("token")
	if activationToken == "" {
		c.JSON(400, gin.H{"error": "missing activation token"})
	}
	token, _ := c.Get("token")

	h.users.ActivateUser(token.(*JWTCredentials).Id)

	c.JSON(200, gin.H{"error": "", "redirect": fmt.Sprintf("/%s", token.(*JWTCredentials).Username)})
	/*
		user, err := h.users.GetByEmail(loginRequest.Email)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		hasher, _ := passwordhash.New("sha256")
		passwordHash := hasher.Hash(h.salt, user.PasswordSalt, loginRequest.Password)

		if passwordHash != user.PasswordHash {
			c.JSON(400, gin.H{"error": "invalid password"})
			return
		}*/

	/*
		err = h.users.CreateUser(signupRequest.Firstname, signupRequest.Lastname, signupRequest.Email, passwordSalt, passwordHash)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
	*/
	//m.Infof("hello %s", c.Param("name"))
}
