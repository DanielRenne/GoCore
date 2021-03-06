package ginServer

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/DanielRenne/GoCore/core/extensions"
	"github.com/DanielRenne/GoCore/core/serverSettings"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

const TimeFormat = "Mon, 02 Jan 2006 15:04:05 GMT"
const StatusNotModified = 304 // RFC 7232, 4.1

var unixEpochTime = time.Unix(0, 0)

type ErrorResponse struct {
	Message string `json:"message"`
}

type LocaleLanguage struct {
	Locale   string
	Language string
}

func GetSessionKey(c *gin.Context, key string) (sessionKey string) {
	//https://github.com/gin-gonic/gin/issues/700
	// defer needed to catch session.Get concurrent map read write.
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()
	session := sessions.Default(c)
	if strings.Contains(c.Request.Host, ".com") {
		session.Options(sessions.Options{MaxAge: 86400 * serverSettings.WebConfig.Application.SessionExpirationDays,
			Secure: serverSettings.WebConfig.Application.SessionSecureCookie,
			Domain: ginCookieDomain})
	} else {
		session.Options(sessions.Options{Path: "/", MaxAge: 86400 * serverSettings.WebConfig.Application.SessionExpirationDays})
	}
	value := session.Get(key)
	if value == nil {
		sessionKey = ""
		return
	} else {
		sessionKey = session.Get(key).(string)
		return
	}
}

func SetSessionKey(c *gin.Context, key string, value string) {
	session := sessions.Default(c)
	if strings.Contains(c.Request.Host, ".com") {
		session.Options(sessions.Options{Path: "/", MaxAge: 86400 * serverSettings.WebConfig.Application.SessionExpirationDays,
			Secure: serverSettings.WebConfig.Application.SessionSecureCookie,
			Domain: ginCookieDomain})
	} else {
		session.Options(sessions.Options{Path: "/", MaxAge: 86400 * serverSettings.WebConfig.Application.SessionExpirationDays})
	}
	session.Set(key, value)
	session.Save()
}

func SaveSession(c *gin.Context) {
	session := sessions.Default(c)
	if strings.Contains(c.Request.Host, ".com") {
		session.Options(sessions.Options{MaxAge: 86400 * serverSettings.WebConfig.Application.SessionExpirationDays,
			Secure: serverSettings.WebConfig.Application.SessionSecureCookie,
			Domain: ginCookieDomain})
	} else {
		session.Options(sessions.Options{Path: "/", MaxAge: 86400 * serverSettings.WebConfig.Application.SessionExpirationDays})
	}
	session.Save()
}

func ClearSession(c *gin.Context) {
	session := sessions.Default(c)
	if strings.Contains(c.Request.Host, ".com") {
		session.Options(sessions.Options{MaxAge: 86400 * serverSettings.WebConfig.Application.SessionExpirationDays,
			Secure: serverSettings.WebConfig.Application.SessionSecureCookie,
			Domain: ginCookieDomain})
	} else {
		session.Options(sessions.Options{Path: "/", MaxAge: 86400 * serverSettings.WebConfig.Application.SessionExpirationDays})
	}
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

//Reads a file and responds with a base64 encoded string.  Primarily used for jquery ajax response binary data blob encoding.
func ReadFileBase64(path string, c *gin.Context) {
	page, err := ioutil.ReadFile(path)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	data := base64.StdEncoding.EncodeToString(page)
	c.Writer.Header().Set("Content-Length", extensions.IntToString(len(data)))
	c.String(http.StatusOK, data)
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

func RespondGzipJSFile(data []byte, modTime time.Time, c *gin.Context) {
	c.Header("Content-Type", "application/javascript")
	c.Header("Content-Encoding", "gzip")
	CheckLastModified(c.Writer, c.Request, modTime)
	c.Writer.Write(data)
}

func RespondJSFile(data []byte, modTime time.Time, c *gin.Context) {
	c.Header("Content-Type", "application/javascript")
	CheckLastModified(c.Writer, c.Request, modTime)
	c.Writer.Write(data)
}

func RespondTtfFile(data []byte, modTime time.Time, c *gin.Context) {
	c.Header("Content-Type", "application/x-font-ttf")
	CheckLastModified(c.Writer, c.Request, modTime)
	c.Writer.Write(data)
}

func RespondOtfFile(data []byte, modTime time.Time, c *gin.Context) {
	c.Header("Content-Type", "application/x-font-otf")
	CheckLastModified(c.Writer, c.Request, modTime)
	c.Writer.Write(data)
}

func RespondWoffFile(data []byte, modTime time.Time, c *gin.Context) {
	c.Header("Content-Type", "application/font-woff")
	CheckLastModified(c.Writer, c.Request, modTime)
	c.Writer.Write(data)
}

func RespondWoff2File(data []byte, modTime time.Time, c *gin.Context) {
	c.Header("Content-Type", "application/font-woff2")
	CheckLastModified(c.Writer, c.Request, modTime)
	c.Writer.Write(data)
}

func RespondEotFile(data []byte, modTime time.Time, c *gin.Context) {
	c.Header("Content-Type", "application/vnd.ms-fontobject")
	CheckLastModified(c.Writer, c.Request, modTime)
	c.Writer.Write(data)
}

func RespondSvgFile(data []byte, modTime time.Time, c *gin.Context) {
	c.Header("Content-Type", "image/svg+xml")
	CheckLastModified(c.Writer, c.Request, modTime)
	c.Writer.Write(data)
}

func ReadGzipCSSFile(path string, c *gin.Context) {

	c.Header("Content-Type", "text/css")
	c.Header("Content-Encoding", "gzip")
	c.File(path)

}

func RespondGzipCSSFile(data []byte, modTime time.Time, c *gin.Context) {
	c.Header("Content-Type", "text/css")
	c.Header("Content-Encoding", "gzip")
	CheckLastModified(c.Writer, c.Request, modTime)
	c.Writer.Write(data)
}

func ReadPngFile(path string, c *gin.Context) {
	c.Header("Content-Type", "image/png")
	c.File(path)
}

func ReadJpgFile(path string, c *gin.Context) {
	c.Header("Content-Type", "image/jpeg")
	c.File(path)
}

// modtime is the modification time of the resource to be served, or IsZero().
// return value is whether this request is now complete.
func CheckLastModified(w http.ResponseWriter, r *http.Request, modtime time.Time) bool {
	if modtime.IsZero() || modtime.Equal(unixEpochTime) {
		// If the file doesn't have a modtime (IsZero), or the modtime
		// is obviously garbage (Unix time == 0), then ignore modtimes
		// and don't process the If-Modified-Since header.
		return false
	}

	// The Date-Modified header truncates sub-second precision, so
	// use mtime < t+1s instead of mtime <= t to check for unmodified.
	if t, err := time.Parse(TimeFormat, r.Header.Get("If-Modified-Since")); err == nil && modtime.Before(t.Add(1*time.Second)) {
		h := w.Header()
		delete(h, "Content-Type")
		delete(h, "Content-Length")
		w.WriteHeader(StatusNotModified)
		return true
	}
	w.Header().Set("Last-Modified", modtime.UTC().Format(TimeFormat))
	return false
}
