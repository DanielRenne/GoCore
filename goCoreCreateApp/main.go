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
	var createGit string
	var createGitUsername string
	var privateRepo string
	var pushGit string
	var useSSH string
	var gitPassword string
	var colorPalette string
	var basePath string

	logger.Message("Welcome to the GoCore createApp tool!  Thank you for using GoCore.", logger.YELLOW)

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

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Path of application install (no trailing /): ")
		basePath, _ = reader.ReadString('\n')
		basePath = strings.Trim(basePath, "\n")
		ok := false
		if strings.Index(basePath, " ") == -1 {
			ok = true
		} else {
			fmt.Println("No spaces please")
		}
		if ok {
			break
		}
	}
	err := os.MkdirAll(basePath, 0777)

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

	if createGit == "y" {

		for {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("If public repo.  Would you like to push this to github.com? (defaults y): ")
			pushGit, _ = reader.ReadString('\n')
			pushGit = strings.Trim(pushGit, "\n")
			if pushGit == "" {
				pushGit = "y"
			}
			ok := false
			if pushGit == "y" || pushGit == "n" {
				ok = true
			} else {
				fmt.Println("Invalid type 'n' or 'y'")
			}
			if ok {
				break
			}
		}
	}

	if pushGit == "y" {

		for {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Use SSH keys locally or password based to push? [y=ssh n=username/password]: ")
			useSSH, _ = reader.ReadString('\n')
			useSSH = strings.Trim(createGitUsername, "\n")
			if useSSH == "" {
				useSSH = "y"
			}
			ok := false
			if useSSH == "y" || useSSH == "n" {
				ok = true
			} else {
				fmt.Println("Enter y or n")
			}
			if ok {
				break
			}
		}

		for {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter github.com username to push as [will send \"" + username + "\"] change if this is a team account to your local login you use for your team: ")
			createGitUsername, _ = reader.ReadString('\n')
			createGitUsername = strings.Trim(createGitUsername, "\n")
			if createGitUsername == "" {
				createGitUsername = username
			}
			ok := false
			if createGitUsername != "" {
				ok = true
			} else {
				fmt.Println("Enter username")
			}
			if ok {
				break
			}
		}

		for {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Private Repo (defaults n): ")
			privateRepo, _ = reader.ReadString('\n')
			privateRepo = strings.Trim(privateRepo, "\n")
			if privateRepo == "" {
				privateRepo = "n"
			}
			ok := false
			if privateRepo == "y" || privateRepo == "n" {
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
			fmt.Print("Enter github.com password for user \"" + createGitUsername + "\": ")
			gitPassword, _ = reader.ReadString('\n')
			gitPassword = strings.Trim(gitPassword, "\n")
			ok := false
			if gitPassword != "" {
				ok = true
			} else {
				fmt.Println("Enter password with at least 4 bytes")
			}
			if ok {
				break
			}
		}
	}

	//logger.Message("Next choose a color palette:", logger.WHITE)
	//logger.Message("(default) BlueGrey and Orange value=bgo", logger.BLUE)
	//logger.Message("(default) Green and White value=irish", logger.GREEN)

	//reader = bufio.NewReader(os.Stdin)
	//fmt.Print("Color value: ")
	//colorPalette, _ = reader.ReadString('\n')
	//colorPalette = strings.Trim(colorPalette, "\n")

	cdPath(basePath)

	camelUpper := strings.ToTitle(string(appName[0])) + string(appName[1:])

	err = extensions.WriteToFile(colorPalette, "/tmp/colorPalette", 0777)
	errorOut("extensions.WriteToFile "+colorPalette+" to /tmp/colorPalette", err, false)

	err = extensions.WriteToFile(humanTitle, "/tmp/humanTitle", 0777)
	errorOut("extensions.WriteToFile "+humanTitle+" to /tmp/humanTitle", err, false)

	err = extensions.WriteToFile(databaseType, "/tmp/databaseType", 0777)
	errorOut("extensions.WriteToFile "+databaseType+" to /tmp/databaseType", err, false)

	path := "github.com/" + username
	err = os.MkdirAll(path, 0777)
	errorOut("os.MkdirAll("+path+", 0644)", err, false)

	talk("Getting all dependencies and the latest version of GoCore App Templates")
	cmd := exec.Command("getAppTemplate")
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

	err = os.MkdirAll(appPath, 0777)
	errorOut("os.MkdirAll("+appPath+", 0644)", err, false)

	modelBuildPath := appPath + "/modelBuild" + camelUpper + "/"

	err = os.MkdirAll(modelBuildPath, 0777)
	errorOut("os.MkdirAll("+modelBuildPath+", 0644)", err, false)

	buildPath := appPath + "/build" + camelUpper + "/"

	err = os.MkdirAll(buildPath, 0777)
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

	fBashProfile, err := os.OpenFile(user.HomeDir+"/.bash_profile",
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

	err = os.Chdir(basePath + "/" + appPath)
	errorOut("cd appPath", err, false)

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

		if useSSH == "y" {
			cmd = exec.Command("git", "remote", "add", "origin", "git@github.com:"+username+"/"+appName+".git")
			err = cmd.Run()
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			errorOut("git remote add", err, false)
		} else {
			cmd = exec.Command("git", "remote", "add", "origin", "https://github.com/"+username+"/"+appName+".git")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err = cmd.Run()
			errorOut("git remote add", err, false)
		}

		if pushGit == "y" {
			talk("Creating repository online")
			pathExec := "/tmp/execCurl"
			var endpoint string
			if username == createGitUsername {
				endpoint = "https://api.github.com/user/repos"
			} else {
				endpoint = "https://api.github.com/orgs/" + username + "/repos"
			}
			payload := `"{\"name\": \"` + appName + `\"}"`
			if privateRepo == "y" {
				payload = `"{\"private\": true, \"name\": \"` + appName + `\"}"`
			}
			err := extensions.WriteToFile("curl -u "+createGitUsername+":"+gitPassword+" "+endpoint+" -d "+payload, pathExec, 777)
			errorOut("extensions.WriteToFile /tmp/execCurl", err, false)
			cmd = exec.Command("bash", pathExec)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err = cmd.Run()
			talk("Done creating repository online")
			errorOut("curl create repo on API", err, true)

			if useSSH == "n" {
				logger.Message("\n\nRun this after completion.\n\ncd "+os.Getenv("GOPATH")+"/"+appPath+"\ngit push -u "+username+" origin master\n\n\nThen enter your password", logger.MAGENTA)
			} else {
				err = os.Chdir(basePath + "/" + appPath)
				errorOut("cd appPath", err, false)
				cmd = exec.Command("git", "push", "origin", "master")
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				err = cmd.Run()
				errorOut("Git push", err, true)
			}
			err = os.Remove(pathExec)
			errorOut("Remove "+pathExec, err, true)
		}
	}

	cmd = exec.Command("go", "install", strings.Replace(modelBuildPath, "src/", "", -1))
	err = cmd.Run()
	errorOut("go install models `"+"go install "+strings.Replace(modelBuildPath, "src/", "", -1)+"`", err, false)

	cmd = exec.Command("bash", basePath+"/"+appPath+"/bin/start_app")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	errorOut("running gocore app server!", err, false)

}
