package ginServer

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

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

func RespondJSON(v interface{}, c *gin.Context) {

	if v == nil {
		c.JSON(http.StatusNotFound, v)
		return
	}
	c.JSON(http.StatusOK, v)
}
