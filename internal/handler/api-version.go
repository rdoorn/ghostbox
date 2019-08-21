package handler

import "github.com/gin-gonic/gin"

func (h *Handler) apiV1Version(c *gin.Context) {
	h.Infof("hello verion request")
}
