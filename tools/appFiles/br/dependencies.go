package br

import (
	"os/exec"
	"runtime"

	"github.com/DanielRenne/GoCore/core/extensions"
	"github.com/DanielRenne/goCoreAppTemplate/constants"
	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
)

type dependenciesBr struct{}

type dependencyPackage struct {
	PackageName string `json:"PackageName"`
	Description string `json:"Description"`
	GOOS        string `json:"GOOS"`
	OSRelease   string `json:"OSRelease"`
	FileName    string `json:"FileName"`
}

type dependencyManifest struct {
	Dependencies []dependencyPackage `json:"Dependencies"`
}

//InstallDependencies will read a directory and read a manifest to install dependencies.
func (dependenciesBr) InstallDependencies(path string) {
	var manifest dependencyManifest
	err := extensions.ReadFileAndParse(path+constants.PATH_SEPARATOR+"manifest.json", &manifest)
	if err != nil {
		session_functions.Log("Error->br->dependencies->InstallDependencies", "Failed to Read manifest file and parse.")
		return
	}

	release := GetOSVersion()

	for i := range manifest.Dependencies {
		dep := manifest.Dependencies[i]
		if dep.GOOS == runtime.GOOS {
			if dep.OSRelease == "" || dep.OSRelease == release {
				session_functions.Log("Installing Package", dep.PackageName)
				err = exec.Command("/usr/bin/sudo", "/usr/bin/apt-get", "install", "-y", dep.PackageName).Run()
				if err != nil {
					session_functions.Log("Error->br->dependencies->InstallDependencies", "Error running apt-get install -y "+dep.PackageName+":  "+err.Error())
					continue
				}
				session_functions.Log("Successfully Installed Package", dep.PackageName)
			}
		}
	}
}
