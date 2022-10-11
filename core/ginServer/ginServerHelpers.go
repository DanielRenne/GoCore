// Package ginServer contains the gin server and ginServer helper functions
package ginServer

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/DanielRenne/GoCore/core/extensions"
	"github.com/DanielRenne/GoCore/core/serverSettings"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

// TImeFormat useful for formatting timestamps
const TimeFormat = "Mon, 02 Jan 2006 15:04:05 GMT"

// StatusNotModified returns a 304 HTTP status code.
const StatusNotModified = 304 // RFC 7232, 4.1
var unixEpochTime = time.Unix(0, 0)
var mux sync.RWMutex

// ErrorResonse is a struct that is used to return errors to the client.
type ErrorResponse struct {
	Message string `json:"message"`
}

// LocaleLanguage is a struct that is used to return the language and locale of the client.
type LocaleLanguage struct {
	Locale   string
	Language string
}

// GetSessionKey returns a thread-safe value of the session key.
func GetSessionKey(c *gin.Context, key string) (sessionKey string) {
	//https://github.com/gin-gonic/gin/issues/700
	// defer needed to catch session.Get concurrent map read write.
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()
	mux.RLock()
	defer mux.RUnlock()

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

// SetSessionKey sets a thread-safe value of the session key.
func SetSessionKey(c *gin.Context, key string, value string) {
	mux.Lock()
	defer mux.Unlock()

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

// SaveSession saves the session of the client.
func SaveSession(c *gin.Context) {
	mux.Lock()
	defer mux.Unlock()
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

// ClearSession clears the session of the client.
func ClearSession(c *gin.Context) {
	mux.Lock()
	defer mux.Unlock()
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

// GetLocaleLanguage returns the locale and language of the client.
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

// GetRequestBody returns the body of the request as a string.
func GetRequestBody(c *gin.Context) ([]byte, error) {
	body := c.Request.Body
	return ioutil.ReadAll(body)
}

// ReadHTMLFile reads a file from the path parameter and returns to the client as text/html.
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

// ReadJSFile reads a file from the path parameter and returns to the client as text/javascript.
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

// ReadFileBase64 reads a file and responds with a base64 encoded string.  Primarily used for jquery ajax response binary data blob encoding.
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

// RenderHTML takes a string and returns to the client as text/html.
func RenderHTML(html string, c *gin.Context) {
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}

// ReadJSONFile reads a file from the path parameter and returns to the client application/json
func ReadJSONFile(path string, c *gin.Context) {
	js, err := ioutil.ReadFile(path)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Header("Content-Type", "application/json")
	c.Writer.Write(js)
}

// RespondJSON returns to the client application/json format for the passed interface.
func RespondJSON(v interface{}, c *gin.Context) {

	if v == nil {
		c.JSON(http.StatusNotFound, v)
		return
	}
	c.JSON(http.StatusOK, v)
}

// RespondError returns an byte array comprised of a JSON formated object with the error message.
func RespondError(message string) []byte {
	var msg ErrorResponse
	msg.Message = message
	b, _ := json.Marshal(msg)
	return b
}

// ReadGzipJSFile reads a file from the path parameter and returns to the client as application/gzip.
func ReadGzipJSFile(path string, c *gin.Context) {

	c.Header("Content-Type", "application/javascript")
	c.Header("Content-Encoding", "gzip")
	c.File(path)
}

// RespondGzipJSFile returns a file to the client as application/gzip.
func RespondGzipJSFile(data []byte, modTime time.Time, c *gin.Context) {
	c.Header("Content-Type", "application/javascript")
	c.Header("Content-Encoding", "gzip")
	CheckLastModified(c.Writer, c.Request, modTime)
	c.Writer.Write(data)
}

// RespondJSFile returns a file to the client as application/javascript.
func RespondJSFile(data []byte, modTime time.Time, c *gin.Context) {
	c.Header("Content-Type", "application/javascript")
	CheckLastModified(c.Writer, c.Request, modTime)
	c.Writer.Write(data)
}

// ResponsdTtfFile returns a file to the client as application/x-font-ttf.
func RespondTtfFile(data []byte, modTime time.Time, c *gin.Context) {
	c.Header("Content-Type", "application/x-font-ttf")
	CheckLastModified(c.Writer, c.Request, modTime)
	c.Writer.Write(data)
}

// RespondOtfFile returns a file to the client as application/x-font-opentype.
func RespondOtfFile(data []byte, modTime time.Time, c *gin.Context) {
	c.Header("Content-Type", "application/x-font-otf")
	CheckLastModified(c.Writer, c.Request, modTime)
	c.Writer.Write(data)
}

// ResondWoffFile returns a file to the client as application/font-woff.
func RespondWoffFile(data []byte, modTime time.Time, c *gin.Context) {
	c.Header("Content-Type", "application/font-woff")
	CheckLastModified(c.Writer, c.Request, modTime)
	c.Writer.Write(data)
}

// ResondWoff2File returns a file to the client as application/font-woff2.
func RespondWoff2File(data []byte, modTime time.Time, c *gin.Context) {
	c.Header("Content-Type", "application/font-woff2")
	CheckLastModified(c.Writer, c.Request, modTime)
	c.Writer.Write(data)
}

// RespondEotFile returns a file to the client as application/vnd.ms-fontobject.
func RespondEotFile(data []byte, modTime time.Time, c *gin.Context) {
	c.Header("Content-Type", "application/vnd.ms-fontobject")
	CheckLastModified(c.Writer, c.Request, modTime)
	c.Writer.Write(data)
}

// ResondSvgFile returns a file to the client as image/svg+xml.
func RespondSvgFile(data []byte, modTime time.Time, c *gin.Context) {
	c.Header("Content-Type", "image/svg+xml")
	CheckLastModified(c.Writer, c.Request, modTime)
	c.Writer.Write(data)
}

// ReadGzipCSSFile reads a file from the path parameter and returns to the client as text/css.
func ReadGzipCSSFile(path string, c *gin.Context) {

	c.Header("Content-Type", "text/css")
	c.Header("Content-Encoding", "gzip")
	c.File(path)

}

// RespondGzipCSSFile returns a file to the client as text/css.
func RespondGzipCSSFile(data []byte, modTime time.Time, c *gin.Context) {
	c.Header("Content-Type", "text/css")
	c.Header("Content-Encoding", "gzip")
	CheckLastModified(c.Writer, c.Request, modTime)
	c.Writer.Write(data)
}

// RespondCSSFile returns a file to the client as text/css.
func RespondCSSFile(data []byte, modTime time.Time, c *gin.Context) {
	c.Header("Content-Type", "text/css")
	CheckLastModified(c.Writer, c.Request, modTime)
	c.Writer.Write(data)
}

// ReadPngFile reads a file from the path parameter and returns to the client as image/png.
func ReadPngFile(path string, c *gin.Context) {
	c.Header("Content-Type", "image/png")
	c.File(path)
}

// ReadJpgFile reads a file from the path parameter and returns to the client as image/jpeg.
func ReadJpgFile(path string, c *gin.Context) {
	c.Header("Content-Type", "image/jpeg")
	c.File(path)
}

// CheckLastModified modtime is the modification time of the resource to be served, or IsZero().
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
