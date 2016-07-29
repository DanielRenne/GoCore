package ginServer

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

// Reads a file from the path parameter and returns to the client as text/html.
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

// Takes a string and returns to the client as text/html.
func RenderHTML(html string, c *gin.Context) {
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}

// Reads a file from the path parameter and returns to the client application/json
func ReadJSONFile(path string, c *gin.Context) {
	js, err := ioutil.ReadFile(path)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Header("Content-Type", "application/json")
	c.Writer.Write(js)
}

// Returns to the client application/json format for the passed interface.
func RespondJSON(v interface{}, c *gin.Context) {

	if v == nil {
		c.JSON(http.StatusNotFound, v)
		return
	}
	c.JSON(http.StatusOK, v)
}

//  Returns an byte array comprised of a JSON formated object with the error message.
func RespondError(message string) []byte {
	var msg ErrorResponse
	msg.Message = message
	b, _ := json.Marshal(msg)
	return b
}

func ReadGzipJSFile(path string, c *gin.Context) {

	c.Header("Content-Type", "application/javascript")
	c.Header("Content-Encoding", "gzip")
	c.File(path)

}
