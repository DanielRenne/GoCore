package main

import (
	"bufio"
	"strings"
	"fmt"
	"os"
	"github.com/cloud-ignite/GoCore/core/utils"
	"github.com/cloud-ignite/GoCore/core/logger"
	"github.com/cloud-ignite/GoCore/core/extensions"
	"github.com/davidrenne/heredoc"
	"os/exec"
)

func talk(msg string) {
	utils.TalkDirty(msg)
	logger.Message("Message: " + msg, logger.GREEN)
}

func cdGoPath() {
	err := os.Chdir(os.Getenv("GOPATH"))
	errorOut("cd gopath", err, false)
}

func errorOut(line string, err error, dontExit bool) {
	if err != nil {
		msg := "Errored out: " + err.Error()
		logger.Message(msg + " " + line, logger.RED)
		utils.TalkDirty(msg)
		if !dontExit {
			os.Exit(2)
		}
	} else {
		logger.Message("Success: " + line, logger.GREEN)
	}
}

func main() {
	var appName string
	var username string
	var humanTitle string

	logger.Message("Welcome to the GoCore createApp tool!  Thank you for using GoCore.", logger.YELLOW)
	logger.Message("We hold these below truths to be self-evident", logger.WHITE)
	logger.Message(fmt.Sprintf("% x", []byte{100,97,118 ,105, 100, 32 ,114, 101, 110, 110 ,101, 32 ,105, 115 ,32 ,99 ,111 ,111, 108, 33, 10}), logger.MAGENTA)

	//also should ensure first char of appName is lower
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("GoCore AppName: ")
		appName, _ = reader.ReadString('\n')
		appName = strings.Trim(appName, "\n")
		ok := false
		if strings.Index(appName, " ") == -1 {
			ok = true
		} else {
			fmt.Println("No spaces please")
		}
		if ok {
			break
		}
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Title of all pages: ")
	humanTitle, _ = reader.ReadString('\n')
	humanTitle = strings.Trim(humanTitle, "\n")

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Github.com Username: ")
		username, _ = reader.ReadString('\n')
		username = strings.Trim(appName, "\n")
		ok := false
		if strings.Index(username, " ") == -1 {
			ok = true
		} else {
			fmt.Println("No spaces please")
		}
		if ok {
			break
		}
	}

	cdGoPath()

	camelUpper := strings.ToTitle(string(appName[0])) + string(appName[1:])

	err := extensions.WriteToFile(humanTitle, "/tmp/humanFile", 644)
	errorOut("extensions.WriteToFile "+ humanTitle + " to /tmp/humanFile", err, false)

	path := "src/github.com/" + username
	err = os.MkdirAll(path, 0644)
	errorOut("os.MkdirAll(" + path + ", 0644)", err, false)

	appPath := path + "/" + appName
	_, err = os.Stat(appPath)
	if err == nil {
		err := extensions.RemoveDirectory(appPath)
		errorOut("extensions.RemoveDirectory(appPath)", err, false)
	}

	err = os.MkdirAll(appPath, 0644)
	errorOut("os.MkdirAll(" + appPath + ", 0644)", err, false)

	modelBuildPath := appPath + "/modelBuild" + camelUpper + "/"

	err = os.MkdirAll(modelBuildPath, 0644)
	errorOut("os.MkdirAll(" + modelBuildPath + ", 0644)", err, false)

	buildPath := appPath + "/build" + camelUpper + "/"

	err = os.MkdirAll(buildPath, 0644)
	errorOut("os.MkdirAll(" + buildPath + ", 0644)", err, false)

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
	%s.Initialize("src/github.com/` + username + "/" + appName + `", *configFile)
}

`
	buildGoFile := buildPath + "build" + camelUpper + ".go"
	err = extensions.WriteAndGoFmt(heredoc.Docf(template, "buildCore", "buildCore"), buildGoFile, true, 644)
	errorOut("extensions.WriteAndGoFormat "+ buildGoFile, err, false)

	modelGoFile := modelBuildPath + "modelBuild" + camelUpper + ".go"
	err = extensions.WriteAndGoFmt(heredoc.Docf(template, "modelBuild", "modelBuild"), modelGoFile, true, 644)
	errorOut("extensions.WriteAndGoFormat " + modelGoFile, err, false)

	talk("Copying app generation files")

	cmd := exec.Command("go", "run", buildGoFile)
	err = cmd.Run()
	errorOut("running go run " + buildGoFile, err, false)

	cmd = exec.Command("go", "run", appPath + "/install" + camelUpper + "/install" + camelUpper + ".go")
	err = cmd.Run()
	errorOut("running " + appPath + "/install" + camelUpper , err, false)

	cdGoPath()
	err = os.Chdir(appPath + "/bin")
	errorOut("cd bin", err, false)

	cmd = exec.Command("bash", "format")
	err = cmd.Start()
	errorOut("formatting all code", err, false)

	cdGoPath()
	err = os.Chdir(appPath)
	errorOut("cd appPath", err, false)

	talk("Adding github files")

	cmd = exec.Command("git", "init")
	err = cmd.Run()
	errorOut("git init", err, false)

	cmd = exec.Command("git", "add", ".")
	err = cmd.Run()
	errorOut("git add", err, false)

	cmd = exec.Command("git", "commit", "-m", "Initial GoCore App Generation")
	err = cmd.Run()
	errorOut("git commit", err, false)

	cmd = exec.Command("git", "remote", "add", "origin", "https://github.com/" + username + "/" + appName + ".git")
	err = cmd.Run()
	errorOut("git commit", err, false)

	cdGoPath()
	cmd = exec.Command("go", "install", strings.Replace(modelBuildPath, "src/", "", -1))
	err = cmd.Run()
	errorOut("go install models", err, false)

	cmd = exec.Command("bash", appPath + "/bin/model_build")
	err = cmd.Run()
	errorOut("bash " + appPath + "/bin/model_build", err, false)

	cmd = exec.Command("bash", appPath + "/bin/start_app")
	err = cmd.Start()
	errorOut("running gocore app server!", err, false)

}

