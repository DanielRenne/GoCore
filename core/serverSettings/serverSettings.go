package serverSettings

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/DanielRenne/GoCore/core/path"
)

type HtmlTemplates struct {
	Enabled         bool   `json:"enabled"`
	Directory       string `json:"directory"`
	DirectoryLevels int    `json:"directoryLevels"`
}

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

type License struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Contact struct {
	Name  string `json:"name"`
	URL   string `json:"url"`
	Email string `json:"email"`
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
	HtmlTemplates            HtmlTemplates `json:"htmlTemplates"`
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
}

type WebConfigType struct {
	DbConnections []DbConnection `json:"dbConnections"`
	Application   Application    `json:"application"`
	DbConnection  DbConnection
}

var WebConfig WebConfigType
var WebConfigMutex sync.RWMutex
var APP_LOCATION string

func Init() {
	Initialize(path.GetBinaryPath(), "webConfig.json")
}

func InitCustomWebConfig(webConfig string) {
	Initialize(path.GetBinaryPath(), webConfig)
}

func InitPath(path string) {
	APP_LOCATION = path
}

func Initialize(path string, configurationFile string) (err error) {
	InitPath(path)
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
	WebConfigMutex.Unlock()

	return
}
