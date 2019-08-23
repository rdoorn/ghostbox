package handler

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/rdoorn/gohelper/passwordhash"
)

type loginRequestV1 struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) apiV1Login(c *gin.Context) {
	//c.JSON(200, gin.H{"name": c.Param("name")})
	var loginRequest loginRequestV1
	if err := c.BindJSON(&loginRequest); err != nil {
		setError(c, 400, err)
		return
	}

	user, err := h.users.GetByEmail(loginRequest.Email)
	if err != nil {
		setError(c, 400, err)
		return
	}

	log.Printf("user by email: %+v", user)

	hasher, _ := passwordhash.New("sha256")
	passwordHash := hasher.Hash(h.salt, user.PasswordSalt, loginRequest.Password)

	if passwordHash != user.PasswordHash {
		setError(c, 400, err)
		return
	}

	token, err := h.NewJWTTokenResponse(user.Id(), user.Username, user.Tokens)
	if err != nil {
		setError(c, 400, err)
		return
	}
	/*
		jwtToken, err := h.NewJWTToken(user.ID, user.Username, user.Tokens)
	*/

	c.JSON(200, token)
}
