package main

import (
	"core/extensions"
	"core/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const goCoreURL = "https://github.com/DanielRenne/GoCore/archive/master.zip"
const fileName = "goCoreMaster.zip"
const gitRepoName = "GoCore-master"
const manifestFileName = "downloadedManifest.json"
const manifestURL = "https://raw.githubusercontent.com/DanielRenne/GoCore/master/coreManifest.json"

type version struct {
	Version       string   `json:"version"`
	ReleaseURL    string   `json:"releaseURL"`
	GoDirectories []string `json:"goDirectories"`
}

type coreManifest struct {
	Versions []version `json:"versions"`
}

func main() {

	downloadManifest()

	manifest, errManifest := readManifest()
	if errManifest != nil {
		fmt.Println("Failed to parse manifest file:  " + errManifest.Error())
		return
	}

	fmt.Printf("%+v\n", manifest)

	return

	downloadRelease()

	fmt.Println("Unzipping goCoreMaster.zip . . .")

	excludedFiles := []string{gitRepoName + "/webConfig.json"}
	excludedFiles = append(excludedFiles, gitRepoName+"/src/core/app/app.go")
	excludedFiles = append(excludedFiles, gitRepoName+"/keys/cert.pem")
	excludedFiles = append(excludedFiles, gitRepoName+"/keys/key.pem")

	errUnzip := zip.Unzip(fileName, "goCoreMaster", excludedFiles)

	if errUnzip != nil {
		fmt.Println("Failed to Unzip goCoreMaster.zip:  " + errUnzip.Error())
		return
	}

	fmt.Println("Unzipped GoCore Master successfully.")

	errRemoveRepoZip := os.Remove(fileName)

	if errRemoveRepoZip != nil {
		fmt.Println("Failed to Remove GoCore Repo Zip File:  " + errRemoveRepoZip.Error())
		return
	}

	fmt.Println("Moving Files . . .")

	dir, _ := os.Getwd()
	dir = strings.Replace(dir, " ", "\\ ", -1)

	//Copy the Files then Remove the Directory
	extensions.CopyFolder("goCoreMaster/"+gitRepoName, dir)

	fmt.Println("Moved Files Successfully.")
	fmt.Println("Started Cleaning Files.")
	errRemoveDecompressedFiles := extensions.RemoveDirectory("goCoreMaster")

	if errRemoveDecompressedFiles != nil {
		fmt.Println("Failed to Remove GoCore Decompressed Files:  " + errRemoveDecompressedFiles.Error())
		return
	}

	fmt.Println("Cleaned up Repo Files Successfully.")
}

func downloadRelease() {

	fmt.Println("Downloading GoCore.  Please Wait . . .")

	out, errCreateFile := os.Create(fileName)

	if errCreateFile != nil {
		fmt.Println("Failed to create file handle:  " + errCreateFile.Error())
		return
	}

	resp, errHttpGet := http.Get(goCoreURL)
	defer resp.Body.Close()

	if errHttpGet != nil {
		fmt.Println("Failed to Download GoCore master zip file:  " + errHttpGet.Error())
		return
	}

	n, errCopyOut := io.Copy(out, resp.Body)

	if errCopyOut != nil {
		fmt.Println("Failed to Output to goCoreMaster.zip:  " + errCopyOut.Error())
		return
	}
	out.Close()

	fmt.Println("Downloaded GoCore Master successfully:  " + extensions.PrintMegaBytes(n))
}

func downloadManifest() {

	out, errCreateFile := os.Create(manifestFileName)

	if errCreateFile != nil {
		fmt.Println("Failed to create file handle:  " + errCreateFile.Error())
		return
	}

	fmt.Println("Downloading Manifest.  Please Wait . . .")

	resp, errHttpGet := http.Get(manifestURL)
	defer resp.Body.Close()

	if errHttpGet != nil {
		fmt.Println("Failed to Download GoCore manifest:  " + errHttpGet.Error())
		return
	}

	n, errCopyOut := io.Copy(out, resp.Body)

	if errCopyOut != nil {
		fmt.Println("Failed to Output to manifest:  " + errCopyOut.Error())
		return
	}
	out.Close()

	fmt.Println("Downloaded GoCore manifest successfully:  " + extensions.PrintMegaBytes(n))
}

func readManifest() (coreManifest, error) {

	var manifest coreManifest

	jsonData, err := ioutil.ReadFile(manifestFileName)
	if err != nil {
		fmt.Println("Reading of " + manifestFileName + " failed:  " + err.Error())
		return manifest, err
	}

	errUnmarshal := json.Unmarshal(jsonData, &manifest)
	if errUnmarshal != nil {
		fmt.Println("Parsing / Unmarshaling of " + manifestFileName + " failed:  " + errUnmarshal.Error())
		return manifest, errUnmarshal
	}

	return manifest, nil
}
