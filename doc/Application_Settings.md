#GoCore Application Settings

There are 2 files in GoCore which must be configured to point to your application:

##App.go

Add your package path to `github.com/DanielRenne/GoCore/core/app/app.go`

	package app
	
	import (
	
		//--------Change Below for the application you want to run --------
		_ "github.com/DanielRenne/GoCoreHelloWorld"
	)
	
	func init() {
	
	}

##AppGenFiles.go

Add your package path INCLUDING src/ to `github.com/DanielRenne/GoCore/core/appGen/appGenFiles.go` to the APP_LOCATION constant.

	package appGen
	
	import (
	
		// "fmt"
		"github.com/DanielRenne/GoCore/core/extensions"
		"github.com/DanielRenne/GoCore/core/log"
		"os"
	)
	
	const APP_LOCATION = "src/github.com/DanielRenne/GoCoreHelloWorld"

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

####httpPort, httpsPort

Tells the application which ports to listen on for http and https.

####releaseMode

Tells the application to debug and run GIN http routing into release mode.  "release" will enable release.  An empty string will place the application in debug mode.

####webServiceOnly

Tells the application only route web service paths.  NO static file routing will be enabled when set to true.

####info

Tells the application details about the application for swagger.io information and schema.

####htmlTemplates

Tells the application to use HTML templates that conform to the GIN Engine.  See [HTML Rendering in GIN](https://github.com/gin-gonic/gin#html-rendering]).


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
