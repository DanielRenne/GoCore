# GoCore Application Settings

GoCore reads a file located in the root directory of your package called webConfig.json If one does not exist `buildCore` will auto generate one for your package with default settings.

## WebConfig.json

There are two root objects to be configured:

### application

```json
{
  "application": {
    "logGophers": true,
    "logGopherInterval": 15,
    "domain": "0.0.0.0",
    "serverFQDN": "0.0.0.0",
    "httpPort": 80,
    "httpsPort": 443,
    "releaseMode": "development",
    "webServiceOnly": false,
    "versionNumeric": 1,
    "versionDot": "0.0.1",
    "productName": "goCoreProductNameMainProduct",
    "disableRootIndex": true,
    "sessionKey": "goCoreSessionKey",
    "sessionName": "goCoreProductName",
    "sessionExpirationDays": 3650,
    "sessionSecureCookie": false,
    "csrfSecret": "goCoreCsrfSecret",
    "bootstrapData": true,
    "htmlTemplates": {
      "enabled": false,
      "directory": "templates",
      "directoryLevels": 1
    }
  },
  "dbConnections": [
    {
      "driver": "mongoDB",
      "connectionString": "mongodb://127.0.0.1:27017/goCoreProductName",
      "database": "goCoreProductName"
    }
  ]
}
```

At the root of application there are the following fields:

#### domain

Tells the application which domain to redirect https traffic to.

#### serverFQDN

Currently only used for bootstrap purposes to compare domainName where you want your data inserted

#### logGophers

Instead of calling go func on anonymous functions. Use logger.GoRoutineLogger() and pass the func with a description. Then setup logGophers to true in your web config to log the time in which a goroutine starts and potentially exits.

#### releaseMode

Tells the application to debug and run GIN http routing into release mode. "release" will enable release. An empty string will place the application in debug mode.

#### webServiceOnly

Tells the application only route web service paths. NO static file routing will be enabled when set to true.

#### productName

A short name (usually not human with spaces). Can be used to control which bootstrap information to seed based on the webConfig.json

#### versionNumeric

This is used primarily for bootstrapping data with the BootstrapMeta struct to tag your seeds with how you want them to run

```
type BootstrapMeta struct {
	Version     int    `json:"Version" bson:"Version"`
	Domain      string `json:"Domain" bson:"Domain"`
	ReleaseMode string `json:"ReleaseMode" bson:"ReleaseMode"`
	InfoTitle   string `json:"PostCode" bson:"PostCode"`
	DeleteRow   bool   `json:"DeleteRow" bson:"DeleteRow"`
}
```

#### versionDot

Useful to show the users a dot-based version

#### htmlTemplates

Tells the application to use HTML templates that conform to the GIN Engine. See [HTML Rendering in GIN](https://github.com/gin-gonic/gin#html-rendering]). See [HTML Templates](https://github.com/DanielRenne/GoCore/blob/master/doc/HTML_Templates.md) for more details and examples.

#### htmlTemplates

Tells the application to use HTML templates that conform to the GIN Engine. See [HTML Rendering in GIN](https://github.com/gin-gonic/gin#html-rendering]). See [HTML Templates](https://github.com/DanielRenne/GoCore/blob/master/doc/HTML_Templates.md) for more details and examples.

### dbConnections

Provides an array of database connections. Currently GoCore only supports a single database connection. Future releases will allow for multiple connections and types.

    "dbConnections":[
    	{
    		"driver" : "boltDB",
    		"connectionString" : "db/helloWorld.db"
    	}
    ]

### Database Connection Examples

### Bolt DB

A NOSQL GOLang native database that runs within your application

    	{
    		"driver" : "boltDB",
    		"connectionString" : "db/helloWorld.db"
    	}

### Mongo DB

A NOSQL database that runs outside your application

    	{
    		"driver" : "mongoDB",
    		"connectionString" : "mongodb://myuser:mypass@localhost:40001,otherhost:40001/mydb"
    	}
