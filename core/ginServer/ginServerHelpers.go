package ginServer

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type LocaleLanguage struct {
	Locale   string
	Language string
}

func GetSessionKey(c *gin.Context, key string) string {
	session := sessions.Default(c)
	value := session.Get(key)
	if value == nil {
		return ""
	} else {
		return session.Get(key).(string)
	}
}

func SetSessionKey(c *gin.Context, key string, value string) {
	session := sessions.Default(c)
	session.Set(key, value)
	session.Save()
}

func SaveSession(c *gin.Context) {
	session := sessions.Default(c)
	session.Save()
}

func ClearSession(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
}

func GetLocaleLanguage(c *gin.Context) (ll LocaleLanguage) {
	header := c.Request.Header.Get("Accept-Language")
	allLanguages := strings.Split(header, ";")

	locals := strings.Split(allLanguages[0], ",")
	localsSplit := strings.Split(locals[0], "-")

	if len(localsSplit) == 1 && len(locals) == 2 {
		localsSplit = strings.Split(locals[1], "-")
	}

	ll.Language = localsSplit[0]
	if len(localsSplit) == 2 {
		ll.Locale = localsSplit[1]
	}
	return
}

func GetRequestBody(c *gin.Context) ([]byte, error) {
	body := c.Request.Body
	return ioutil.ReadAll(body)
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

func ReadJSFile(path string, c *gin.Context) {
	page, err := ioutil.ReadFile(path)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	pageHTML := string(page)

	c.Header("Content-Type", "text/javascript")
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

func RespondGzipJSFile(data []byte, c *gin.Context) {
	c.Header("Content-Type", "application/javascript")
	c.Header("Content-Encoding", "gzip")
	c.Writer.Write(data)
}

func ReadGzipCSSFile(path string, c *gin.Context) {

	c.Header("Content-Type", "text/css")
	c.Header("Content-Encoding", "gzip")
	c.File(path)

}

func RespondGzipCSSFile(data []byte, c *gin.Context) {
	c.Header("Content-Type", "text/css")
	c.Header("Content-Encoding", "gzip")
	c.Writer.Write(data)
}

func ReadPngFile(path string, c *gin.Context) {

	c.Header("Content-Type", "image/png")
	c.File(path)

}
