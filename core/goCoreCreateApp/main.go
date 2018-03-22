package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/DanielRenne/GoCore/core/extensions"
	"github.com/DanielRenne/GoCore/core/logger"
	"github.com/DanielRenne/GoCore/core/utils"
	"github.com/davidrenne/heredoc"
)

func talk(msg string) {
	utils.TalkDirty(msg)
	logger.Message("Message: "+msg, logger.GREEN)
}

func cdGoPath() {
	err := os.Chdir(os.Getenv("GOPATH"))
	errorOut("cd gopath", err, false)
}

func errorOut(line string, err error, dontExit bool) {
	if err != nil {
		msg := "Errored out: " + err.Error()
		logger.Message(msg+" "+line, logger.RED)
		utils.TalkDirty(msg)
		if !dontExit {
			os.Exit(2)
		}
	} else {
		logger.Message("Success: "+line, logger.GREEN)
	}
}

func main() {
	var appName string
	var username string
	var databaseType string
	var humanTitle string
	var colorPalette string

	logger.Message("Welcome to the GoCore createApp tool!  Thank you for using GoCore.", logger.YELLOW)
	logger.Message("We hold these below truths to be self-evident", logger.WHITE)
	logger.Message(fmt.Sprintf("% x", []byte{100, 97, 118, 105, 100, 32, 114, 101, 110, 110, 101, 32, 105, 115, 32, 99, 111, 111, 108, 33, 10}), logger.MAGENTA)

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
		username = strings.Trim(username, "\n")
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

	databaseType = "mongo"
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("'mongo' or 'bolt' (defaults mongo due to most support): ")
		databaseType, _ = reader.ReadString('\n')
		databaseType = strings.Trim(databaseType, "\n")
		ok := false
		if databaseType == "mongo" || databaseType == "bolt" {
			ok = true
		} else {
			fmt.Println("Invalid type 'mongo' or 'bolt'")
		}
		if ok {
			break
		}
	}
	//logger.Message("Next choose a color palette:", logger.WHITE)
	//logger.Message("(default) BlueGrey and Orange value=bgo", logger.BLUE)
	//logger.Message("(default) Green and White value=irish", logger.GREEN)

	//reader = bufio.NewReader(os.Stdin)
	//fmt.Print("Color value: ")
	//colorPalette, _ = reader.ReadString('\n')
	//colorPalette = strings.Trim(colorPalette, "\n")

	cdGoPath()

	camelUpper := strings.ToTitle(string(appName[0])) + string(appName[1:])

	err := extensions.WriteToFile(colorPalette, "/tmp/colorPalette", 644)
	errorOut("extensions.WriteToFile "+colorPalette+" to /tmp/colorPalette", err, false)

	err = extensions.WriteToFile(humanTitle, "/tmp/humanTitle", 644)
	errorOut("extensions.WriteToFile "+humanTitle+" to /tmp/humanTitle", err, false)

	err = extensions.WriteToFile(databaseType, "/tmp/databaseType", 644)
	errorOut("extensions.WriteToFile "+databaseType+" to /tmp/databaseType", err, false)

	path := "src/github.com/" + username
	err = os.MkdirAll(path, 0644)
	errorOut("os.MkdirAll("+path+", 0644)", err, false)

	fmt.Println("adsffdas name :", appName)
	fmt.Println(path)
	appPath := path + "/" + appName

	_, err = os.Stat(appPath)
	if err == nil {
		err := extensions.RemoveDirectory(appPath)
		errorOut("extensions.RemoveDirectory(appPath)", err, false)
	}

	err = os.MkdirAll(appPath, 0644)
	errorOut("os.MkdirAll("+appPath+", 0644)", err, false)

	modelBuildPath := appPath + "/modelBuild" + camelUpper + "/"

	err = os.MkdirAll(modelBuildPath, 0644)
	errorOut("os.MkdirAll("+modelBuildPath+", 0644)", err, false)

	buildPath := appPath + "/build" + camelUpper + "/"

	err = os.MkdirAll(buildPath, 0644)
	errorOut("os.MkdirAll("+buildPath+", 0644)", err, false)

	template := `
package main

import (
	"flag"
	"github.com/DanielRenne/GoCore/%s"
)

func main() {
	// allow -configFile=test.json to be passed to build different configs other than webConfig.json
	configFile := flag.String("configFile", "webConfig.json", "Configuration File Name.  Ex...  webConfig.json")
	flag.Parse()
	%s.Initialize("src/github.com/` + username + "/" + appName + `", *configFile)
}

`
	buildGoFile := buildPath + "build" + camelUpper + ".go"
	err = extensions.WriteAndGoFormat(heredoc.Docf(template, "buildCore", "buildCore"), buildGoFile)
	errorOut("extensions.WriteAndGoFormat "+buildGoFile, err, false)

	modelGoFile := modelBuildPath + "modelBuild" + camelUpper + ".go"
	err = extensions.WriteAndGoFormat(heredoc.Docf(template, "modelBuild", "modelBuild"), modelGoFile)
	errorOut("extensions.WriteAndGoFormat "+modelGoFile, err, false)

	talk("Copying app generation files")

	cmd := exec.Command("go", "run", buildGoFile)
	err = cmd.Run()
	errorOut("running go run "+buildGoFile, err, false)

	cmd = exec.Command("go", "run", appPath+"/install"+camelUpper+"/install"+camelUpper+".go")
	err = cmd.Run()
	errorOut("running "+appPath+"/install"+camelUpper, err, false)

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

	cmd = exec.Command("git", "remote", "add", "origin", "https://github.com/"+username+"/"+appName+".git")
	err = cmd.Run()
	errorOut("git commit", err, false)

	cdGoPath()
	cmd = exec.Command("go", "install", strings.Replace(modelBuildPath, "src/", "", -1))
	err = cmd.Run()
	errorOut("go install models", err, false)

	cmd = exec.Command("bash", appPath+"/bin/start_app")
	err = cmd.Start()
	errorOut("running gocore app server!", err, false)

}
