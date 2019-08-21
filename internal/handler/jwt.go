package handler

import (
	"math/rand"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rdoorn/gohelper/jwthelper"
)

type JWTCredentials struct {
	*jwthelper.Credentials
	Username string
	Expire   time.Time
}

func (h *Handler) NewJWTToken(id, username string, tokens []string) (string, error) {
	cred := &JWTCredentials{
		Credentials: &jwthelper.Credentials{
			Id:     id,
			Tokens: tokens,
			Nonce:  rand.Float64(),
		},
		Username: username,
		Expire:   time.Now().Add(24 * time.Hour),
	}
	signedToken, err := jwthelper.Sign(cred)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

/*
func (h *Handler) GetJWTSession(c *gin.Context) (*JWTCredentials, error) {
	auth := c.GetHeader("Authorization")
	a := strings.Split(auth, " ")
	if a[0] != "BEARER" && len(a) < 2 {
		return nil, fmt.Errorf("invalid authorization header")
	}

	token := &JWTCredentials{}
	if err := jwthelper.Validate(a[1], token); err != nil {
		return nil, fmt.Errorf("failed to validate token")
	}

	return token, nil
}
*/

func JWTAuthenticationRequired(tokens ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		a := strings.Split(auth, " ")
		if a[0] != "BEARER" && len(a) < 2 {
			c.JSON(400, gin.H{"error": "invalid authorization header", "redirect": "/login"})
		}

		token := &JWTCredentials{}
		if err := jwthelper.Validate(a[1], token); err != nil {
			c.JSON(400, gin.H{"error": "invalid authentication token", "redirect": "/login"})
		}

		if token.Expire.Before(time.Now()) {
			c.JSON(400, gin.H{"error": "token expired", "redirect": "/login"})
		}

		tokensRequired := len(tokens)
		for _, at := range token.Tokens { // active tokens
			for _, rt := range tokens { // requested tokens
				if at == rt {
					tokensRequired--
				}
			}
		}

		if tokensRequired != 0 {
			c.JSON(400, gin.H{"error": "not authorized"})
		}

		c.Set("token", token)
		c.Next()
	}
}
