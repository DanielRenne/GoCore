// Package settings provides settings for go files to reference paths and other constants to know where files are located.
package settings

import (
	"encoding/json"
	"github.com/DanielRenne/GoCore/core/serverSettings"
	"io/ioutil"
	"log"
	"sync"
)

const webRoot = "/web/app"
const markup = "/web/app/markup"

var Version string
var ProductName string
var Environment string

type appSettings struct {
	ServerAdmins                       []string `json:"serverAdmins"`
	DeveloperMode                      bool     `json:"developerMode"`
	DeveloperGoTrace                   bool     `json:"developerGoTrace"`
	DeveloperLogState                  bool     `json:"developerLogState"`
	DeveloperLogTheseObjects           []string `json:"developerLogTheseObjects"`
	DeveloperSuppressTheseObjects      []string `json:"developerSuppressTheseObjects"`
	DeveloperSuppressThesePages        []string `json:"developerSuppressThesePages"`
	DeveloperLogStateChangePerformance bool     `json:"developerLogStateChangePerformance"`
	DeveloperLogReact                  bool     `json:"developerLogReact"`
	DemoMode                           bool     `json:"demoMode"`
}

type webConfig struct {
	AppSettings appSettings `json:"appSettings"`
}

type FullWebConfig struct {
	Application struct {
		Domain                      string `json:"domain"`
		HTTPPort                    int    `json:"httpPort"`
		HTTPSPort                   int    `json:"httpsPort"`
		ReleaseMode                 string `json:"releaseMode"`
		WebServiceOnly              bool   `json:"webServiceOnly"`
		DisableRootIndex            bool   `json:"disableRootIndex"`
		SessionKey                  string `json:"sessionKey"`
		SessionName                 string `json:"sessionName"`
		SessionExpirationDays       int    `json:"sessionExpirationDays"`
		SessionSecureCookie         bool   `json:"sessionSecureCookie"`
		CsrfSecret                  string `json:"csrfSecret"`
		VersionNumeric              int    `json:"versionNumeric"`
		VersionDot                  string `json:"versionDot"`
		ProductName                 string `json:"productName"`
		BootstrapData               bool   `json:"bootstrapData"`
		LogQueries                  bool   `json:"logQueries"`
		LogQueryStackTraces         bool   `json:"logQueryStackTraces"`
		LogJoinQueries              bool   `json:"logJoinQueries"`
		LogQueryTimes               bool   `json:"logQueryTimes"`
		FlushCoreDebugToStandardOut bool   `json:"flushCoreDebugToStandardOut"`
		CoreDebugStackTrace         bool   `json:"coreDebugStackTrace"`
		Info                        struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			Contact     struct {
				Name  string `json:"name"`
				Email string `json:"email"`
				URL   string `json:"url"`
			} `json:"contact"`
			License struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"license"`
			TermsOfService string `json:"termsOfService"`
		} `json:"info"`
		HTMLTemplates struct {
			Enabled         bool   `json:"enabled"`
			Directory       string `json:"directory"`
			DirectoryLevels int    `json:"directoryLevels"`
		} `json:"htmlTemplates"`
	} `json:"application"`
	DbConnections []struct {
		Driver              string `json:"driver"`
		ConnectionString    string `json:"connectionString"`
		Database            string `json:"database"`
		TransactionSizeMax  int    `json:"transactionSizeMax"`
		AuditHistorySizeMax int    `json:"auditHistorySizeMax"`
		Replication         struct {
			Enabled    bool     `json:"enabled"`
			ReplicaSet string   `json:"replicaSet"`
			Master     string   `json:"master"`
			Slaves     []string `json:"slaves"`
		} `json:"replication"`
	} `json:"dbConnections"`
	AppSettings struct {
		DeveloperMode                      bool          `json:"developerMode"`
		DeveloperLogState                  bool          `json:"developerLogState"`
		DeveloperLogStateChangePerformance bool          `json:"developerLogStateChangePerformance"`
		DeveloperLogTheseObjects           []interface{} `json:"developerLogTheseObjects"`
		DeveloperLogReact                  bool          `json:"developerLogReact"`
		IsDemo                             bool          `json:"isDemo"`
	} `json:"appSettings"`
}

var WebRoot string
var WebUI string
var AppSettings appSettings
var ServerSettings serverSettings.Application

var AppSettingsSync sync.RWMutex

func Initialize() {
	log.Println("Settings Initialized.")
	if serverSettings.APP_LOCATION == "" {
		log.Println("server settings APP_LOCATION is blank.  This is not right!")
		return
	}
	WebRoot = serverSettings.APP_LOCATION + webRoot
	WebUI = serverSettings.APP_LOCATION + markup
	ServerSettings = serverSettings.WebConfig.Application

	jsonData, err := ioutil.ReadFile(serverSettings.APP_LOCATION + "/webConfig.json")
	if err != nil {
		log.Println("Reading of webConfig.json failed at settings.init():  " + err.Error())
		return
	}

	var config webConfig

	errUnmarshal := json.Unmarshal(jsonData, &config)
	if errUnmarshal != nil {
		log.Println("Parsing / Unmarshaling of webConfig.json failed:  " + errUnmarshal.Error())
		return
	}

	AppSettingsSync.Lock()
	AppSettings = config.AppSettings
	AppSettingsSync.Unlock()
}
