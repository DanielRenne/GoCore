package cookie

import (
	"strings"
	"time"

	// "github.com/DanielRenne/goCoreAppTemplate/web/app/javascript/globals"
	"github.com/gopherjs/gopherjs/js"
	"honnef.co/go/js/dom"
)

var (
	doc = dom.GetWindow().Document().(dom.HTMLDocument)
)

func uriDecode(str string) string {
	return js.Global.Call("decodeURIComponent", str).String()
}

func uriEncode(str string) string {
	return js.Global.Call("encodeURIComponent", str).String()
}

// Get returns a given cookie by name. If the cookie is not set, ok will be
// set to false
func Get(name string) (value string, ok bool) {
	cookieStr := doc.Cookie()
	if cookieStr == "" {
		return "", false
	}
	cookiePairs := strings.Split(cookieStr, "; ")
	for _, c := range cookiePairs {
		equalIndex := strings.IndexByte(c, '=')
		cookieName := c[:equalIndex]
		if cookieName == name {
			cookieValue := c[equalIndex+1:]
			return uriDecode(cookieValue), true
		}
	}
	return "", false
}

// SetString sets a cookie given a correctly formatted cookie string
// i.e "username=John Smith; expires=Thu, 18 Dec 2013 12:00:00 UTC; path=/"
func SetString(cookie string) {
	doc.SetCookie(cookie)
}

// Set adds a cookie to a user's browser with a name, value, expiry and path
// value, path and expires can be omitted
func Set(name string, value string, expires *time.Time, path string) {
	if name == "" {
		return
	}
	var expiry string
	if expires != nil {
		e := *expires
		e = e.UTC()
		t := e.Format("Mon, 02 Jan 2006 15:04:05 UTC")
		expiry = "expires=" + t + "; "
	}
	if path != "" {
		path = " path=" + path
	}
	c := name + "=" + uriEncode(value) + "; " + expiry + path
	doc.SetCookie(c)
}

// Delete removes a cookie specified by name
func Delete(name string) {
	c := name + "=; expires=Thu, 01 Jan 1970 00:00:01 UTC;"
	doc.SetCookie(c)
}
