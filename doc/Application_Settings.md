#GoCore Application Settings

There are 2 components to GoCore which must be configured within your application:

##buildCore

Create a build package for your application with the following:

	package main

	import (
		"github.com/DanielRenne/GoCore/buildCore"
	)
	
	func main() {
		buildCore.Initialize("src/github.com/DanielRenne/GoCoreHelloWorld")
	}

##app

The GoCore/core/app package is what runs your application.  You must first Initialize() it with the root path of your application.  Then call the Run() method.
	
	package main
	
	import (
		"github.com/DanielRenne/GoCore/core/app"
		_ "github.com/DanielRenne/GoCoreHelloWorld/webAPIs/v1/webAPI"
	)
	
	func main() {
		//Run First.
		app.Initialize("src/github.com/DanielRenne/GoCoreHelloWorld")
	
		//Add your Application Code here.
	
		//Run Last.
		app.Run()
	}

#App Settings

GoCore reads a file located in the root directory of your package called WebConfig.json  If one does not exist `buildCore` will auto generate one for your package with default settings.

##WebConfig.json

There are two root objects to be configured:

###application



	"application":{
	    "domain": "127.0.0.1",
	    "httpPort": 80,
	    "httpsPort": 443, 
	    "releaseMode":"release",
	    "webServiceOnly":false,
	    "info":{
	    	"title": "Hello World Playground",
	    	"description":"A web site to try GoCore.",
	    	"contact":{
	    		"name":"DRenne",
	    		"email":"support@myWebSite.com",
	    		"url":"myWebSite.com"
	    	},
	    	"license": {
	    		"name": "Apache 2.0",
	  			"url": "http://www.apache.org/licenses/LICENSE-2.0.html"
	    	},
	    	"termsOfService":"http://127.0.0.1/terms"
	    },
		"htmlTemplates":{
			"enabled":false,
			"directory":"templates",
			"directoryLevels": 1
		}
	}

At the root of application there are the following fields:

####domain

Tells the application which domain to redirect https traffic to.

####serverFQDN

Currently only used for bootstrap purposes to compare domainName where you want your data inserted

####logGophers

Instead of calling go func on anonymous functions.  Use logger.GoRoutineLogger() and pass the func with a description.  Then setup logGophers to true in your web config to log the time in which a goroutine starts and potentially exits.

####releaseMode

Tells the application to debug and run GIN http routing into release mode.  "release" will enable release.  An empty string will place the application in debug mode.

####webServiceOnly

Tells the application only route web service paths.  NO static file routing will be enabled when set to true.

####customGinLogger

If you plan to write and .Use a custom gin logger in your AppIndex, set to true.  Otherwise the default of false will use the default logger and recovery handler.

####productName

A short name (usually not human with spaces).  Can be used to control which bootstrap information to seed based on the webConfig.json

####versionNumeric

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

####versionDot

Useful to show the users a dot-based version

####info

Tells the application details about the application for swagger.io information and schema.

####htmlTemplates

Tells the application to use HTML templates that conform to the GIN Engine.  See [HTML Rendering in GIN](https://github.com/gin-gonic/gin#html-rendering]).  See [HTML Templates](https://github.com/DanielRenne/GoCore/blob/master/doc/HTML_Templates.md) for more details and examples.

####htmlTemplates

Tells the application to use HTML templates that conform to the GIN Engine.  See [HTML Rendering in GIN](https://github.com/gin-gonic/gin#html-rendering]).  See [HTML Templates](https://github.com/DanielRenne/GoCore/blob/master/doc/HTML_Templates.md) for more details and examples.


###dbConnections

Provides an array of database connections.  Currently GoCore only supports a single database connection.  Future releases will allow for multiple connections and types.

	"dbConnections":[
		{
			"driver" : "boltDB",
			"connectionString" : "db/helloWorld.db"
		}
	]
###Database Connection Examples

###Bolt DB

A NOSQL GOLang native database that runs within your application

		{
			"driver" : "boltDB",
			"connectionString" : "db/helloWorld.db"
		}

###Mongo DB

A NOSQL database that runs outside your application

		{
			"driver" : "mongoDB",
			"connectionString" : "mongodb://myuser:mypass@localhost:40001,otherhost:40001/mydb"
		}

###SQLite3

A SQL Database instance running within your application

		{
			"driver" : "sqlite3",
			"connectionString" : "db/helloWorld.db"
		}

###MYSQL

A SQL Database instance running external to your application

		{
			"driver" : "mysql",
			"connectionString" : " myUsername:myPassword@/HelloWorld"
		}

###MS SQL Server

A SQL Database instance running external to your application

		{
			"driver" : "mssql",
			"connectionString" : "server=myServerAddress;Database=HelloWorld;user id=myUsername;Password=myPassword;Connection Timeout=3000;"
		}
