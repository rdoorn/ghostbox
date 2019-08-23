package handler

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rdoorn/gohelper/jwthelper"
)

var (
	JWTExpireDuration time.Duration = 24 * time.Hour
)

type apiTokenResponseV1 struct {
	*ApiResponse
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"` // Bearer
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"` // not used yet
}

type JWTCredentials struct {
	*jwthelper.Credentials
	Username string
	Expire   time.Time
}

func (h *Handler) NewJWTTokenResponse(id, username string, tokens []string) (*apiTokenResponseV1, error) {
	h.Debugf("creating token for", "id", id, "type", fmt.Sprintf("%T", id))
	accessToken, err := h.NewJWTToken(id, username, tokens)
	if err != nil {
		return nil, err
	}
	response := &apiTokenResponseV1{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   int(JWTExpireDuration / time.Second),
	}
	return response, nil
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
			setErrorRedirect(c, 403, fmt.Errorf("invalid authorization header"), "/login")
			c.Abort()
			return
		}

		token := &JWTCredentials{}
		if err := jwthelper.Validate(a[1], token); err != nil {
			setErrorRedirect(c, 403, fmt.Errorf("invalid authentication token"), "/login")
			c.Abort()
			return
		}

		if token.Expire.Before(time.Now()) {
			setErrorRedirect(c, 403, fmt.Errorf("token expired"), "/login")
			c.Abort()
			return
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
			setError(c, 403, fmt.Errorf("not authorized"))
			c.Abort()
			return
		}

		c.Set("token", token)
		c.Next()
	}
}
