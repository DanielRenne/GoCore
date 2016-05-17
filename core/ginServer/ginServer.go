package ginServer

import (
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"sync"
)

type routeGroup struct {
	enabled     bool
	routerGroup *gin.RouterGroup
}

var Router *gin.Engine

var groupRoutesSynced = struct {
	sync.RWMutex
	m map[string]routeGroup
}{m: make(map[string]routeGroup)}

func init() {
	Router = gin.Default()
}

func AddRouterGroup(group string, route string, method string, fp func(*gin.Context)) {

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
