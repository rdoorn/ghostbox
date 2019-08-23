package handler

import "github.com/gin-gonic/gin"

type ApiResponse struct {
	Error string
}

func setErrorRedirect(c *gin.Context, status int, err error, redirect string) error {
	c.Error(err)
	c.JSON(400, gin.H{"error": err.Error(), "redirect": redirect})
	return err
}

func setError(c *gin.Context, status int, err error) error {
	c.Error(err)
	c.JSON(400, gin.H{"error": err.Error()})
	return err
}

/*func setResponse(c *gin.Context, i interface{}) error {
	c.JSON(200, gin.H{"error": "", "data": i})
	return err
}*/
