# dbServices Models PlayGround module

This module is a playground for the dbServices module. It is a simple example of how to use the dbServices mongo database abstraction ORM. The ORM for goCore mongo/boltdb

Not all of our methods are documented as to how to filter your data in your mongo app. We built this ORM for a specific project and if you plan on using this, please [reach out to Dave](mailto:dnxglya4@duck.com) and we can help you understand things better.

## Getting Started

First ensure you have mongo installed (It also supports boltdb but this package is deprecated now by the authors). The webConfig.json example points to localhost on the default port to the `test` database that will be generated.

For this example, we committed the models based off the sample JSON schema. You `should not do this in your app`. Developers and your build services to build your app will always just run the below command to generate the models into the models/v1/model directory prior to compiling a test run or build of your application. But since we committed the files for these examples you dont have to run this unless you are changing the example db/bootstrap or db/schemas

```bash
go build -o genModels modelsGenerate/main.go
./genModels
```

Anytime you change something inside of db/bootstrap or db/schemas, you must re-run this so that the embedded data gets compiled into the go model generation. [Please see the docs on model ORMs before playing with stuff](https://github.com/DanielRenne/GoCore/blob/master/doc/NOSQL_Schema_Model.md)

To load the database filtering and insertion main, run `go build main.go && ./main`

You will see a bunch of dumped data starting with

```
  ___                               _     _
 |_ _|  _ __    ___    ___   _ __  | |_  (_)   ___    _ __    ___
  | |  | '_ \  / __|  / _ \ | '__| | __| | |  / _ \  | '_ \  / __|
  | |  | | | | \__ \ |  __/ | |    | |_  | | | (_) | | | | | \__ \
 |___| |_| |_| |___/  \___| |_|     \__| |_|  \___/  |_| |_| |___/

```

Try to gain and understanding by [reading the main source code](https://github.com/DanielRenne/GoCore/blob/master/core/dbServices/example/main.go) and also playing with some sample keys I provided in webConfig.json in which will bootstrap new records based upon how you have your application configured bootstrapping data based upon the product you are building or whether its in development mode.
