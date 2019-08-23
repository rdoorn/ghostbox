package handler

import (
	"log"

	"github.com/gin-gonic/gin"
)

func (h *Handler) apiV1AccountActivate(c *gin.Context) {

	activationToken := c.Param("token")
	if activationToken == "" {
		c.JSON(400, gin.H{"error": "missing activation token"})
	}
	token, _ := c.Get("token")
	log.Printf("token: %+v credendials: %+v", token, token.(*JWTCredentials).Credentials)

	// update account to be active
	err := h.users.ActivateUser(token.(*JWTCredentials).Id, activationToken)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = h.users.SetStorageLimit(token.(*JWTCredentials).Id, FreeTeerStorageLimitMB)

	// update account with new token now that user is active
	err = h.users.AddToken(token.(*JWTCredentials).Id, "user")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// get new user settings
	user, err := h.users.GetByID(token.(*JWTCredentials).Id)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// re-issue token
	newToken, err := h.NewJWTTokenResponse(user.Id(), user.Username, user.Tokens)
	if err != nil {
		c.JSON(400, gin.H{"error": "failed to create a token"})
		return
	}

	c.JSON(200, newToken)
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
