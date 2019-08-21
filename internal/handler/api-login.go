package handler

import (
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
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

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
	}

	jwtToken, err := h.NewJWTToken(user.ID, user.Username, user.Tokens)
	if err != nil {
		c.JSON(400, gin.H{"error": "failed to create a token"})
		return
	}

	c.JSON(200, gin.H{"error": "", "token": jwtToken})
}
