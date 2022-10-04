## GoCore Lite

Create a main.go

```go
package main

import (
	"mime"
	"os"
	"os/signal"
	"path"
	"time"

	"github.com/DanielRenne/GoCore/core/app"
	"github.com/DanielRenne/GoCore/core/extensions"
	"github.com/DanielRenne/GoCore/core/ginServer"
	"github.com/gin-gonic/gin"
)

func addSecureHeaders(c *gin.Context) {
	c.Writer.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
	c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
	c.Writer.Header().Set("X-Frame-Options", "SAMEORIGIN")
	serve_SSL_only := os.Getenv("SERVE_SSL_ONLY")
	if serve_SSL_only == "true" {
		c.Writer.Header().Set("Content-Security-Policy", "default-src https: 'unsafe-eval' 'unsafe-inline'; font-src 'self' data: https:; img-src 'self' blob: data: https:;media-src 'self' blob: data: https:; object-src 'self' blob: data: https; connect-src 'self' wss: https:")
	}
	c.Writer.Header().Set("X-XSS-Protection", "1")
}

func main() {
	port := os.Getenv("PORT")
	app.InitializeLite(false, []string{})

	ginServer.Router.NoRoute(func(c *gin.Context) {
		addSecureHeaders(c)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	ginServer.Router.GET("/", func(c *gin.Context) {
		addSecureHeaders(c)
		var data []byte
		data = append(data, []byte("log.Println('Hello GoCore Lite');")...)
		// Possibly read assets like an index.html template through the use of an awesome project https://github.com/go-bindata/go-bindata which compiles all assets into your go binary.
		// data, err := assets.Asset("web/dist/index.html")
		// if err != nil {
		// 	ginServer.RenderHTML("Error: "+err.Error(), c)
		// 	return
		// }
		ext := path.Ext(c.Request.URL.String())

		c.Writer.Header().Set("Content-Type", mime.TypeByExtension(ext))
		c.Writer.Header().Set("Content-Length", extensions.IntToString(len(data)))

		ginServer.RespondJSFile(data, time.Now(), c)

	})

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

Then run `go run main.go $(pwd)`

Follow the steps outlined and parameters to generate your backend only goCore application

Then run `go get -d ./...` to download all the dependencies of your main.go

If you look at the output it will show you your command to go run like: `go run main.go $(pwd)`

When you build your main.go you dont have to pass the `pwd` as the first parameter. It's only needed in go run due to temporary compile directories and needing to know where your webConfig and referenced files are located.

## GoCore core/app package

The GoCore/core/app package is what runs your application. You must first Init() it with the root path of your application. Then call the Run() method which will block on the HTTP server initialization.

Please note, that if you use go run, you must call `go run main.go $(pwd)` because GoCore needs to know the directory of your project to read the webConfig.json file and associated paths for things like models, keys, db/bootstrap etc. You can also call app.InitCustomWebConfig() and pass a custom file name for webConfig.json based on your own logic.

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
