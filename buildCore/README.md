# Build Core

Build Core is a package for GoCore that will generate SQL Tables, indexes, primary & foreign keys, NOSQL Collections & Buckets, auto-generate Golang model code files, auto-generate REST web api code files, and swagger.io schema files for your application.

###Configuration Files

This routine will read the webConfig.json file and read versioned schema files located in db\\{appName}\schemas\{version}\ to create database files, models, web api's, and swagger.io definitions.  {version} must conform to standard semantic versioning format 1.0.0.  Major, Minor, and Revision.

####WebConfig.json
...


###Web API

All auto-generated web api's will route to \api\v{Major}.  For a major version of 1 the web api routes would be set to \api\v1\ 

Semantic Versioning within the db\{appName}\schemas directory will produce multiple routes of versioning.


To use the model package you will import into your source files the following:

	import( 
		"helloWorld/model"
	)