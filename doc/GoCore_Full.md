## Use the buildCore package to generate a webConfig.json file and related files for your new GoCore backend only application

Create a build binary inside a new modelBuild directory `mkdir modelBuild; cd modelBuild; go mod init go mod init github.com/davidrenne/GoCoreFullExample/modelBuild` (preferably named modelBuild.go in your base application folder like the example snippet)

The purpose of this binary is that anytime your schema JSON files have changed to regerate the model structs and functions. So once you create the file, go install it and run it anytime your models change:

```go
    package main

    import (
    	"github.com/DanielRenne/GoCore/buildCore"
    )

    func main() {
    	buildCore.Init()
    }
```

Then run `go mod tidy` which will download buildCore package

Then run `cd ..; go get github.com/DanielRenne/GoCore/buildCore; go build -o appModelBuild modelBuild/modelBuild.go && ./appModelBuild`

Follow the steps outlined and parameters to generate your backend only goCore application:

```
Please enter the camelCase name of your app
goCoreFullExample
We are now attempting to generate SSL self signed certificates.  Add your full cert information like this: "/CN=www.mydom.com/O=My Company Name LTD./C=US" (defaults to this if you just press enter)
What do you want to call your main package fileName?
myApp
What is your go module name?
github.com/davidrenne/GoCoreFullExample
Do you want to include cron jobs to your main.go? ('y' or 'n')
y
```

Then run `go get -d ./...` to download all the dependencies of your main.go.  Disregard `go get github.com/gomodule/redigo@latest` being retracted.  We still need to resolve this issue in our dependencies.

If you look at the output above it will look something like this `GoCore myApp.go and other files generated in your module successfully.  Please 'go build myApp.go && ./myApp' to get started running your app for the first time` it will show you your command to run your app for the first time.

## GoCore core/app package

The GoCore/core/app package is what runs your application. You must first Init() it with the root path of your application. Then call the Run() method which will block on the HTTP server initialization.

Please note, that if you use go run, you must call `go run main.go $(pwd)` because GoCore needs to know the directory of your project to read the webConfig.json file and associated paths for things like models, keys, db/bootstrap etc. You can also call app.InitCustomWebConfig() and pass a custom file name for webConfig.json based on your own logic. It's best to just call `go build` on your main server and then run it vs. go run and passing the current working directory.

```go
    package main
    import (
      "github.com/DanielRenne/GoCore/core/app"
    )
    func main() {
      //Run First.
      app.Init()
      //Add your Application Code here.
      //Run Last.
      app.Run()
    }
```

#### How to build your own web project in GoCore

See [Application Settings](https://github.com/DanielRenne/GoCore/blob/master/doc/Application_Settings.md) within docs for information on what webConfig.json allows for.

## References

For information on how to build out noSQL models in our ORM, see [this markdown](https://github.com/DanielRenne/GoCore/blob/master/doc/NOSQL_Schema_Model.md)
