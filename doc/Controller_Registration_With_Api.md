## Utilizing controller reflection with GoCore/core/app/api

Add a file in the package controllers called controllers.go

```go
package controllers

import (
	_ "yourgomodapp/controllers/reports"
)
```

Add a directory called reports in the controllers directory and add this file

```go
package reports

import (
	"github.com/DanielRenne/GoCore/core/app/api"
)

type Get struct {
	Id  string `json:"Id"`
}

type ErrorResponse struct {
	Error struct {
		Message    string `json:"Message"`
		Code       string `json:"Code"`
		Stacktrace string `json:"StackTrace"`
	} `json:"Error"`
}

//Reports controller provides http GET and POST requests for the reports.
type Reports struct{}

func init() {
	api.RegisterController(&Reports{})
}

func (d *Reports) FilteredReportById(obj Get) (y interface{}) {
  if obj.Id == "" {
		response := ErrorResponse{}
		response.Error.Message = "Failed to fetch user id"
		y = response
		return
	}
  type report struct {
    Date      time.Time          `json:"Date"`
  }
  y = report{Date: time.Now()}
  return y
}
```

[//]: # Todo... add more examples of the javascript side for this

```go
package main

import (
	"os"
	"os/signal"
  _ "yourgomodapp/controllers"
	"github.com/DanielRenne/GoCore/core/app"
	"github.com/DanielRenne/GoCore/core/app/api"
	"github.com/DanielRenne/GoCore/core/extensions"
	"github.com/DanielRenne/GoCore/core/ginServer"
)

func main() {
	port := os.Getenv("PORT")
	app.InitializeLite(false, []string{})

	ginServer.Router.GET("/apiGET", api.APICallback)
	ginServer.Router.POST("/apiPOST", api.APICallback)

	if port == "" {
		go app.RunLite(80)
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, os.Kill)
		<-c
	} else {
		app.RunLite(extensions.StringToInt(port))
	}
}

```
