package handler

import (
	"github.com/gin-gonic/gin"
)

type getUserResponseSelfV1 struct {
	*ApiResponse
	ID        string `json:"id"`
	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
}

func (h *Handler) apiV1AccountGet(c *gin.Context) {

	token, _ := c.Get("token")
	/*if !ok {
		setError(c, 400, fmt.Errorf("token not set"))
		return
	}*/

	// we are getting our own information
	h.Debugf("creds", "usernameparam", c.Param("username"), "usernametoken", token.(*JWTCredentials).Username)
	if c.Param("username") == token.(*JWTCredentials).Username {
		user, err := h.users.GetByID(token.(*JWTCredentials).Id)
		if err != nil {
			setError(c, 500, err)
			return
		}

		response := getUserResponseSelfV1{
			ID:        user.Id(),
			Username:  user.Username,
			Firstname: user.Firstname,
			Lastname:  user.Lastname,
		}
		for _, email := range user.Emails {
			if email.Primary == true {
				response.Email = email.Email
			}
		}

		c.JSON(200, response)

		return
	}
	h.Infof("get for user %s tokens: %v", token.(*JWTCredentials).Username, token.(*JWTCredentials).Tokens)
}
