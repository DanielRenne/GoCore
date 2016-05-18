# Build Core

BuildCore is a package for GoCore that will generate SQL Tables, indexes, primary & foreign keys, NOSQL Collections & Buckets, auto-generate Golang model code files, auto-generate REST web api code files, and swagger.io schema files for your application.

###Configuration Files

This routine will read the webConfig.json file and read versioned schema files located in db/{appName}/schemas/{version}/ to create database files, models, web api's, and swagger.io definitions.  {version} must conform to standard semantic versioning format 1.0.0.  Major, Minor, and Revision.

###WebConfig.json

For information on how to configure webConfig.json see [Application Settings](https://github.com/DanielRenne/GoCore/blob/master/doc/Application_Settings.md).


###Model

If database schema files exist buildCore will process the files and create a folder in your application called models.  Within the models folder your {version}/model folder will be created along with all the model go files.  To import your model into your application here is an example below:

	import (
		"github.com/DanielRenne/GoCoreHelloWorld/models/v1/model"
	)

For more information on the NOSQL model package see [NOSQL Model](https://github.com/DanielRenne/GoCore/blob/master/doc/NOSQL_Schema_Model.md).

###Web API

All auto-generated web api's will route to \api\v{Major}.  For a major version of 1 the web api routes would be set to \api\v1\ 

Semantic Versioning within the db\{appName}\schemas directory will produce multiple routes for each version of the schema.


To enable and initialize your web api's you must import a version of the api into one of your package go files:

	import (
		_ "github.com/DanielRenne/GoCoreHelloWorld/webAPIs/v1/webAPI"
	)

This will add all routes to your schema objects GET, POST, PUT, and DELETE Web API.

###Swagger.io 2.0 

