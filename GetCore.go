package main

import (
	"core/extensions"
	"core/zip"
	"fmt"
	"io"
	"net/http"
	"os"
)

const goCoreURL = "https://github.com/DanielRenne/GoCore/archive/master.zip"
const fileName = "goCoreMaster.zip"
const gitRepoName = "GoCore-master"

func main() {

	out, errCreateFile := os.Create(fileName)

	if errCreateFile != nil {
		fmt.Println("Failed to create file handle:  " + errCreateFile.Error())
		return
	}

	fmt.Println("Downloading GoCore.  Please Wait . . .")

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

	fmt.Println("Unzipping goCoreMaster.zip . . .")

	excludedFiles := []string{gitRepoName + "/webConfig.json"}
	excludedFiles = append(excludedFiles, gitRepoName+"/src/core/app/app.go")

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

	//Copy the Files then Remove the Directory
	extensions.CopyFolder("goCoreMaster/"+gitRepoName, "")

	fmt.Println("Moved Files Successfully.")
	fmt.Println("Started Cleaning Files.")
	errRemoveDecompressedFiles := extensions.RemoveDirectory("goCoreMaster")

	if errRemoveDecompressedFiles != nil {
		fmt.Println("Failed to Remove GoCore Decompressed Files:  " + errRemoveDecompressedFiles.Error())
		return
	}

	fmt.Println("Cleaned up Repo Files Successfully.")
}
