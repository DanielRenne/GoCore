package ginServer

import (
	"sync"

	"github.com/DanielRenne/GoCore/core/serverSettings"
	"github.com/davidrenne/professor"
	"github.com/fatih/color"
	"github.com/gin-contrib/secure"

	// todo, this is legacy and deprecated, we need to move to something else to replace it
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

const (
	noWritten     = -1
	defaultStatus = 200
)

type routeGroup struct {
	enabled     bool
	routerGroup *gin.RouterGroup
}

// Router is the main gin.Engine that is used to serve all requests
var Router *gin.Engine

var groupRoutesSynced = struct {
	sync.RWMutex
	m map[string]routeGroup
}{m: make(map[string]routeGroup)}

type routerGroup struct {
	group  string
	route  string
	method string
	fp     func(*gin.Context)
}

var initializedRouterGroups []routerGroup
var hasInitialized bool
var ginCookieDomain string

// SessionConfiguration is used to configure the session cookie
type SessionConfiguration struct {
	Enabled               bool
	SessionKey            string
	SessionName           string
	SessionExpirationDays int
	SessionSecureCookie   bool
}

// ConfigureGin sets up all the GoCore features we currently offer
func ConfigureGin(mode string, cookieDomain string, secureHeaders bool, allowedHosts []string, csrfMiddleWareSecret string, cookieSession SessionConfiguration, launchPprof bool) {
	if launchPprof {
		// Run a safe pprof localhost server.
		professor.Launch("localhost:6897")
	}
	if mode == "release" || mode == "debug" {
		gin.SetMode(mode)
	} else {
		gin.SetMode("debug")
	}
	if cookieDomain != "" {
		ginCookieDomain = cookieDomain
	}
	Router = gin.Default()

	if secureHeaders {
		Router.Use(secure.New(secure.Config{
			AllowedHosts:          allowedHosts,
			SSLRedirect:           true,
			STSSeconds:            315360000,
			STSIncludeSubdomains:  true,
			FrameDeny:             true,
			ContentTypeNosniff:    true,
			BrowserXssFilter:      true,
			ContentSecurityPolicy: "default-src 'self'",
			IENoOpen:              true,
			ReferrerPolicy:        "strict-origin-when-cross-origin",
			SSLProxyHeaders:       map[string]string{"X-Forwarded-Proto": "https"},
		}))
	}

	if cookieSession.Enabled {
		store := sessions.NewCookieStore([]byte(cookieSession.SessionKey))
		store.Options(sessions.Options{MaxAge: 86400 * cookieSession.SessionExpirationDays,
			Secure: cookieSession.SessionSecureCookie})

		if cookieSession.SessionName != "" {
			Router.Use(sessions.Sessions(cookieSession.SessionName, store))
		} else {
			Router.Use(sessions.Sessions("defaultSession", store))
		}
	}

	if csrfMiddleWareSecret != "" {
		//Protect from CSRF Hacking
		Router.Use(csrf.Middleware(csrf.Options{
			Secret: csrfMiddleWareSecret,
			ErrorFunc: func(c *gin.Context) {
				c.String(400, "CSRF token mismatch")
				c.Abort()
			},
		}))
	}

	hasInitialized = true

	for _, group := range initializedRouterGroups {
		AddRouterGroup(group.group, group.route, group.method, group.fp)
	}
	initializedRouterGroups = nil
}

// Initialize is an internal export used in app/app.go based on your webConfig.json
func Initialize(mode string, cookieDomain string) {
	serverSettings.WebConfigMutex.RLock()
	csrf := serverSettings.WebConfig.Application.CSRFSecret
	cookieConfig := SessionConfiguration{
		Enabled:               true,
		SessionKey:            serverSettings.WebConfig.Application.SessionKey,
		SessionName:           serverSettings.WebConfig.Application.SessionName,
		SessionExpirationDays: serverSettings.WebConfig.Application.SessionExpirationDays,
		SessionSecureCookie:   serverSettings.WebConfig.Application.SessionSecureCookie,
	}
	serverSettings.WebConfigMutex.RUnlock()
	ConfigureGin(mode, cookieDomain, false, []string{}, csrf, cookieConfig, true)
}

// InitializeLite is an internal export used in app/app.go based on if you use app.InitializeLite()
func InitializeLite(mode string, secureHeaders bool, allowedHosts []string) {
	ConfigureGin(mode, "", secureHeaders, allowedHosts, "", SessionConfiguration{}, true)
}

// AddRouterGroup adds a gin-gonic router group to the gin server
func AddRouterGroup(group string, route string, method string, fp func(*gin.Context)) {

	if !hasInitialized {
		rgroup := routerGroup{group: group, route: route, method: method, fp: fp}
		initializedRouterGroups = append(initializedRouterGroups, rgroup)
		return
	}

	rg := getRouterSyncGroup(group)

	switch method {
	case "GET":
		rg.GET(route, fp)
		color.Green("Added GET to " + group + " with route " + route)
	case "POST":
		rg.POST(route, fp)
		color.Green("Added POST to " + group + " with route " + route)
	case "PUT":
		rg.PUT(route, fp)
		color.Green("Added PUT to " + group + " with route " + route)
	case "DELETE":
		rg.DELETE(route, fp)
		color.Green("Added DELETE to " + group + " with route " + route)
	case "PATCH":
		rg.PATCH(route, fp)
		color.Green("Added PATCH to " + group + " with route " + route)
	case "HEAD":
		rg.HEAD(route, fp)
		color.Green("Added HEAD to " + group + " with route " + route)
	case "OPTIONS":
		rg.OPTIONS(route, fp)
		color.Green("Added OPTIONS to " + group + " with route " + route)
	}
}

func getRouterSyncGroup(group string) *gin.RouterGroup {
	groupRoutesSynced.RLock()
	n := groupRoutesSynced.m[group]
	groupRoutesSynced.RUnlock()

	if n.enabled == false {
		r := Router.Group(group)
		rg := routeGroup{enabled: true, routerGroup: r}
		addRouterSyncGroup(group, rg)
		color.Green("Added Routing Group " + group)
		return rg.routerGroup
	}

	return n.routerGroup
}

func addRouterSyncGroup(group string, rg routeGroup) {
	groupRoutesSynced.Lock()
	groupRoutesSynced.m[group] = rg
	groupRoutesSynced.Unlock()
}
