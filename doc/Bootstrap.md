# Bootstrapping the database

## Bootstrapping basics

GoCore's dbServices package (which is included in a GoCore full app) provides a facility which will preload your database before any reads will be allowed. It waits until each collection has been fully bootstrapped before letting in any reads when the system boots up. It also can allow you to have some granular control over what records are inserted based on situations through the "BootstrapMeta" key you will inject inside your raw JSON data of your records

## Folders

To bootstrap data, you must first create a db/ folder to house db/bootstrap and db/schemas. A small collection bootstrap must be inside the subfolder of the collection name. So if your collection is called accounts, the bootstrap file must be called db/bootstrap/accounts/accounts.json and the contents will be a slice of many objects. These files will be cached inside the models after you execute something [like this](https://github.com/DanielRenne/GoCore/blob/master/core/dbServices/example/modelsGenerate/main.go) see the [example](https://github.com/DanielRenne/GoCore/tree/master/core/dbServices/example) for more details about dbServices and models and bootstrapping.

## dist/ folder

Alternatively if you have a large amount of json to bootstrap, you may simply have one record per file and dump it into db/bootstrap/{collectionName}/dist/folders. We will quickly pickup all files and get them into the db as fast as possible.

## Bootstrap control

One neat thing you can do with our bootstrapping of data is you can have some control to insert different things based upon your webConfig.json

So inside of a bootstrap record you can always add the following struct

```
"BootstrapMeta": {
	"Version": 0,
	"Domain": "",
	"ReleaseMode": "",
	"ProductName": "",
	"Domains": [""],
	"ProductNames": [""],
	"DeleteRow": false,
	"AlwaysUpdate": false,
}
```

DeleteRow and AlwaysUpdate allow you to control when your server boots up, to update the record or to delete it if its deprecated. But the other keys are comparing what is inside your webConfig.json with the values in each row as to whether to bootstrap the record.

Here is a mapping of what key maps to what webConfig.json key

| BootStrapMeta | WebConfig.json "application" object |
| ------------- | ----------------------------------- |
| Version       | versionNumeric                      |
| Domain        | serverFQDN                          |
| Domains       | serverFQDN                          |
| ReleaseMode   | releaseMode                         |
| ProductName   | productName                         |
| ProductNames  | productName                         |

For more details about dbServices model generation see [the example](https://github.com/DanielRenne/GoCore/tree/master/core/dbServices/example) for more details and a working example outside of GoCore Full
