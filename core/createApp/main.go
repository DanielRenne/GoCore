package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"github.com/cloud-ignite/GoCore/core/utils"
)

func main() {
	var appName string
	var username string
	var pushGithubYN string
	var humanTitle string
	var pushGithub bool
	fmt.Println(pushGithub)
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

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Push To Github.com (y or n): ")
		pushGithubYN, _ = reader.ReadString('\n')
		pushGithubYN = strings.Trim(pushGithubYN, "\n")
		ok := false
		if strings.ToUpper(pushGithubYN) == "Y" || strings.ToUpper(pushGithubYN) == "N" {
			ok = true
			pushGithub = strings.ToUpper(pushGithubYN) == "Y"
		} else {
			fmt.Println("y or n please")
		}
		if ok {
			break
		}
	}

	utils.TalkDirtyToMe()

}

