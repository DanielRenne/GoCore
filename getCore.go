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
)

const gitRepoName = "GoCore-master"
const manifestFileName = "downloadedManifest.json"
const manifestURL = "https://raw.githubusercontent.com/DanielRenne/GoCore/master/coreManifest.json"

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

	downloadManifest()

	manifest, errManifest := readManifest()
	if errManifest != nil {
		fmt.Println("Failed to parse manifest file:  " + errManifest.Error())
		return
	}

	v := cleanGoCore(manifest)
	fmt.Println("")
	cleanupFiles()

	downloadRelease(v.ReleaseURL, v.ReleaseFileName+".zip")

	fmt.Println("Unzipping goCoreMaster.zip . . .")

	excludedFiles := []string{v.ReleaseFileName + "/webConfig.json"}
	excludedFiles = append(excludedFiles, v.ReleaseFileName+"/src/core/app/app.go")
	excludedFiles = append(excludedFiles, v.ReleaseFileName+"/keys/cert.pem")
	excludedFiles = append(excludedFiles, v.ReleaseFileName+"/keys/key.pem")

	errUnzip := zip.Unzip(v.ReleaseFileName+".zip", v.ReleaseFileName, excludedFiles)

	if errUnzip != nil {
		fmt.Println("Failed to Unzip goCoreMaster.zip:  " + errUnzip.Error())
		return
	}

	fmt.Println("Unzipped GoCore Master successfully.")

	errRemoveRepoZip := os.Remove(v.ReleaseFileName + ".zip")

	if errRemoveRepoZip != nil {
		fmt.Println("Failed to Remove GoCore Repo Zip File:  " + errRemoveRepoZip.Error())
		return
	}

	fmt.Println("\nMoving Files . . .")

	//Copy the Files then Remove the Directory
	extensions.CopyFolder(v.ReleaseFileName+"/"+v.ReleaseFileName, "")

	fmt.Println("Moved Files Successfully.")
	fmt.Println("\nStarted Cleaning Files.")
	errRemoveDecompressedFiles := extensions.RemoveDirectory(v.ReleaseFileName)

	if errRemoveDecompressedFiles != nil {
		fmt.Println("Failed to Remove GoCore Decompressed Files:  " + errRemoveDecompressedFiles.Error())
		return
	}

	fmt.Println("Cleaned up Repo Files Successfully.")
}

func downloadRelease(url string, fileName string) {

	fmt.Println("Downloading GoCore.  Please Wait . . .")

	out, errCreateFile := os.Create(fileName)

	if errCreateFile != nil {
		fmt.Println("Failed to create file handle:  " + errCreateFile.Error())
		return
	}

	resp, errHttpGet := http.Get(url)
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

func cleanupFiles() {
	err := os.Remove(manifestFileName)

	if err != nil {
		fmt.Println("Failed to clean up " + manifestFileName + ":  " + err.Error())
	}

}

func cleanGoCore(manifest coreManifest) version {

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

	for _, directory := range v.GoDirectories {
		fmt.Println("Removing Directory:  " + directory)
		extensions.RemoveDirectory(directory)
	}

	return v

}
