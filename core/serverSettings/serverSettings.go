package serverSettings

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
)

type htmlTemplates struct {
	Enabled         bool   `json:"enabled"`
	Directory       string `json:"directory"`
	DirectoryLevels int    `json:"directoryLevels"`
}

type dbConnection struct {
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

type license struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type contact struct {
	Name  string `json:"name"`
	URL   string `json:"url"`
	Email string `json:"email"`
}

type info struct {
	Title          string  `json:"title"`
	Description    string  `json:"description"`
	Contact        contact `json:"contact"`
	License        license `json:"license"`
	TermsOfService string  `json:termsOfService"`
}

type Application struct {
	Name                     string        `json:"name"`
	Domain                   string        `json:"domain"`
	ServerFQDN               string        `json:"serverFQDN"`
	VersionNumeric           int           `json:"versionNumeric"`
	VersionDot               string        `json:"versionDot"`
	ProductName              string        `json:"productName"`
	HttpPort                 int           `json:"httpPort"`
	HttpsPort                int           `json:"httpsPort"`
	CookieDomain             string        `json:"cookieDomain"`
	ReleaseMode              string        `json:"releaseMode"`
	WebServiceOnly           bool          `json:"webServiceOnly"`
	MountGitWebHooks         bool          `json:"mountGitWebHooks"`
	GitWebHookSecretKey      string        `json:"gitWebHookSecretKey"`
	GitWebHookPort           string        `json:"gitWebHookServerPort"`
	GitWebHookPath           string        `json:"gitWebHookPath"`
	Info                     info          `json:"info"`
	HtmlTemplates            htmlTemplates `json:"htmlTemplates"`
	RootIndexPath            string        `json:"rootIndexPath"`
	DisableRootIndex         bool          `json:"disableRootIndex"`
	CustomGinLogger          bool          `json:"customGinLogger"`
	SessionKey               string        `json:"sessionKey"`
	SessionName              string        `json:"sessionName"`
	SessionExpirationDays    int           `json:"sessionExpirationDays"`
	SessionSecureCookie      bool          `json:"sessionSecureCookie"`
	CSRFSecret               string        `json:"csrfSecret"`
	BootstrapData            bool          `json:"bootstrapData"`
	LogQueries               bool          `json:"logQueries"`
	LogQueryStackTraces      bool          `json:"logQueryStackTraces"`
	LogJoinQueries           bool          `json:"logJoinQueries"`
	LogQueryTimes            bool          `json:"logQueryTimes"`
	LogGophers               bool          `json:"logGophers"`
	CoreDebugStackTrace      bool          `json:"coreDebugStackTrace"`
	AllowCrossOriginRequests bool          `json:"allowCrossOriginRequests"`
	SecondaryWebsocketPath   string        `json:"secondaryWebsocketPath"`
}

type webConfigObj struct {
	DbConnections []dbConnection `json:"dbConnections"`
	Application   Application    `json:"application"`
	DbConnection  dbConnection
}

var WebConfig webConfigObj
var WebConfigMutex sync.RWMutex
var APP_LOCATION string
var SWAGGER_UI_PATH string

func Initialize(path string, configurationFile string) (err error) {

	APP_LOCATION = path
	SWAGGER_UI_PATH = APP_LOCATION + "/web/swagger/dist"
	fmt.Println("core serverSettings initialized.")

	jsonData, err := ioutil.ReadFile(APP_LOCATION + "/" + configurationFile)
	if err != nil {
		fmt.Println("Reading of webConfig.json failed:  " + err.Error())
		return
	}

	WebConfigMutex.Lock()
	errUnmarshal := json.Unmarshal(jsonData, &WebConfig)
	if errUnmarshal != nil {
		fmt.Println("Parsing / Unmarshaling of webConfig.json failed:  " + errUnmarshal.Error())
		return
	}

	for _, dbConnection := range WebConfig.DbConnections {
		WebConfig.DbConnection = dbConnection
	}
	WebConfigMutex.Unlock()

	return
}
