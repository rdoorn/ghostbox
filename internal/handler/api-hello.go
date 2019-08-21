package handler

import "github.com/gin-gonic/gin"

func (h *Handler) apiV1Hello(c *gin.Context) {
	c.JSON(200, gin.H{"name": c.Param("name")})
	//m.Infof("hello %s", c.Param("name"))
}
