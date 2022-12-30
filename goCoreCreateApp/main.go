package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/user"
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

func cdPath(path string) {
	err := os.Chdir(path)
	errorOut("cd gopath", err, false)
}

func errorOut(line string, err error, dontExit bool) {
	if err != nil {
		msg := "Errored out: " + err.Error()
		logger.Message(msg+" "+line, logger.RED)
		utils.TalkDirty(msg)
		cdPath(basePath)
		extensions.RemoveDirectory(appName)
		if !dontExit {
			os.Exit(2)
		}
	} else {
		logger.Message("Success: "+line, logger.GREEN)
	}
}

var appName string
var databaseType string
var humanTitle string
var colorPalette string
var basePath string
var profileFile string
var mainCNKeys string
var createGit string
var username string

func main() {

	logger.Message("Welcome to the GoCore createApp tool!  Thank you for using GoCore.", logger.YELLOW)

	//also should ensure first char of appName is lower
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("GoCore AppName (lowercase first byte camel case proect name): ")
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

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Path of module install (please no trailing /) It will create three subdirectories to access your app {your-path-you-enter-here}/github.com/" + username + "/" + appName + "/" + appName + "/: ")
		basePath, _ = reader.ReadString('\n')
		basePath = strings.Trim(basePath, "\n")
		ok := false
		if strings.Index(basePath, " ") == -1 && basePath[len(basePath)-1:] != "/" {
			ok = true
		} else {
			fmt.Println("No spaces please and dont end in / ")
		}
		if ok {
			break
		}
	}
	err := extensions.MkDir(basePath)

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Title of all pages: ")
	humanTitle, _ = reader.ReadString('\n')
	humanTitle = strings.Trim(humanTitle, "\n")

	for {
		reader := bufio.NewReader(os.Stdin)
		logger.Message("For your self signed SSL cert.  Add your full cert information like this: \"/CN=www.mydom.com/O=My Company Name LTD./C=US\" (defaults to this if you just press enter)", logger.GREEN)
		mainCNKeys, _ = reader.ReadString('\n')
		mainCNKeys = strings.Trim(mainCNKeys, "\n")
		if mainCNKeys == "" {
			mainCNKeys = "/CN=www.mydom.com/O=My Company Name LTD./C=US"
		}
		break
	}

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("'mongo' or 'bolt' (defaults mongo due to most support): ")
		databaseType, _ = reader.ReadString('\n')
		databaseType = strings.Trim(databaseType, "\n")
		if databaseType == "" {
			databaseType = "mongo"
		}
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

	createGit = "y"
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Create and commit initial git repository (y or n) (defaults y): ")
		createGit, _ = reader.ReadString('\n')
		createGit = strings.Trim(createGit, "\n")
		if createGit == "" {
			createGit = "y"
		}
		ok := false
		if createGit == "y" || createGit == "n" {
			ok = true
		} else {
			fmt.Println("Invalid type 'n' or 'y'")
		}
		if ok {
			break
		}
	}

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter the name of the file which loads your shell environment (~/.bash_profile) is defaulted if you leave blank: ")
		profileFile, _ = reader.ReadString('\n')
		profileFile = strings.Trim(profileFile, "\n")
		if profileFile == "" {
			profileFile = ".bash_profile"
		}
		ok := false
		if strings.Index(profileFile, "/home") != -1 && strings.Index(profileFile, "/Users") != -1 && strings.Index(profileFile, "~") != -1 {
			ok = true
		} else {
			fmt.Println("Just the .filename, not the path such as ~, /home or /Users/")
		}
		if ok {
			break
		}
		break
	}

	cdPath(basePath)

	camelUpper := strings.ToTitle(string(appName[0])) + string(appName[1:])

	err = extensions.Write(colorPalette, "/tmp/colorPalette")
	errorOut("extensions.WriteToFile "+colorPalette+" to /tmp/colorPalette", err, false)

	err = extensions.Write(humanTitle, "/tmp/humanTitle")
	errorOut("extensions.WriteToFile "+humanTitle+" to /tmp/humanTitle", err, false)

	err = extensions.Write(mainCNKeys, "/tmp/mainCNKeys")
	errorOut("extensions.WriteToFile "+mainCNKeys+" to /tmp/mainCNKeys", err, false)

	err = extensions.Write(databaseType, "/tmp/databaseType")
	errorOut("extensions.WriteToFile "+databaseType+" to /tmp/databaseType", err, false)

	err = extensions.Write(username, "/tmp/username")
	errorOut("extensions.WriteToFile "+databaseType+" to /tmp/username", err, false)

	path := "github.com/" + username + "/" + appName
	err = extensions.MkDir(path)
	errorOut("extensions.MkDir("+path+", 0644)", err, false)

	cmd := exec.Command("go", "install", "github.com/DanielRenne/GoCore/getAppTemplate@latest")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	errorOut("running go install github.com/DanielRenne/GoCore/getAppTemplate@latest", err, false)

	talk("Getting all dependencies and the latest version of GoCore App Templates")
	cmd = exec.Command("getAppTemplate")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	errorOut("running getAppTemplate", err, false)

	fmt.Println("App name :", appName)
	appPath := path + "/" + appName

	_, err = os.Stat(appPath)
	if err == nil {
		extensions.RemoveDirectoryShell(appPath)
	}

	modelBuildPath := appPath + "/modelBuild" + camelUpper + "/"
	err = extensions.MkDir(modelBuildPath)
	errorOut("os.MkdirAll("+modelBuildPath+", 0644)", err, false)

	buildPath := appPath + "/build" + camelUpper + "/"

	err = extensions.MkDir(buildPath)
	errorOut("os.MkdirAll("+buildPath+", 0644)", err, false)

	template := `
package main

import (
	"flag"
	"os"
	"github.com/DanielRenne/GoCore/%s"
)

func main() {
	// allow -configFile=test.json to be passed to build different configs other than webConfig.json
	configFile := flag.String("configFile", "webConfig.json", "Configuration File Name.  Ex...  webConfig.json")
	flag.Parse()
	%s.Initialize(os.Getenv("%s_path"), *configFile)
}

`
	buildGoFile := buildPath + "build" + camelUpper + ".go"

	err = extensions.WriteAndGoFormat(heredoc.Docf(template, "buildCore", "buildCore", appName), buildGoFile)
	errorOut("extensions.WriteAndGoFormat "+buildGoFile, err, false)

	modelGoFile := modelBuildPath + "modelBuild" + camelUpper + ".go"
	err = extensions.WriteAndGoFormat(heredoc.Docf(template, "modelBuild", "modelBuild", appName), modelGoFile)
	errorOut("extensions.WriteAndGoFormat "+modelGoFile, err, false)

	talk("Copying app generation files")

	cdPath(basePath + "/" + appPath)

	os.Setenv(appName+"_path", basePath+"/"+appPath)

	user, err := user.Current()
	if err != nil {
		errorOut("couldnt process user", err, false)
	}

	fBashProfile, err := os.OpenFile(user.HomeDir+"/"+profileFile,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		errorOut("couldnt create environment variable in the "+user.HomeDir+".bash_profile", err, false)
	}
	defer fBashProfile.Close()
	if _, err := fBashProfile.WriteString("\nexport " + appName + "_path=" + basePath + "/" + appPath + "\n"); err != nil {
		errorOut("couldnt create environment variable in .bash_profile", err, false)
	}

	cmd = exec.Command("go", "mod", "init", "github.com/"+username+"/"+appName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	errorOut("running go mod init", err, false)

	if createGit == "y" {
		talk("Adding github files")

		cmd = exec.Command("git", "init")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		errorOut("git init", err, false)

		cmd = exec.Command("git", "add", ".")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		errorOut("git add", err, false)

		cmd = exec.Command("git", "commit", "-m", "Initial GoCore App Generation")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		errorOut("git commit", err, false)
	}

	cmd = exec.Command("go", "get", "github.com/DanielRenne/GoCore")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	errorOut("running go get on GoCore", err, false)

	cmd = exec.Command("go", "install", "github.com/DanielRenne/GoCore/getAppTemplate")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	errorOut("running go install on getAppTemplate", err, false)

	cdPath(basePath + "/" + appPath + "/build" + camelUpper + "/")

	cmd = exec.Command("go", "mod", "init", "scaffoldGoCoreApp")
	err = cmd.Run()
	errorOut("go mod init scaffoldGoCoreApp", err, false)

	cmd = exec.Command("go", "get", "github.com/DanielRenne/GoCore/buildCore@a114cdfbeccce193d17f900e919f9b69b1dc9ef9")
	err = cmd.Run()
	errorOut("github.com/DanielRenne/GoCore/buildCore@a114cdfbeccce193d17f900e919f9b69b1dc9ef9", err, false)

	cmd = exec.Command("go", "run", "build"+camelUpper+"/build"+camelUpper+".go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	errorOut("running go run "+buildGoFile, err, false)

	cmd = exec.Command("go", "run", "install"+camelUpper+"/install"+camelUpper+".go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	errorOut("running "+appPath+"/install"+camelUpper, err, false)

	err = os.Chdir(basePath + "/" + appPath + "/bin")
	errorOut("cd bin", err, false)

	cmd = exec.Command("bash", "format")
	err = cmd.Start()
	errorOut("formatting all code", err, false)

	cdPath(basePath + "/" + appPath + "/modelBuild" + camelUpper + "/")

	cmd = exec.Command("go", "install", ".")
	err = cmd.Run()
	errorOut("go install models `"+"go install modelBuild"+camelUpper+"/build"+camelUpper+".go`", err, false)

	cmd = exec.Command("bash", basePath+"/"+appPath+"/bin/start_app")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	errorOut("running gocore app server!", err, false)

}
