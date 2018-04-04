# GoCore

A Golang solution of tools for building a full stack web application.

## Goals of the Project ##

Below are some targeted goals:

* Http & Https Redirection & HTTP 2.0 with Golang 1.8 and [gin-gonic/gin](https://github.com/gin-gonic/gin)
* Extension packages for common functions including:
	* File IO Management & Manipulation
	* Zip File Compression & Decompression 

----------

* Database extensions and drivers for the following:
	* Supported databases:
		* MongoDB 
		* BoltDB
	* Create SQL Schema (DDL) from JSON Configuration.
	* Create Golang ORM packages for RDBMS Transactions & Queries.

----------
## Getting Started ##

### Install Golang 1.8 or greater and setup your GOPATH ###
[Windows Golang Setup](http://www.wadewegner.com/2014/12/easy-go-programming-setup-for-windows/ "Windows Golang Setup")

[Linux & MAC Setup](https://golang.org/doc/install/ "Linux & MAC Setup")

### Set Your GOPATH to your Go Workspace

### Add to your Operating Systems Path the GOPATH/bin directory.


### Get GoCore
1.  To start a new project with GoCore run the following steps in a new console window.

	`go get github.com/DanielRenne/GoCore/...`

NOTE:  You will see an output message (no buildable Go source files in ....  Please disregaurd.

2.  GoCore comes with a getCore package which can be used to get all external dependencies and custom files as well as example applications.

	`go install github.com/DanielRenne/GoCore/getCore`

Then run

    getCore

### Build GoCore App

1.  A sample project is available to build a GoCore project.

	`go install github.com/DanielRenne/GoCore/core/goCoreCreateApp`

Then run

    goCoreCreateApp

Follow the prompts to generate your app and it will copy the templates found in GoCore/tools/appFiles into your project directory github.com/username/appName.  Note you probably should install nodejs, npm and nvm before generating an app.

GoCore has built in functions to read json configuration files to generate SQL Tables, indexes, primary & foreign keys, NOSQL Collections & Buckets, auto-generated Golang model code files, auto-generated REST web api code files, and swagger.io schema files application.  See [BuildCore Readme](https://github.com/DanielRenne/GoCore/blob/master/buildCore/README.md) for details.

### Run Your GoCore App

1.  Run the following to start your app

	`bash bin/start_app`

Open a web browser to:  [http://127.0.0.1:8080](http://127.0.0.1:8080)

Login as admin/admin and setup application roles etc.  More documentation to come later on application specific setup.

#### How to build your own web project in GoCore

See [Application Settings](https://github.com/DanielRenne/GoCore/blob/master/doc/Application_Settings.md) within docs for information on what webConfig.json allows for.

## IMPORTANT NOTE for HTTPS (TLS) Security
GoCore comes default with an open ssl generated Cert and PEM files that are **NOT** secret as they are available via the open source repository.  Make sure you replace both these files located in the `keys` directory.  To do this we recommend in Linux running the following command and copying the output files to the `keys` directory.  Alternatively cert and pem files generated by a valid Certificate Authority like GoDaddy or Verisign when you reach production with an online domain.

	openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout key.pem -out cert.pem

NOTE:  We also recommend **NOT** to store your secret key files in source control.  We recommend ignoring the keys directory for source control.

NOTE:  key files are ingored when running `getCore` as to not overwrite your keys.

Additional Info on Golang https:  [https://www.kaihag.com/https-and-go/](https://www.kaihag.com/https-and-go/ "HTTPS&GO")



## Building a Database Model with SQLite3

####NOTE: Because SQLite3 requires gcc externally we separated the driver for compiling reasons.  Windows users we recommend installing gcc as a prerequisite for the sqlite3 golang module to compile via [tdb-gcc](http://tdm-gcc.tdragon.net/download).  Be sure to install 64 bit for 64 bit machines. 

####More SQLite tools to verify your data in Windows [SQLite Studio](http://sqlitestudio.pl/)

To create a SQLite3 Database schema and model package for your application run the following:

	go install github.com/DanielRenne/GoCore/buildCoreLite

Then run

	buildCoreLite  

## References

* [NOSQL Database Schema Model API](https://github.com/DanielRenne/GoCore/blob/master/doc/NOSQL_Schema_Model.md)
