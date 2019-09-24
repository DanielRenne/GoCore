package main

import (
	"fmt"
	"github.com/DanielRenne/GoCore/core/app"
	"github.com/DanielRenne/GoCore/core/extensions"
	"github.com/atlonaeng/studio/settings"
	"strings"
	"syscall"
)

func main() {
	app.Initialize("src/github.com/atlonaeng/studio", "webConfig.json")
	settings.Initialize()
	page, err := extensions.ReadFile(settings.WebRoot + "/javascript/pages/template/{page}.js")
	if err != nil {
		fmt.Println(err)
		syscall.Exit(1)
	}
	newPage := strings.Replace(string(page), "page", "adsfad", 0)
	err = extensions.WriteToFile(string(newPage), settings.WebRoot+"/javascript/pages/template/{page}-tmp.js", 777)
	if err != nil {
		fmt.Println(err)
		syscall.Exit(1)
	}
	fmt.Println(string(newPage))
}
