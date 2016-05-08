package ginServer

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

var Router *gin.Engine

func init() {
	Router = gin.Default()
}

func ReadHTMLFile(path string, c *gin.Context) {
	page, err := ioutil.ReadFile(path)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	pageHTML := string(page)

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, pageHTML)
}
