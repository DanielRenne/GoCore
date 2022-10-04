// Package serverSettings provides a global object for storing server settings or local GoCore configurations for package usages
// Most of this is used with GoCore full mode.	In GoCore lite mode, this is not used.
package serverSettings

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/DanielRenne/GoCore/core/path"
)

// HtmlTemplates is a map of all the html static templates that are loaded into memory for gin-gonic
type HtmlTemplates struct {
	Enabled         bool   `json:"enabled"`
	Directory       string `json:"directory"`
	DirectoryLevels int    `json:"directoryLevels"`
}

// DbConnection is a struct for storing the database connection information for mongo or boltdb.  Replication support is very limited
type DbConnection struct {
	ConnectionString    string `json:"connectionString"`
	EnableTLS           bool   `json:"enableTLS"`
	Driver              string `json:"driver"`
	Database            string `json:"database"`
	AuthServer          bool   `json:"authServer"`
	TransactionSizeMax  int    `json:"transactionSizeMax"`
	AuditHistorySizeMax int    `json:"auditHistorySizeMax"`
	Replication         struct {
		Enabled    bool     `json:"enabled"`
		ReplicaSet string   `json:"replicaSet"`
		Master     string   `json:"master"`
		Slaves     []string `json:"slaves"`
	} `json:"replication"`
}

// License is legacy (for when GoCore did swagger generation) and needs to be deprecated soon.
type License struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Contact is legacy (for when GoCore did swagger generation) and needs to be deprecated soon.
type Contact struct {
	Name  string `json:"name"`
	URL   string `json:"url"`
	Email string `json:"email"`
}

// Application is a struct for storing the application settings
type Application struct {
	// Name is the name of the application
	Name string `json:"name"`
	// Domain Tells the application which domain to redirect https traffic to.
	Domain string `json:"domain"`
	// ServerFQDN is the fully qualified domain name of the server (mostly used in our react apps for where to point the webpack server)
	// But can be also used for bootstrapping mongo or boltdb data for purposes to compare domainName where you want your data inserted to insert different data based on a domain  (See core/dbServices/createDBServices.go#genNoSQLBootstrap)
	ServerFQDN string `json:"serverFQDN"`
	// VersionNumeric is the numeric version of the application (typically unused)
	VersionNumeric int `json:"versionNumeric"`
	// VersionDot is the dot version of the application for use to compile your builds against (typically unused but can be used to seed and bootstrap different data during bootstrap See core/dbServices/createDBServices.go#genNoSQLBootstrap)
	VersionDot string `json:"versionDot"`
	// ProductName is the name of the product (typically unused but can be used to seed and bootstrap different data during bootstrap See core/dbServices/createDBServices.go#genNoSQLBootstrap)
	ProductName string `json:"productName"`
	// HttpPort is the port that the application will listen on for http traffic
	HttpPort int `json:"httpPort"`
	// HttpsPort is the port that the application will listen on for https traffic
	HttpsPort int `json:"httpsPort"`
	// CookieDomain is the domain that the application will set cookies for
	CookieDomain string `json:"cookieDomain"`
	// ReleaseMode is the release mode of the application for gin to be in production or debug mode (can be used to seed and bootstrap different data during bootstrap See core/dbServices/createDBServices.go#genNoSQLBootstrap)
	ReleaseMode string `json:"releaseMode"`
	// WebServiceOnly is a flag to tell the application to only run the web service. NO static file routing will be enabled when set to true.
	WebServiceOnly bool `json:"webServiceOnly"`
	// MountGitWebHooks - Deprecated - use github.com actions instead
	MountGitWebHooks bool `json:"mountGitWebHooks"`
	// GitWebHookSecretKey - Deprecated - use github.com actions instead
	GitWebHookSecretKey string `json:"gitWebHookSecretKey"`
	// GitWebHookPort - Deprecated - use github.com actions instead
	GitWebHookPort string `json:"gitWebHookServerPort"`
	// GitWebHookPath - Deprecated - use github.com actions instead
	GitWebHookPath string `json:"gitWebHookPath"`
	// HtmlTemplates is a struct for storing the html template settings
	HtmlTemplates HtmlTemplates `json:"htmlTemplates"`
	// RootIndexPath is the path to the index.html file for the application if Application/HtmlTemplates/Enabled is false it will be a file path you wish to load in serverSettings.APP_LOCATION+"/web/" RootIndexPath is the file name of your index.html
	RootIndexPath string `json:"rootIndexPath"`
	// If HtmlTemplates is off and you dont want anythin in the index path, you can set this to true and it will serve nothing but a 404
	DisableRootIndex bool `json:"disableRootIndex"`
	// DisableWebSockets is a flag to tell the application to disable websockets
	DisableWebSockets bool `json:"disableWebSockets"`
	// SessionKey is the key used to encrypt the session cookie
	SessionKey string `json:"sessionKey"`
	// SessionName is the name of the session cookie
	SessionName string `json:"sessionName"`
	// SessionExpirationDays is the number of days the session cookie will expire
	SessionExpirationDays int `json:"sessionExpirationDays"`
	// SessionSecureCookie is a flag to tell the application to use a secure cookie
	SessionSecureCookie bool `json:"sessionSecureCookie"`
	// CSRFSecret is the secret used to encrypt the csrf token.  Please dont leave blank or something guessable
	CSRFSecret string `json:"csrfSecret"`
	// BootstrapData is a flag to tell the application to bootstrap data into the database, definitely set to true
	BootstrapData bool `json:"bootstrapData"`
	// LogQueries is a flag to tell the application to log detailed queries to the log
	LogQueries bool `json:"logQueries"`
	// LogQueryStackTraces is a flag to tell the application to log stack traces for queries to the log
	LogQueryStackTraces bool `json:"logQueryStackTraces"`
	// LogJoinQueries is a flag to tell the application to log detailed queries to the log
	LogJoinQueries bool `json:"logJoinQueries"`
	// LogQueryTimes is a flag to tell the application to log detailed queries to the log
	LogQueryTimes bool `json:"logQueryTimes"`
	// LogGophers is a flag to tell the application to log gophers who are wrapped with logger.GoRoutineLogger to the log
	LogGophers bool `json:"logGophers"`
	// LogGopherInterval is the interval in seconds to log gophers who are wrapped with logger.GoRoutineLogger to the log
	LogGopherInterval int `json:"logGopherInterval"`
	// CoreDebugStackTrace is a flag to tell the application to log stack traces for queries to the log if you use core.Debug.* functions to log
	CoreDebugStackTrace bool `json:"coreDebugStackTrace"`
	// AllowCrossOriginRequests is a flag to tell the application to allow cross origin requests
	AllowCrossOriginRequests bool `json:"allowCrossOriginRequests"`
}

// WebConfigType is a struct for storing the web configuration settings
type WebConfigType struct {
	// DbConnection is the database connection information for mongo or boltdb.  Replication support is very limited
	DbConnections []DbConnection `json:"dbConnections"`
	// APplication is the application settings
	Application Application `json:"application"`
	// DbConnection is mostly for internal stuff.  Dont set this yourself
	DbConnection DbConnection
}

// WebConfig is the global object for storing server settings or local GoCore configurations for package usages
var WebConfig WebConfigType

// WebConfigMutex is a mutex for the webConfig object
var WebConfigMutex sync.RWMutex

// APP_LOCATION is the path to the application.  This is set by the Initialize function and should not be set by you.  Only read.
var APP_LOCATION string

// Init initializes the webConfig object
func Init() {
	Initialize(path.GetBinaryPath(), "webConfig.json")
}

// InitCustomWebConfig initializes the webConfig object with a custom webConfig.json file
func InitCustomWebConfig(webConfig string) {
	Initialize(path.GetBinaryPath(), webConfig)
}

func initPath(path string) {
	APP_LOCATION = path
}

// Initialize initializes the webConfig object (typically this is handled outside of your code in buildCore and app packages so it needs to be exported for them)
func Initialize(path string, configurationFile string) (err error) {
	initPath(path)
	fmt.Println("core serverSettings initialized.")

	jsonData, err := ioutil.ReadFile(APP_LOCATION + "/" + configurationFile)
	if err != nil {
		fmt.Println("Reading of webConfig.json failed:  " + err.Error())
	}

	WebConfigMutex.Lock()
	errUnmarshal := json.Unmarshal(jsonData, &WebConfig)
	if errUnmarshal != nil {
		fmt.Println("Parsing / Unmarshaling of webConfig.json failed:  " + errUnmarshal.Error())
	}

	for _, dbConnection := range WebConfig.DbConnections {
		WebConfig.DbConnection = dbConnection
	}
	if WebConfig.Application.LogGopherInterval == 0 {
		WebConfig.Application.LogGopherInterval = 15
	}
	WebConfigMutex.Unlock()

	return
}
