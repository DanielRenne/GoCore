package main

import (
	//"bufio"
	//"strings"
	"fmt"
	"os"
	"github.com/cloud-ignite/GoCore/core/utils"
	"github.com/cloud-ignite/GoCore/core/logger"
	"github.com/cloud-ignite/GoCore/core/extensions"
	"strings"
	"github.com/davidrenne/heredoc"
)

func errorOut(line string, err error) {
	if err != nil {
		msg := "Errored out: " + err.Error()
		logger.Message(msg + line, logger.RED)
		utils.TalkDirty(msg)
		os.Exit(2)
	}
}

func main() {
	var appName string
	var username string
	var pushGithubYN string
	var humanTitle string
	var pushGithub bool
	fmt.Println(pushGithub, appName, username, pushGithubYN, humanTitle)
	appName = "goCoreTest"
	username = "davidrenne"
	pushGithubYN = "Y"
	humanTitle = "david renne site"
	pushGithub = false

	//also should ensure first char of appName is lower

	//for {
	//	reader := bufio.NewReader(os.Stdin)
	//	fmt.Print("GoCore AppName: ")
	//	appName, _ = reader.ReadString('\n')
	//	appName = strings.Trim(appName, "\n")
	//	ok := false
	//	if strings.Index(appName, " ") == -1 {
	//		ok = true
	//	} else {
	//		fmt.Println("No spaces please")
	//	}
	//	if ok {
	//		break
	//	}
	//}
	//
	//reader := bufio.NewReader(os.Stdin)
	//fmt.Print("Title of all pages: ")
	//humanTitle, _ = reader.ReadString('\n')
	//humanTitle = strings.Trim(humanTitle, "\n")
	//
	//for {
	//	reader := bufio.NewReader(os.Stdin)
	//	fmt.Print("Github.com Username: ")
	//	username, _ = reader.ReadString('\n')
	//	username = strings.Trim(appName, "\n")
	//	ok := false
	//	if strings.Index(username, " ") == -1 {
	//		ok = true
	//	} else {
	//		fmt.Println("No spaces please")
	//	}
	//	if ok {
	//		break
	//	}
	//}
	//
	//for {
	//	reader := bufio.NewReader(os.Stdin)
	//	fmt.Print("Push To Github.com (y or n): ")
	//	pushGithubYN, _ = reader.ReadString('\n')
	//	pushGithubYN = strings.Trim(pushGithubYN, "\n")
	//	ok := false
	//	if strings.ToUpper(pushGithubYN) == "Y" || strings.ToUpper(pushGithubYN) == "N" {
	//		ok = true
	//		pushGithub = strings.ToUpper(pushGithubYN) == "Y"
	//	} else {
	//		fmt.Println("y or n please")
	//	}
	//	if ok {
	//		break
	//	}
	//}
	camelUpper := strings.ToTitle(string(appName[0])) + string(appName[1:])

	path := "src/github.com/" + username
	err := os.MkdirAll(path, 0644)
	errorOut("os.MkdirAll(path, 0644)", err)

	appPath := path + "/" + appName
	_, err = os.Stat(appPath)
	if err == nil {
		err := extensions.RemoveDirectory(appPath)
		errorOut("extensions.RemoveDirectory(appPath)", err)
	}

	err = os.MkdirAll(appPath, 0644)
	errorOut("os.MkdirAll(appPath, 0644)", err)

	modelBuildPath := appPath + "/modelBuild" + camelUpper + "/"

	err = os.MkdirAll(modelBuildPath, 0644)
	errorOut("os.MkdirAll(modelBuildPath, 0644)", err)

	buildPath := appPath + "/build" + camelUpper + "/"

	err = os.MkdirAll(buildPath, 0644)
	errorOut("os.MkdirAll(buildPath, 0644)", err)

	template := `
package main

import (
	"flag"
	"github.com/cloud-ignite/GoCore/%s"
)

func main() {
	// allow -configFile=test.json to be passed to build different configs other than webConfig.json
	configFile := flag.String("configFile", "webConfig.json", "Configuration File Name.  Ex...  webConfig.json")
	flag.Parse()
	%s.Initialize("src/github.com/frezadev/rezaSandbox", *configFile)
}

`
	buildGoFile := buildPath + "build" + camelUpper + ".go"
	err = extensions.WriteAndGoFormat(heredoc.Docf(template, "buildCore", "buildCore"), buildGoFile)
	errorOut("extensions.WriteAndGoFormat build", err)

	modelGoFile := modelBuildPath + "modelBuild" + camelUpper + ".go"
	err = extensions.WriteAndGoFormat(heredoc.Docf(template, "modelBuild", "modelBuild"), modelGoFile)
	errorOut("extensions.WriteAndGoFormat model", err)


}

