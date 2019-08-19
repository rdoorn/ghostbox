package ghostbox

import "github.com/gin-gonic/gin"

func (g *Ghostbox) apiV1Version(c *gin.Context) {
	g.Infof("hello verion request")
}
