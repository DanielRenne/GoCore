package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/DanielRenne/GoCore/core/extensions"
	"github.com/DanielRenne/GoCore/core/zip"
)

const gitRepoName = "GoCoreDep-master"
const manifestFileName = "downloadedManifest.json"
const manifestURL = "https://raw.githubusercontent.com/DanielRenne/GoCoreDep/master/coreManifest.json"

type version struct {
	Version         string   `json:"version"`
	ReleaseURL      string   `json:"releaseURL"`
	ReleaseFileName string   `json:"releaseFileName"`
	GoDirectories   []string `json:"goDirectories"`
}

type coreManifest struct {
	Versions []version `json:"versions"`
}

func main() {

	appGenerationFileName := "appGeneration"
	downloadRelease("https://github.com/davidrenne/GoCoreAppTemplate/archive/master.zip", appGenerationFileName+".zip")

	fmt.Println("Unzipping appGenerationFileName.zip . . .")

	excludedFiles := []string{}
	errUnzip := zip.Unzip(appGenerationFileName+".zip", appGenerationFileName, excludedFiles)

	if errUnzip != nil {
		fmt.Println("Failed to Unzip appGenerationFileName.zip:  " + errUnzip.Error())
		return
	}

	fmt.Println("Unzipped appGeneration successfully.")

	errRemoveRepoZip := os.Remove(appGenerationFileName + ".zip")

	if errRemoveRepoZip != nil {
		fmt.Println("Failed to Remove appGeneration zip File:  " + errRemoveRepoZip.Error())
		return
	}

	os.RemoveAll("/tmp/tools/appFiles")
	os.MkdirAll("/tmp/tools/appFiles", 0777)

	//Copy the Files then Remove the Directory
	fmt.Println(appGenerationFileName + "/master")
	fmt.Println("/tmp/tools/appFiles")
	extensions.CopyFolder(appGenerationFileName+"/GoCoreAppTemplate-master", "/tmp/tools/appFiles")

	fmt.Println("Moved Files Successfully.")
	fmt.Println("\nStarted Cleaning Files.")
	errRemoveDecompressedFiles := extensions.RemoveDirectory(appGenerationFileName)

	if errRemoveDecompressedFiles != nil {
		fmt.Println("Failed to Remove appGeneration Files:  " + errRemoveDecompressedFiles.Error())
		return
	}
}

func downloadRelease(url string, fileName string) {

	fmt.Println("Downloading GoCore Dependencies. (" + url + ")  Please Wait . . .")

	out, errCreateFile := os.Create(fileName)

	if errCreateFile != nil {
		fmt.Println("Failed to create file handle:  " + errCreateFile.Error())
		return
	}

	resp, errHttpGet := http.Get(url)

	if errHttpGet != nil {
		fmt.Println("Failed to Download GoCoreDep master zip file:  " + errHttpGet.Error())
		return
	}
	defer resp.Body.Close()

	n, errCopyOut := io.Copy(out, resp.Body)

	if errCopyOut != nil {
		fmt.Println("Failed to Output to goCoreMaster.zip:  " + errCopyOut.Error())
		return
	}
	out.Close()

	fmt.Println("Downloaded GoCoreDep Master successfully:  " + extensions.PrintMegaBytes(n))
}

func downloadManifest() {

	out, errCreateFile := os.Create(manifestFileName)

	if errCreateFile != nil {
		fmt.Println("Failed to create file handle:  " + errCreateFile.Error())
		return
	}

	fmt.Println("Downloading Manifest.  Please Wait . . .")

	resp, errHttpGet := http.Get(manifestURL)

	if errHttpGet != nil {
		fmt.Println("Failed to Download GoCoreDep manifest:  " + errHttpGet.Error())
		return
	}
	defer resp.Body.Close()

	n, errCopyOut := io.Copy(out, resp.Body)

	if errCopyOut != nil {
		fmt.Println("Failed to Output to manifest:  " + errCopyOut.Error())
		return
	}
	out.Close()

	fmt.Println("Downloaded GoCoreDep manifest successfully:  " + extensions.PrintMegaBytes(n))
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

func cleanupFiles() {
	err := os.Remove(manifestFileName)

	if err != nil {
		fmt.Println("Failed to clean up " + manifestFileName + ":  " + err.Error())
	}

}

func cleanGoCore(v version) {

	for _, directory := range v.GoDirectories {
		fmt.Println("Removing Directory:  " + directory)
		os.RemoveAll(directory)
	}
}

func getVersion(manifest coreManifest) version {
	v := manifest.Versions[len(manifest.Versions)-1]

	if len(os.Args) > 1 {
		argVersion := os.Args[1]

		for _, vn := range manifest.Versions {
			if vn.Version == argVersion {
				v = vn
				break
			}
		}
	}

	return v
}

func excludeFiles(excludedFiles []string, v version) {

	_, err := os.Stat("webConfig.json")
	if err == nil {
		excludedFiles = append(excludedFiles, v.ReleaseFileName+"/webConfig.json")
	}

	_, err1 := os.Stat("src/github.com/DanielRenne/GoCore/core/app/app.go")
	if err1 == nil {
		excludedFiles = append(excludedFiles, v.ReleaseFileName+"/src/github.com/DanielRenne/GoCore/core/app/app.go")
	}

	_, err2 := os.Stat("/keys/cert.pem")
	if err2 == nil {
		excludedFiles = append(excludedFiles, v.ReleaseFileName+"/keys/cert.pem")
	}

	_, err3 := os.Stat("/keys/key.pem")
	if err3 == nil {
		excludedFiles = append(excludedFiles, v.ReleaseFileName+"/keys/key.pem")
	}
}
