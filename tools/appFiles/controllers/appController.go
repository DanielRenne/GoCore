//  Package app provides html routing & rendering for a go core app
package controllers

import (
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/globalsign/mgo/bson"

	"github.com/DanielRenne/GoCore/core"
	"github.com/DanielRenne/GoCore/core/app"
	"github.com/DanielRenne/GoCore/core/extensions"
	"github.com/DanielRenne/GoCore/core/fileCache"
	"github.com/DanielRenne/GoCore/core/ginServer"
	"github.com/DanielRenne/GoCore/core/logger"
	"github.com/DanielRenne/goCoreAppTemplate/br"
	"github.com/DanielRenne/goCoreAppTemplate/constants"
	appErrors "github.com/DanielRenne/goCoreAppTemplate/errors"
	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
	"github.com/DanielRenne/goCoreAppTemplate/payloads"
	"github.com/DanielRenne/goCoreAppTemplate/scheduleEngine"
	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
	"github.com/DanielRenne/goCoreAppTemplate/settings"
	"github.com/DanielRenne/goCoreAppTemplate/viewModel"
	"github.com/gin-gonic/gin"
	"github.com/mattn/go-isatty"
	"net/http"
	"reflect"
	"xojoc.pw/useragent"
)

var versionShort string

var (
	green        = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white        = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow       = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	red          = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blue         = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magenta      = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyan         = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	reset        = string([]byte{27, 91, 48, 109})
	disableColor = false
)

// html struct is used to extract any element within a content htm file.
type html struct {
	Head         head      `json:"head" xml:"head"`
	Body         body      `json:"body" xml:"body"`
	HotReloadCss cssStruct `json:"hotReloadCss" xml:"hotReloadCss"`
	HotReloadJs  jsStruct  `json:"hotReloadJs" xml:"hotReloadJs"`
}

type jsStruct struct {
	Content string `xml:",innerxml"`
}

type cssStruct struct {
	Content string `xml:",innerxml"`
}

type head struct {
	Content string `xml:",innerxml"`
}

type body struct {
	Content string `xml:",innerxml"`
}

type SocketAPIRequest struct {
	CallbackId      int    `json:"callBackId"`
	Context         string `json:"context"`
	ProxyGinContext []byte `json:"proxyGinContext"`
	// Data            []byte              `json:"Data"`

	// ModTime         time.Time           `json:"ModTime"`
	ApiRequest payloads.ApiRequest `json:"data"`
}

type SocketAPIResponse struct {
	CallbackId  int                  `json:"callBackId"`
	ApiResponse payloads.ApiResponse `json:"data"`
}

func Initialize() {
	ginServer.Router.Use(Logger())
	versionShort = strings.Replace(settings.Version, ".", "", -1)
	loadRoutes()
}

// Logger instances a Logger middleware that will write the logs to gin.DefaultWriter.
// By default gin.DefaultWriter = os.Stdout.
func Logger() gin.HandlerFunc {
	return LoggerWithWriter(gin.DefaultWriter)
}

func colorForStatus(code int) string {
	switch {
	case code >= 200 && code < 300:
		return green
	case code >= 300 && code < 400:
		return white
	case code >= 400 && code < 500:
		return yellow
	default:
		return red
	}
}

func colorForMethod(method string) string {
	switch method {
	case "GET":
		return blue
	case "POST":
		return cyan
	case "PUT":
		return yellow
	case "DELETE":
		return red
	case "PATCH":
		return green
	case "HEAD":
		return magenta
	case "OPTIONS":
		return white
	default:
		return reset
	}
}

// LoggerWithWriter instance a Logger middleware with the specified writter buffer.
// Example: os.Stdout, a file opened in write mode, a socket...
func LoggerWithWriter(out io.Writer, notlogged ...string) gin.HandlerFunc {
	isTerm := true

	if w, ok := out.(*os.File); !ok ||
		(os.Getenv("TERM") == "dumb" || (!isatty.IsTerminal(w.Fd()) && !isatty.IsCygwinTerminal(w.Fd()))) {
		isTerm = false
	}

	var skip map[string]struct{}

	if length := len(notlogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range notlogged {
			skip[path] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		// Start timer
		start := scheduleEngine.GetLocalTime(time.Now())
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log only when path is not being skipped
		if _, ok := skip[path]; !ok {
			// Stop timer
			end := scheduleEngine.GetLocalTime(time.Now())
			latency := end.Sub(start)

			clientIP := c.ClientIP()
			method := c.Request.Method
			statusCode := c.Writer.Status()
			var statusColor, methodColor string
			if isTerm {
				statusColor = colorForStatus(statusCode)
				methodColor = colorForMethod(method)
			}
			comment := c.Errors.ByType(gin.ErrorTypePrivate).String()

			if raw != "" {
				path = path + "?" + raw
			}

			fmt.Fprintf(out, "[GOCORE] %v |%s %3d %s| %13v | %15s |%s %-7s %s %s\n%s",
				end.Format("2006/01/02 - 15:04:05"),
				statusColor, statusCode, reset,
				latency,
				clientIP,
				methodColor, method, reset,
				path,
				comment,
			)
		}
	}
}

// Loads all routes for the package.
func loadRoutes() {
	ginServer.Router.GET("/", loadApp)
	ginServer.Router.POST("/api", handleApi)

	ginServer.Router.GET("/dist/javascript/AppInit.js", handleInit)
	ginServer.Router.GET("/dist/javascript/json.js", handleJsonInit)
	if settings.AppSettings.DeveloperMode {
		ginServer.Router.GET("/dist/javascript/go-core-app.js.gz", handleGzip)
		ginServer.Router.GET("/dist/javascript/go-core-app.js.map", handleMap)
		ginServer.Router.GET("/dist/css/go-core-app.css.gz", handleCss)
		ginServer.Router.GET("/dist/css/go-core-app.css.map", handleCssMap)
	} else {
		ginServer.Router.GET("/dist/javascript/go-core-app-"+versionShort+".js.gz", handleGzip)
		ginServer.Router.GET("/dist/javascript/go-core-app-"+versionShort+".js.map", handleMap)
		ginServer.Router.GET("/dist/css/go-core-app-"+versionShort+".css.gz", handleCss)
		ginServer.Router.GET("/dist/css/go-core-app-"+versionShort+".css.map", handleCssMap)
	}

	ginServer.Router.GET("/dist/javascript/libphonenumber.js.gz", handleLibPhoneGzip)
	ginServer.Router.GET("/dist/javascript/gopherjs.js.gz", handleGopherGzip)
	ginServer.Router.GET("/dist/javascript/gopherjs.js.map", handleGopherMap)
	ginServer.Router.GET("/dist/css/flags@2x.png", handleFlag2X)
	ginServer.Router.GET("/dist/css/flags.png", handleFlag)

	ginServer.Router.GET("/dist/javascript/polyfills.js", handlePolyfills)
	ginServer.Router.GET("/dist/markup", handleMarkupMiddleWare)

	//remark

	ginServer.Router.GET("/dist/css/remark-core.css.gz", handleRemarkGzip)
	ginServer.Router.GET("/dist/css/remark-experimental.css.gz", handleRemarkGzip2)

	ginServer.Router.GET("/dist/javascript/jquery.min.js.gz.js", handleRemarkJsGzipMin)

	// this is busted between linux and mac somehow.  prod will just use the min
	ginServer.Router.GET("/dist/javascript/bootstrap.min.js.gz.js", handleRemarkJsGzipBootstrap)
	ginServer.Router.GET("/dist/javascript/bootstrap.min.js", handleRemarkJsBootstrap)
	ginServer.Router.GET("/dist/javascript/animsition.min.js.gz.js", handleRemarkJsGzipAnimsition)
	ginServer.Router.GET("/dist/javascript/jquery-asScroll.min.js.gz.js", handleRemarkJsGzipJqueryAsscroll)
	ginServer.Router.GET("/dist/javascript/jquery.mousewheel.min.js.gz.js", handleRemarkJsGzipMousewheel)
	ginServer.Router.GET("/dist/javascript/jquery.asScrollable.all.min.js.gz.js", handleRemarkJsGzipAsscrollable)
	ginServer.Router.GET("/dist/javascript/jquery-asHoverScroll.min.js.gz.js", handleRemarkJsGzipJqueryAshoverscroll)
	ginServer.Router.GET("/dist/javascript/waves.min.js.gz.js", handleRemarkJsGzipWaves)
	ginServer.Router.GET("/dist/javascript/switchery.min.js.gz.js", handleRemarkJsGzipSwitchery)
	ginServer.Router.GET("/dist/javascript/intro.min.js.gz.js", handleRemarkJsGzipIntro)
	ginServer.Router.GET("/dist/javascript/screenfull.min.js.gz.js", handleRemarkJsGzipScreenfull)
	ginServer.Router.GET("/dist/javascript/jquery-slidePanel.min.js.gz.js", handleRemarkJsGzipJquerySlidepanel)
	ginServer.Router.GET("/dist/javascript/menu.min.js.gz.js", handleRemarkJsGzipMenu)
	ginServer.Router.GET("/dist/javascript/menubar.min.js.gz.js", handleRemarkJsGzipMenubar)
	ginServer.Router.GET("/dist/javascript/sidebar.min.js.gz.js", handleRemarkJsGzipSidebar)
	ginServer.Router.GET("/dist/javascript/config-colors.min.js.gz.js", handleRemarkJsGzipConfigColors)
	ginServer.Router.GET("/dist/javascript/config-tour.min.js.gz.js", handleRemarkJsGzipConfigTour)
	ginServer.Router.GET("/dist/javascript/asscrollable.min.js.Component.gz.js", handleRemarkJsGzipAsscrollableComponent)
	ginServer.Router.GET("/dist/javascript/animsition.min.js.Component.gz.js", handleRemarkJsGzipAnimsitionComponent)
	ginServer.Router.GET("/dist/javascript/slidepanel.min.js.Component.gz.js", handleRemarkJsGzipSlidepanelComponent)
	ginServer.Router.GET("/dist/javascript/switchery.min.js.Component.gz.js", handleRemarkJsGzipSwitcheryComponent)
	ginServer.Router.GET("/dist/javascript/tabs.min.js.Component.gz.js", handleRemarkJsGzipTabsComponent)
	ginServer.Router.GET("/dist/javascript/material-design.min.css.gz.js", handleRemarkJsGzipMaterialDesign)
	ginServer.Router.GET("/dist/javascript/brand-icons.min.css.gz.js", handleRemarkJsGzipBrandIcons)
	ginServer.Router.GET("/dist/javascript/html5shiv.min.js.gz.js", handleRemarkJsGzipHtml5shiv)
	ginServer.Router.GET("/dist/javascript/media.match.min.js.gz.js", handleRemarkJsGzipMedia)
	ginServer.Router.GET("/dist/javascript/respond.min.js.gz.js", handleRemarkJsGzipRespond)
	ginServer.Router.GET("/dist/javascript/breakpoints.min.js.gz.js", handleRemarkJsGzipBreakpoints)
	ginServer.Router.GET("/dist/javascript/asscrollable.min.js.jsComponent.gz.js", handleRemarkJsGzipAsscrollableComponent)
	ginServer.Router.GET("/dist/javascript/animsition.min.js.jsComponent.gz.js", handleRemarkJsGzipAnimsitionComponent)
	ginServer.Router.GET("/dist/javascript/slidepanel.min.js.jsComponent.gz.js", handleRemarkJsGzipSlidepanelComponent)
	ginServer.Router.GET("/dist/javascript/tabs.min.js.jsComponent.gz.js", handleRemarkJsGzipTabsComponent)
	ginServer.Router.GET("/dist/javascript/modernizr.min.js", handleRemarkJsModernizr)

	// this is busted between linux and mac somehow.  prod will just use the min
	ginServer.Router.GET("/dist/javascript/modernizr.min.js.gz.js", handleRemarkJsGzipModernizr)
	ginServer.Router.GET("/dist/javascript/core.min.js.gz.js", handleRemarkJsGzipCore)
	ginServer.Router.GET("/dist/javascript/site.min.js.gz.js", handleRemarkJsGzipSite)
	ginServer.Router.GET("/dist/javascript/moment.min.js.gz.js", handleRemarkJsGzipMoment)
	ginServer.Router.GET("/dist/javascript/moment-timezone.js.gz.js", handleRemarkJsGzipMomentTimeZone)
	ginServer.Router.GET("/dist/javascript/Material-Design-Iconic-Font.eot", handleRemarkJsMaterialDesignIconicFontEot)
	ginServer.Router.GET("/dist/javascript/Material-Design-Iconic-Font.svg", handleRemarkJsMaterialDesignIconicFontSvg)
	ginServer.Router.GET("/dist/javascript/Material-Design-Iconic-Font.ttf", handleRemarkJsMaterialDesignIconicFontTtf)
	ginServer.Router.GET("/dist/javascript/Material-Design-Iconic-Font.woff", handleRemarkJsMaterialDesignIconicFontWoff)
	ginServer.Router.GET("/dist/javascript/Material-Design-Iconic-Font.woff2", handleRemarkJsMaterialDesignIconicFontWoff2)
	ginServer.Router.GET("/dist/javascript/brand-icons.svg", handleRemarkJsBrandIconsSvg)
	ginServer.Router.GET("/dist/javascript/brand-icons.ttf", handleRemarkJsBrandIconsTtf)
	ginServer.Router.GET("/dist/javascript/brand-icons.woff", handleRemarkJsBrandIconsWoff)
	ginServer.Router.GET("/dist/javascript/brand-icons.woff2", handleRemarkJsBrandIconsWoff2)
	ginServer.Router.GET("/fileObject/:Id", handleFileObject)
	app.RegisterWebSocketDataCallback(handleWebSocketData)

}

func Rollbar() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		//log.Println("Every Request")
		for range c.Errors {
			// rollbar.RequestError(rollbar.ERR, c.Request, errors.New(e.Err))
			//log.Println("rollbar" + e.Error())
		}

	}
}

func loadApp(c *gin.Context) {
	htmlContent := AppIndex(c)
	handleRespondHTML(c, []byte(htmlContent), time.Now())
}

// app handles routing to the / path.
func AppIndex(c *gin.Context) (htmlContent string) {

	defer func() {
		if r := recover(); r != nil {
			session_functions.Print("\n\nPanic Stack: " + string(debug.Stack()))
			session_functions.Log("Panic Recovered at appController.AppIndex", fmt.Sprintf("%+v", r))
			return
		}
	}()

	markupData, _, err := readProductionCachedFile(settings.WebUI + "/app/index.htm")

	if err != nil {
		session_functions.Log("loadApp", "Failed to Read "+settings.WebUI+"/app/index.htm:  "+err.Error())
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	var redirectPage = ""
	var htm html
	err = xml.Unmarshal(markupData, &htm)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	reloadPage := false

	var appState viewModel.AppViewModel
	appState.LoadDefaultState()
	ua := useragent.Parse(c.Request.Header.Get("User-Agent"))
	if ua != nil {
		appState.UserAgent = ua
	}
	appState.SideBarMenu = viewModel.GetSideBarViewModel(c)
	appState.HTTPPort = settings.ServerSettings.HttpPort
	appState.DisplayVersion = settings.Version
	appState.ProductName = settings.ProductName
	appState.DeveloperMode = settings.AppSettings.DeveloperMode
	appState.DeveloperLogReact = settings.AppSettings.DeveloperLogReact
	appState.DeveloperLogState = settings.AppSettings.DeveloperLogState
	appState.DeveloperLogTheseObjects = settings.AppSettings.DeveloperLogTheseObjects
	appState.DeveloperSuppressTheseObjects = settings.AppSettings.DeveloperSuppressTheseObjects
	appState.DeveloperSuppressThesePages = settings.AppSettings.DeveloperSuppressThesePages
	appState.DeveloperLogStateChangePerformance = settings.AppSettings.DeveloperLogStateChangePerformance
	appState.ShowDialogSubmitBug2 = true

	roles := make(map[string]bool, 0)
	roles["NONE"] = true
	appState.HasRole = roles

	if session_functions.GetSessionAuthToken(c) == constants.COOKIE_AUTHED {
		appState.LoggedIn = true
		account, err := session_functions.GetSessionAccount(c)

		if err != nil {
			appState.SnackBarOpen = true
			appState.SnackBarMessage = "Failed to get account from Session."
			appState.SnackBarType = SNACKBAR_TYPE_ERROR
			appState.DialogMessage = appState.SnackBarMessage + "\n\n" + err.Error()
			appState.DialogTranslationTitle = constants.VIEWMODEL_DIALOG_ERROR_TITLE
			appState.DialogOpen = true

		} else {
			appState.AccountId = account.Id.Hex()
			appState.AccountName = account.AccountName
			appState.AccountTypeShort = account.AccountTypeShort
			appState.IsSystemAccount = account.IsSystemAccount
			appState.AccountUsername = account.Email
			appState.Banner.AccountName = account.AccountName
		}

		user, err := session_functions.GetSessionUser(c)
		if err != nil {
			appState.SnackBarOpen = true
			appState.SnackBarMessage = "Failed to get user from Session."
			appState.SnackBarType = SNACKBAR_TYPE_ERROR
			appState.DialogMessage = appState.SnackBarMessage + "\n\n" + err.Error()
			appState.DialogTranslationTitle = constants.VIEWMODEL_DIALOG_ERROR_TITLE
			appState.DialogOpen = true

		} else {

			appState.UserInitials = user.First[:1] + user.Last[:1]
			appState.UserFirst = user.First
			appState.UserLast = user.Last
			appState.UserId = user.Id.Hex()
			appState.UserLanguage = user.Language
			appState.UserEnforcePasswordChange = user.EnforcePasswordChange
			appState.UserEmail = user.Email
			appState.UserPrimaryAccount = user.DefaultAccountId
			appState.UserPreferences = user.Preferences

			if account.Id.Hex() != user.DefaultAccountId {
				appState.Banner.Color = constants.BANNER_COLOR_OTHER
				appState.Banner.IsSecondaryAccount = true
			} else {
				appState.Banner.Color = constants.BANNER_COLOR_DEFAULT
				appState.Banner.IsSecondaryAccount = false
			}

			if user.Language == "" {
				localeLanguage := ginServer.GetLocaleLanguage(c)
				user.Language = model.GetDefaultLocale(localeLanguage.Locale)
				user.DateFormat = "mm/dd/yyyy"
				user.TimeZone = "US/Eastern"
				t, _ := session_functions.StartTransaction(c)
				user.SaveWithTran(t)
				t.Commit()
			}

			session_functions.StoreDataFormat(c, user.Language, user.TimeZone, user.DateFormat)

			//If the User Account Role is a Dedicated Device then redirect to the Room
			accRole, err := session_functions.GetSessionAccountRole(c)
			if err == nil {
				appState.AccountRoleId = accRole.RoleId
			}

			var roleFeatures []model.RoleFeature
			model.RoleFeatures.Query().Filter(model.Q(model.FIELD_ROLEFEATURE_ROLEID, accRole.RoleId)).Join("Feature").All(&roleFeatures)
			for _, roleFeature := range roleFeatures {
				if roleFeature.Joins.Feature != nil {
					roles[roleFeature.Joins.Feature.Key] = true
				}
			}
			var features []model.Feature
			model.Features.Query().All(&features)
			for _, feature := range features {
				_, ok := roles[feature.Key]
				if !ok {
					roles[feature.Key] = false
				}
			}
			appState.HasRole = roles
		}

	} else {
		ginServer.SetSessionKey(c, constants.COOKIE_AUTH_TOKEN, "")
	}

	renderStandardSideBar(session_functions.PassContext(c), &appState.SideBarMenu)

	stateData, err := json.Marshal(appState)
	if err != nil {
		session_functions.Log("loadApp", "Failed to Marshal appState:  "+err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	htm.Head.Content += "<script type=\"text/javascript\">window.appState = JSON.parse(window.atob(\"" + base64.StdEncoding.EncodeToString(stateData) + "\"));</script>"
	var contentData []byte
	if appState.LoggedIn {
		user, err := session_functions.GetSessionUser(c)

		var langFile string
		if user.Language == "en" {
			langFile = "US"
		} else {
			langFile = user.Language
		}
		contentData, _, err = readProductionCachedFile(settings.WebRoot + "/globalization/translations/app/" + user.Language + "/" + langFile + ".json")

		if err != nil {
			contentData, _, err = readProductionCachedFile(settings.WebRoot + "/globalization/translations/app/en/US.json")
		}
	} else {
		localeLanguage := ginServer.GetLocaleLanguage(c)
		contentData, _, err = readProductionCachedFile(settings.WebRoot + "/globalization/translations/app/" + localeLanguage.Language + "/" + localeLanguage.Language + ".json")

		if err != nil {
			contentData, _, err = readProductionCachedFile(settings.WebRoot + "/globalization/translations/app/en/US.json")
		}
	}

	htm.Head.Content += "<script type=\"text/javascript\">window.appContent = JSON.parse(window.atob(\"" + base64.StdEncoding.EncodeToString(contentData) + "\"));</script>"
	htm.Body.Content = strings.Replace(htm.Body.Content, "//RedirectPage", redirectPage, -1)

	if settings.AppSettings.DeveloperMode {
		var domain string
		if settings.ServerSettings.ServerFQDN != "127.0.0.1" && settings.ServerSettings.ServerFQDN != "localhost" && settings.ServerSettings.ServerFQDN != "0.0.0.0" && settings.ServerSettings.ServerFQDN != "" {
			domain = settings.ServerSettings.ServerFQDN
		} else {
			domain = settings.ServerSettings.Domain
		}

		htm.HotReloadJs.Content = "<script src=\"http://" + domain + ":3000/dist/javascript/go-core-app.js\" type=\"text/javascript\"></script>"
		htm.HotReloadCss.Content = "<link rel=\"stylesheet\" href=\"http://" + domain + ":3000/dist/css/go-core-app.css\"/>"
	} else {
		htm.HotReloadJs.Content = "<script src=\"/dist/javascript/go-core-app.js.gz\" type=\"text/javascript\"></script>"
		htm.HotReloadCss.Content = "<link rel=\"stylesheet\" href=\"/dist/css/go-core-app.css.gz\"/>"

	}

	if reloadPage {
		reloadScript := "<script type=\"text/javascript\">location.reload();</script>"
		htm.Body.Content = htm.Body.Content + "\n" + reloadScript
	}

	data, err := xml.Marshal(htm)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	htmlContent = strings.Replace(string(data[:]), "<html>", "<!DOCTYPE html>", 1)
	return
}

func readProductionCachedFile(path string) (data []byte, modTime time.Time, err error) {

	f, err := os.Open(path)
	d, err := f.Stat()
	if err != nil {
		modTime = bod(time.Now())
	} else {
		modTime = d.ModTime()
	}
	defer f.Close()

	data, err = fileCache.GetFile(path)
	return
}

func bod(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

func handleApi(c *gin.Context) {
	partialPageRequest(c)
}

func partialPageRequest(c *gin.Context) {
	start := time.Now()
	var err error
	var body []byte
	body, err = ginServer.GetRequestBody(c)

	path := c.Query("path")

	if err != nil {
		//Handle Error
		return
	}
	var apiRequest payloads.ApiRequest
	errMarshal := json.Unmarshal(body, &apiRequest)
	if errMarshal != nil {
		//Handle Error
		return
	}

	if settings.AppSettings.DeveloperMode {
		if settings.ServerSettings.ReleaseMode == "development" {
			core.TransactionLogMutex.Lock()
			core.TransactionLog = ""
			core.TransactionLogMutex.Unlock()
		}
	}
	defer func() {
		session_functions.Println(logger.TimeTrack(start, path+"#"+apiRequest.Action))
	}()
	responseHandler := clientResponse(c)

	callState(path, apiRequest.Action, string(apiRequest.State[:]), c, responseHandler)
}

func clientResponse(c *gin.Context) session_functions.ServerResponse {

	return func(redirect string, globalMessage string, globalMessageType string, trace error, transactionId string, v interface{}) {
		respondFinal(c, redirect, globalMessage, globalMessageType, trace, transactionId, v)
	}
}

func respondFinal(c *gin.Context, redirect string, globalMessage string, globalMessageType string, trace error, transactionId string, v interface{}) {

	if globalMessageType == PARAM_SNACKBAR_TYPE_DOWNLOAD_CONTENT {

		base64Data := []byte(base64.StdEncoding.EncodeToString(v.([]byte)))

		c.Writer.Header().Set("Content-Disposition", "attachment; filename=\""+globalMessage+"\"")
		c.Writer.Header().Set("Content-Type", "application/octet-stream")
		c.Writer.Header().Set("Content-Length", extensions.IntToString(len(base64Data)))
		c.Writer.Header().Set("Content-Transfer-Encoding", "binary")
		c.Writer.Header().Set("Expires", "0")
		c.Writer.Write(base64Data)

		return
	} else if globalMessageType == PARAM_SNACKBAR_TYPE_DOWNLOAD_FILE {
		c.Writer.Header().Set("Content-Disposition", "attachment; filename=\""+globalMessage+"\"")
		c.Writer.Header().Set("Content-Type", "application/octet-stream")
		c.Writer.Header().Set("Content-Transfer-Encoding", "binary")
		c.Writer.Header().Set("Expires", "0")
		ginServer.ReadFileBase64(v.(string), c)
		return
	}

	var apiResponse payloads.ApiResponse
	apiResponse.Redirect = redirect
	apiResponse.GlobalMessage = globalMessage
	apiResponse.GlobalMessageType = globalMessageType
	apiResponse.Transactionid = transactionId
	if trace != nil {
		apiResponse.Trace = appErrors.PrintStackTrace(trace, 15) // Convert to full trace.
	}

	stateData, _ := json.Marshal(v)
	apiResponse.State = string(stateData[:])
	if settings.ServerSettings.ReleaseMode == "development" {
		core.TransactionLogMutex.RLock()
		devLog := "\"DeveloperLog\": \"" + base64.StdEncoding.EncodeToString([]byte(core.TransactionLog)) + "\""
		core.TransactionLogMutex.RUnlock()
		if apiResponse.State == "{}" {
			apiResponse.State = "{" + devLog + "}"
		} else {
			apiResponse.State = "{" + devLog + ", " + apiResponse.State[1:]
		}
	}

	ginServer.RespondJSON(&apiResponse, c)
}

func callState(controller string, action string, state string, c *gin.Context, responseHandler session_functions.ServerResponse) {
	ctl := getController(controller)
	methodToCall := ctl.MethodByName(action)

	if methodToCall.String() == "<invalid Value>" {
		//Handle Error and Return
		responseHandler("", "Failed to Call "+action+" for "+controller+".", PARAM_SNACKBAR_TYPE_ERROR, nil, "", base64.StdEncoding.EncodeToString([]byte(state)))
		return
	}

	contextHandler := session_functions.PassContext(c)

	in := []reflect.Value{}
	in = append(in, reflect.ValueOf(contextHandler))
	in = append(in, reflect.ValueOf(state))
	in = append(in, reflect.ValueOf(responseHandler))

	methodToCall.Call(in)

}

// Remark

func handleRemarkJsGzipMin(c *gin.Context) {
	data, modTime, _ := readProductionCachedFile(settings.WebRoot + "/remark/material/global/vendor/jquery/jquery.min.js.gz.js")
	ginServer.RespondGzipJSFile(data, modTime, c)
}
func handleRemarkJsGzipBootstrap(c *gin.Context) {
	data, modTime, _ := readProductionCachedFile(settings.WebRoot + "/remark/material/global/vendor/bootstrap/bootstrap.min.js.gz.js")
	ginServer.RespondGzipJSFile(data, modTime, c)
}
func handleRemarkJsBootstrap(c *gin.Context) {
	data, modTime, _ := readProductionCachedFile(settings.WebRoot + "/remark/material/global/vendor/bootstrap/bootstrap.min.js")
	ginServer.RespondJSFile(data, modTime, c)
}
func handleRemarkJsGzipAnimsition(c *gin.Context) {
	data, modTime, _ := readProductionCachedFile(settings.WebRoot + "/remark/material/global/vendor/animsition/animsition.min.js.gz.js")
	ginServer.RespondGzipJSFile(data, modTime, c)
}
func handleRemarkJsGzipJqueryAsscroll(c *gin.Context) {
	data, modTime, _ := readProductionCachedFile(settings.WebRoot + "/remark/material/global/vendor/asscroll/jquery-asScroll.min.js.gz.js")
	ginServer.RespondGzipJSFile(data, modTime, c)
}
func handleRemarkJsGzipMousewheel(c *gin.Context) {
	data, modTime, _ := readProductionCachedFile(settings.WebRoot + "/remark/material/global/vendor/mousewheel/jquery.mousewheel.min.js.gz.js")
	ginServer.RespondGzipJSFile(data, modTime, c)
}
func handleRemarkJsGzipAsscrollable(c *gin.Context) {
	data, modTime, _ := readProductionCachedFile(settings.WebRoot + "/remark/material/global/vendor/asscrollable/jquery.asScrollable.all.min.js.gz.js")
	ginServer.RespondGzipJSFile(data, modTime, c)
}
func handleRemarkJsGzipJqueryAshoverscroll(c *gin.Context) {
	data, modTime, _ := readProductionCachedFile(settings.WebRoot + "/remark/material/global/vendor/ashoverscroll/jquery-asHoverScroll.min.js.gz.js")
	ginServer.RespondGzipJSFile(data, modTime, c)
}
func handleRemarkJsGzipWaves(c *gin.Context) {
	data, modTime, _ := readProductionCachedFile(settings.WebRoot + "/remark/material/global/vendor/waves/waves.min.js.gz.js")
	ginServer.RespondGzipJSFile(data, modTime, c)
}
func handleRemarkJsGzipSwitchery(c *gin.Context) {
	data, modTime, _ := readProductionCachedFile(settings.WebRoot + "/remark/material/global/vendor/switchery/switchery.min.js.gz.js")
	ginServer.RespondGzipJSFile(data, modTime, c)
}
func handleRemarkJsGzipIntro(c *gin.Context) {
	data, modTime, _ := readProductionCachedFile(settings.WebRoot + "/remark/material/global/vendor/intro-js/intro.min.js.gz.js")
	ginServer.RespondGzipJSFile(data, modTime, c)
}
func handleRemarkJsGzipScreenfull(c *gin.Context) {
	data, modTime, _ := readProductionCachedFile(settings.WebRoot + "/remark/material/global/vendor/screenfull/screenfull.min.js.gz.js")
	ginServer.RespondGzipJSFile(data, modTime, c)
}
func handleRemarkJsGzipJquerySlidepanel(c *gin.Context) {
	data, modTime, _ := readProductionCachedFile(settings.WebRoot + "/remark/material/global/vendor/slidepanel/jquery-slidePanel.min.js.gz.js")
	ginServer.RespondGzipJSFile(data, modTime, c)
}
func handleRemarkJsGzipMenu(c *gin.Context) {
	data, modTime, _ := readProductionCachedFile(settings.WebRoot + "/remark/material/center/assets/js/sections/menu.min.js.gz.js")
	ginServer.RespondGzipJSFile(data, modTime, c)
}
func handleRemarkJsGzipMenubar(c *gin.Context) {
	data, modTime, _ := readProductionCachedFile(settings.WebRoot + "/remark/material/center/assets/js/sections/menubar.min.js.gz.js")
	ginServer.RespondGzipJSFile(data, modTime, c)
}
func handleRemarkJsGzipSidebar(c *gin.Context) {
	data, modTime, _ := readProductionCachedFile(settings.WebRoot + "/remark/material/center/assets/js/sections/sidebar.min.js.gz.js")
	ginServer.RespondGzipJSFile(data, modTime, c)
}
func handleRemarkJsGzipConfigColors(c *gin.Context) {
	data, modTime, _ := readProductionCachedFile(settings.WebRoot + "/remark/material/global/js/configs/config-colors.min.js.gz.js")
	ginServer.RespondGzipJSFile(data, modTime, c)
}
func handleRemarkJsGzipConfigTour(c *gin.Context) {
	data, modTime, _ := readProductionCachedFile(settings.WebRoot + "/remark/material/center/assets/js/configs/config-tour.min.js.gz.js")
	ginServer.RespondGzipJSFile(data, modTime, c)
}
func handleRemarkJsGzipAsscrollableComponent(c *gin.Context) {
	data, modTime, _ := readProductionCachedFile(settings.WebRoot + "/remark/material/global/js/components/asscrollable.min.jsComponent.gz.js")
	ginServer.RespondGzipJSFile(data, modTime, c)
}
func handleRemarkJsGzipAnimsitionComponent(c *gin.Context) {
	data, modTime, _ := readProductionCachedFile(settings.WebRoot + "/remark/material/global/js/components/animsition.min.jsComponent.gz.js")
	ginServer.RespondGzipJSFile(data, modTime, c)
}
func handleRemarkJsGzipSlidepanelComponent(c *gin.Context) {
	data, modTime, _ := readProductionCachedFile(settings.WebRoot + "/remark/material/global/js/components/slidepanel.min.jsComponent.gz.js")
	ginServer.RespondGzipJSFile(data, modTime, c)
}
func handleRemarkJsGzipSwitcheryComponent(c *gin.Context) {
	data, modTime, _ := readProductionCachedFile(settings.WebRoot + "/remark/material/global/js/components/switchery.min.jsComponent.gz.js")
	ginServer.RespondGzipJSFile(data, modTime, c)
}
func handleRemarkJsGzipTabsComponent(c *gin.Context) {
	data, modTime, _ := readProductionCachedFile(settings.WebRoot + "/remark/material/global/js/components/tabs.min.jsComponent.gz.js")
	ginServer.RespondGzipJSFile(data, modTime, c)
}
func handleRemarkJsGzipMaterialDesign(c *gin.Context) {
	data, modTime, _ := readProductionCachedFile(settings.WebRoot + "/remark/material/global/fonts/material-design/material-design.min.css.gz.js")
	ginServer.RespondGzipCSSFile(data, modTime, c)
}
func handleRemarkJsGzipBrandIcons(c *gin.Context) {
	data, modTime, _ := readProductionCachedFile(settings.WebRoot + "/remark/material/global/fonts/brand-icons/brand-icons.min.css.gz.js")
	ginServer.RespondGzipCSSFile(data, modTime, c)
}
func handleRemarkJsGzipHtml5shiv(c *gin.Context) {
	data, modTime, _ := readProductionCachedFile(settings.WebRoot + "/remark/material/global/vendor/html5shiv/html5shiv.min.js.gz.js")
	ginServer.RespondGzipJSFile(data, modTime, c)
}
func handleRemarkJsGzipMedia(c *gin.Context) {
	data, modTime, _ := readProductionCachedFile(settings.WebRoot + "/remark/material/global/vendor/media-match/media.match.min.js.gz.js")
	ginServer.RespondGzipJSFile(data, modTime, c)
}
func handleRemarkJsGzipRespond(c *gin.Context) {
	data, modTime, _ := readProductionCachedFile(settings.WebRoot + "/remark/material/global/vendor/respond/respond.min.js.gz.js")
	ginServer.RespondGzipJSFile(data, modTime, c)
}
func handleRemarkJsGzipBreakpoints(c *gin.Context) {
	data, modTime, _ := readProductionCachedFile(settings.WebRoot + "/remark/material/global/vendor/breakpoints/breakpoints.min.js.gz.js")
	ginServer.RespondGzipJSFile(data, modTime, c)
}
func handleRemarkJsModernizr(c *gin.Context) {
	data, modTime, _ := readProductionCachedFile(settings.WebRoot + "/remark/material/global/vendor/modernizr/modernizr.min.js")
	ginServer.RespondJSFile(data, modTime, c)
}
func handleRemarkJsGzipModernizr(c *gin.Context) {
	data, modTime, _ := readProductionCachedFile(settings.WebRoot + "/remark/material/global/vendor/modernizr/modernizr.min.js.gz.js")
	ginServer.RespondGzipJSFile(data, modTime, c)
}
func handleRemarkJsGzipCore(c *gin.Context) {
	data, modTime, _ := readProductionCachedFile(settings.WebRoot + "/remark/material/global/js/core.min.js.gz.js")
	ginServer.RespondGzipJSFile(data, modTime, c)
}
func handleRemarkJsGzipSite(c *gin.Context) {
	data, modTime, _ := readProductionCachedFile(settings.WebRoot + "/remark/material/center/assets/js/site.min.js.gz.js")
	ginServer.RespondGzipJSFile(data, modTime, c)
}
func handleRemarkJsGzipMoment(c *gin.Context) {
	data, modTime, _ := readProductionCachedFile(settings.WebRoot + "/remark/material/global/vendor/moment/moment.min.js.gz.js")
	ginServer.RespondGzipJSFile(data, modTime, c)
}
func handleRemarkJsGzipMomentTimeZone(c *gin.Context) {
	data, modTime, _ := readProductionCachedFile(settings.WebRoot + "/remark/material/global/vendor/moment/moment-timezone.js.gz.js")
	ginServer.RespondGzipJSFile(data, modTime, c)
}

// fonts
func handleRemarkJsMaterialDesignIconicFontEot(c *gin.Context) {
	data, modTime, _ := readProductionCachedFile(settings.WebRoot + "/remark/material/global/fonts/material-design/Material-Design-Iconic-Font.eot")
	ginServer.RespondEotFile(data, modTime, c)
}

func handleRemarkJsMaterialDesignIconicFontSvg(c *gin.Context) {
	data, modTime, _ := readProductionCachedFile(settings.WebRoot + "/remark/material/global/fonts/material-design/Material-Design-Iconic-Font.svg")
	ginServer.RespondSvgFile(data, modTime, c)
}

func handleRemarkJsMaterialDesignIconicFontTtf(c *gin.Context) {
	data, modTime, _ := readProductionCachedFile(settings.WebRoot + "/remark/material/global/fonts/material-design/Material-Design-Iconic-Font.ttf")
	ginServer.RespondTtfFile(data, modTime, c)
}

func handleRemarkJsMaterialDesignIconicFontWoff(c *gin.Context) {
	data, modTime, _ := readProductionCachedFile(settings.WebRoot + "/remark/material/global/fonts/material-design/Material-Design-Iconic-Font.woff")
	ginServer.RespondWoffFile(data, modTime, c)
}

func handleRemarkJsMaterialDesignIconicFontWoff2(c *gin.Context) {
	data, modTime, _ := readProductionCachedFile(settings.WebRoot + "/remark/material/global/fonts/material-design/Material-Design-Iconic-Font.woff2")
	ginServer.RespondWoff2File(data, modTime, c)
}

func handleRemarkJsBrandIconsSvg(c *gin.Context) {
	data, modTime, _ := readProductionCachedFile(settings.WebRoot + "/remark/material/global/fonts/brand-icons/brand-icons.svg")
	ginServer.RespondSvgFile(data, modTime, c)
}

func handleRemarkJsBrandIconsTtf(c *gin.Context) {
	data, modTime, _ := readProductionCachedFile(settings.WebRoot + "/remark/material/global/fonts/brand-icons/brand-icons.ttf")
	ginServer.RespondTtfFile(data, modTime, c)
}

func handleRemarkJsBrandIconsWoff(c *gin.Context) {
	data, modTime, _ := readProductionCachedFile(settings.WebRoot + "/remark/material/global/fonts/brand-icons/brand-icons.woff")
	ginServer.RespondWoffFile(data, modTime, c)
}

func handleRemarkJsBrandIconsWoff2(c *gin.Context) {
	data, modTime, _ := readProductionCachedFile(settings.WebRoot + "/remark/material/global/fonts/brand-icons/brand-icons.woff2")
	ginServer.RespondWoff2File(data, modTime, c)
}

func handleFileObject(c *gin.Context) {
	id := c.Param("Id")
	acct, err := session_functions.GetSessionAccount(c)
	if err != nil {
		ginServer.ReadJpgFile(settings.WebRoot+"/images/no-image-found.jpg", c)
		return
	}

	if !bson.IsObjectIdHex(id) {
		ginServer.ReadJpgFile(settings.WebRoot+"/images/no-image-found.jpg", c)
		return
	}

	var fileObj model.FileObject
	filter := make(map[string]interface{}, 2)
	filter[model.FIELD_FILEOBJECT_ID] = id
	filter[model.FIELD_FILEOBJECT_ACCOUNTID] = acct.Id.Hex()
	err = model.FileObjects.Query().Filter(filter).One(&fileObj)
	if err == nil {
		data, err := base64.StdEncoding.DecodeString(fileObj.Content)
		if err == nil {

			if fileObj.Path != "" {
				data, err := extensions.ReadFile(fileObj.Path)
				if err == nil {
					c.Writer.Header().Set("Content-Type", fileObj.Type)
					c.Writer.Header().Set("Content-Length", extensions.IntToString(len(data)))
					c.Writer.Write(data)
					if fileObj.SingleDownload {
						err = fileObj.Delete()
						if err != nil {
							session_functions.Log("Error->appController->handleFileObject", err.Error())
						}
					}
					return
				}
				if err != nil {
					session_functions.Log("Error->appController->handleFileObject", err.Error())
				}

				return
			}
			c.Writer.Header().Set("Content-Type", fileObj.Type)
			widthString := c.Request.URL.Query().Get("width")
			heightString := c.Request.URL.Query().Get("height")
			width := extensions.StringToInt(widthString)
			height := extensions.StringToInt(heightString)
			var resizeContent []byte
			var errResize error
			if width > 0 && height > 0 {
				resizeContent, errResize = br.FileObjects.Resize(data, fileObj.Name, uint(width), uint(height))
				if errResize == nil {
					c.Writer.Header().Set("Content-Length", extensions.IntToString(len(resizeContent)))
					c.Writer.Header().Set("Cache-Control", "no-cache")
				} else {
					ginServer.CheckLastModified(c.Writer, c.Request, fileObj.Modified)
					c.Writer.Header().Set("Cache-Control", "max-age=31536000")
					c.Writer.Header().Set("Content-Length", extensions.IntToString(fileObj.Size))
				}
			} else {
				ginServer.CheckLastModified(c.Writer, c.Request, fileObj.Modified)
				c.Writer.Header().Set("Cache-Control", "max-age=31536000")
				c.Writer.Header().Set("Content-Length", extensions.IntToString(fileObj.Size))
			}
			session_functions.Dump("resizeErr", errResize)

			if width > 0 && height > 0 && errResize == nil {
				c.Writer.Write(resizeContent)
			} else {
				c.Writer.Write(data)
			}
			return
		} else {
			ginServer.ReadJpgFile(settings.WebRoot+"/images/no-image-found.jpg", c)
			session_functions.Dump(err.Error())
		}
	}

	ginServer.ReadJpgFile(settings.WebRoot+"/images/no-image-found.jpg", c)
}

func handleGopherGzip(c *gin.Context) {
	ginServer.ReadGzipJSFile(settings.WebRoot+"/dist/javascript/gopherjs.js.gz", c)
}

func handleGopherMap(c *gin.Context) {
	ginServer.ReadHTMLFile(settings.WebRoot+"/dist/javascript/gopherjs.js.map", c)
}

func handleInit(c *gin.Context) {
	ginServer.ReadJSFile(settings.WebRoot+"/dist/javascript/AppInit.js", c)
}

func handleJsonInit(c *gin.Context) {
	ginServer.ReadJSFile(settings.WebRoot+"/dist/javascript/json.js", c)
}

func handleGzip(c *gin.Context) {
	data, modTime, _ := readProductionCachedFile(settings.WebRoot + "/dist/javascript/go-core-app.js.gz")
	ginServer.RespondGzipJSFile(data, modTime, c)
}

func handleRemarkGzip(c *gin.Context) {
	data, modTime, _ := readProductionCachedFile(settings.WebRoot + "/dist/css/remark-core.css.gz")
	ginServer.RespondGzipCSSFile(data, modTime, c)
}

func handleRemarkGzip2(c *gin.Context) {
	data, modTime, _ := readProductionCachedFile(settings.WebRoot + "/dist/css/remark-experimental.css.gz")
	ginServer.RespondGzipCSSFile(data, modTime, c)
}

func handleLibPhoneGzip(c *gin.Context) {
	data, modTime, _ := readProductionCachedFile(settings.WebRoot + "/dist/javascript/libphonenumber.js.gz")
	ginServer.RespondGzipJSFile(data, modTime, c)
}

func handleMap(c *gin.Context) {
	ginServer.ReadHTMLFile(settings.WebRoot+"/dist/javascript/go-core-app.js.map", c)
}

func handleCss(c *gin.Context) {
	data, modTime, _ := readProductionCachedFile(settings.WebRoot + "/dist/css/go-core-app.css.gz")
	ginServer.RespondGzipCSSFile(data, modTime, c)
}

func handleCssMap(c *gin.Context) {
	ginServer.ReadHTMLFile(settings.WebRoot+"/dist/css/go-core-app.css.map", c)
}

func handlePolyfills(c *gin.Context) {
	ginServer.ReadJSFile(settings.WebRoot+"/dist/javascript/polyfills.js", c)
}

func handleFlag(c *gin.Context) {
	ginServer.ReadPngFile(settings.WebRoot+"/dist/css/flags.png", c)
}

func handleFlag2X(c *gin.Context) {
	ginServer.ReadPngFile(settings.WebRoot+"/dist/css/flags@2x.png", c)
}

func handleMarkupMiddleWare(c *gin.Context) {
	start := time.Now()
	path := c.Query("path")
	file_name := c.Query("file")
	action := c.Query("action")
	uriParams := c.Query("uriParams")

	if settings.ServerSettings.ReleaseMode == "development" {
		core.TransactionLogMutex.Lock()
		core.TransactionLog = ""
		core.TransactionLogMutex.Unlock()
	}

	defer func() {
		session_functions.Println(logger.TimeTrack(start, "MarkupMiddleWare Done for "+path))
	}()

	if path == "" {
		c.AbortWithError(http.StatusNotAcceptable, nil)
		return
	}

	//Get the Session Cookie and Only Authorize /home requests for requests without a session.

	token := session_functions.GetSessionAuthToken(c)
	if token == "" {
		if !isPageSecurityException(path) {
			ginServer.RenderHTML(constants.HTTP_NOT_AUTHORIZED, c)
			c.AbortWithError(http.StatusNotFound, errors.New("Not Authorized for partial page."))
			return
		}
	}

	markup, err := extensions.ReadFile(settings.WebRoot + "/markup/" + path + "/" + file_name + ".htm")
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	var middleWare payloads.MarkupMiddleware
	middleWare.Html = string(markup[:])

	var contentData []byte
	user, err := session_functions.GetSessionUser(c)
	if err == nil {
		var langFile string
		if user.Language == "en" {
			langFile = "US"
		} else {
			langFile = user.Language
		}
		contentData, _, err = readProductionCachedFile(settings.WebRoot + "/globalization/translations/" + path + "/" + user.Language + "/" + langFile + ".json")
		if err != nil {
			contentData, err = extensions.ReadFile(settings.WebRoot + "/globalization/translations/" + path + "/en/US.json")
			if err != nil {
				session_functions.Dump(err)
			}
		}
	} else {
		localeLanguage := ginServer.GetLocaleLanguage(c)
		contentData, _, err = readProductionCachedFile(settings.WebRoot + "/globalization/translations/" + path + "/" + localeLanguage.Language + "/" + localeLanguage.Language + ".json")
		if err != nil {
			contentData, err = extensions.ReadFile(settings.WebRoot + "/globalization/translations/" + path + "/en/US.json")
			if err != nil {
				session_functions.Dump(err)
			}
		}
	}
	middleWare.PageContent = string(contentData[:])
	handler := getMiddlewareState(c, middleWare)
	contextHandler := session_functions.PassContext(c)
	//Call the controller action if it is present
	if action != "" {

		ctl := getController(path)
		methodToCall := ctl.MethodByName(action)

		if methodToCall.String() == "<invalid Value>" {
			getMiddleWareErrorJSONPayload(c, path, "Failed to Call "+action+" for "+path+".", handler)
			return
		}

		//Load params into a map

		uriParamsData, err := base64.StdEncoding.DecodeString(uriParams)

		if err != nil {
			getMiddleWareErrorJSONPayload(c, path, "Failed to decode uri parameters:  "+err.Error(), handler)
			return
		}

		var obj interface{}
		err = json.Unmarshal(uriParamsData, &obj)

		if err != nil {
			getMiddleWareErrorJSONPayload(c, path, "Failed to unmarshal uri parameters:  "+err.Error(), handler)
			return
		}

		m := obj.(map[string]interface{})

		uriParamsMap := make(map[string]string)

		for key, value := range m {

			switch value.(type) {
			case string:
				uriParamsMap[key] = value.(string)
			case int:
				uriParamsMap[key] = extensions.IntToString(value.(int))
			case float64:
				uriParamsMap[key] = extensions.FloatToString(value.(float64), 10)
			case bool:
				uriParamsMap[key] = extensions.BoolToString(value.(bool))
			case []interface{}:
				uriParamsMap[key] = "undefined"
			default:
				uriParamsMap[key] = "undefined"
			}

		}

		in := []reflect.Value{}
		in = append(in, reflect.ValueOf(contextHandler))
		in = append(in, reflect.ValueOf(uriParamsMap))
		in = append(in, reflect.ValueOf(handler))
		if len(uriParamsMap) > 0 {
			session_functions.Dump(uriParamsMap)
		}
		methodToCall.Call(in)

	} else {

		ctl := getController(path)
		methodToCall := ctl.MethodByName(constants.GET_INDEX_CONTROLLER_METHOD)

		if methodToCall.String() == "<invalid Value>" {
			getMiddleWareJSONPayload(c, path, handler)
			return
		}

		uriParamsMap := make(map[string]string)

		in := []reflect.Value{}
		in = append(in, reflect.ValueOf(contextHandler))
		in = append(in, reflect.ValueOf(uriParamsMap))
		in = append(in, reflect.ValueOf(handler))
		if len(uriParamsMap) > 0 {
			session_functions.Dump(uriParamsMap)
		}
		methodToCall.Call(in)

	}
}

func isPageSecurityException(page string) bool {
	if page == "home" {
		return true
	} else if page == "invitation" {
		return true
	} else if page == "registration" {
		return true
	} else if page == "passwordReset" {
		return true
	} else if page == "recovery" {
		return true
	}
	return false
}

func getMiddlewareState(c *gin.Context, middleWare payloads.MarkupMiddleware) session_functions.ServerResponse {

	return func(redirect string, globalMessage string, globalMessageType string, trace error, transactionId string, v interface{}) {

		stateData, _ := json.Marshal(v)
		middleWare.Json = string(stateData[:])
		if settings.ServerSettings.ReleaseMode == "development" {
			core.TransactionLogMutex.RLock()
			devLog := "\"DeveloperLog\": \"" + base64.StdEncoding.EncodeToString([]byte(core.TransactionLog)) + "\""
			core.TransactionLogMutex.RUnlock()
			if middleWare.Json == "{}" {
				middleWare.Json = "{" + devLog + "}"
			} else {
				middleWare.Json = "{" + devLog + ", " + middleWare.Json[1:]
			}
		}

		middleWare.GlobalMessage = globalMessage
		middleWare.GlobalMessageType = globalMessageType
		if trace != nil {
			middleWare.Trace = appErrors.PrintStackTrace(trace, 3) // Convert to full trace.
		}
		middleWare.Redirect = redirect

		ginServer.RespondJSON(middleWare, c)
	}
}

func getMiddleWareJSONPayload(c *gin.Context, path string, handler session_functions.ServerResponse) {
	// pass context to business logic callers which will return JSON data under "data" key

	vm := viewModel.GetViewModel(path)
	vm.LoadDefaultState()
	handler(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_SUCCESS, nil, PARAM_TRANSACTION_ID_NONE, vm)

}

func getMiddleWareErrorJSONPayload(c *gin.Context, path string, err string, handler session_functions.ServerResponse) {
	// pass context to business logic callers which will return JSON data under "data" key

	vm := viewModel.GetViewModel(path)
	vm.LoadDefaultState()
	handler(PARAM_REDIRECT_NONE, err, PARAM_SNACKBAR_TYPE_SUCCESS, nil, PARAM_TRANSACTION_ID_NONE, vm)

}

func handleWebSocketData(conn *app.WebSocketConnection, c *gin.Context, messageType int, id string, data []byte) {

	var apiRequest SocketAPIRequest

	if strings.Contains(string(data), "\"Thank\"") {
		return
	}

	if strings.Contains(string(data), "\"data\":{}") { //Empty Request
		return
	}

	errMarshal := json.Unmarshal(data, &apiRequest)
	if errMarshal != nil {
		session_functions.Log("handleWebSocketData err in marshall", errMarshal.Error())
		return
	}

	resp := func(c *gin.Context, callbackId int) session_functions.ServerResponse {
		return func(redirect string, globalMessage string, globalMessageType string, trace error, transactionId string, v interface{}) {

			var apiResponse payloads.ApiResponse
			apiResponse.Redirect = redirect
			apiResponse.GlobalMessage = globalMessage
			apiResponse.GlobalMessageType = globalMessageType
			apiResponse.Transactionid = transactionId
			if trace != nil {
				apiResponse.Trace = appErrors.PrintStackTrace(trace, 15) // Convert to full trace.
			}

			stateData, _ := json.Marshal(v)

			apiResponse.State = string(stateData[:])

			var response SocketAPIResponse
			response.ApiResponse = apiResponse
			response.CallbackId = callbackId

			app.ReplyToWebSocketJSON(conn, response)
		}
	}

	if apiRequest.ApiRequest.Action == "SetCurrentPage" && apiRequest.ApiRequest.Controller == "App" {

		meta, ok := app.GetWebSocketMeta(id)
		if ok == false {
			return
		}

		meta.ContextString = apiRequest.ApiRequest.State
		meta.ContextType = "ClientStatus"
		app.SetWebSocketMeta(id, meta)

		return
	}

	responseHandler := resp(c, apiRequest.CallbackId)
	callState(apiRequest.ApiRequest.Controller, apiRequest.ApiRequest.Action, string(apiRequest.ApiRequest.State[:]), c, responseHandler)

}

func handleRespondGzipJSFile(c *gin.Context, data []byte, modTime time.Time) {
	ginServer.RespondGzipJSFile(data, modTime, c)
}

func handleRespondGzipCSSFile(c *gin.Context, data []byte, modTime time.Time) {
	ginServer.RespondGzipJSFile(data, modTime, c)
}

func handleRespondHTML(c *gin.Context, data []byte, modTime time.Time) {
	ginServer.RenderHTML(string(data), c)
}

func handleRespondJSON(c *gin.Context, data []byte, modTime time.Time) {
	c.Header("Content-Type", "application/json")
	c.Writer.Write(data)
}
