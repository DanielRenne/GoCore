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

func main() {

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

	fmt.Println("Unzipping goCoreMaster.zip . . .")

	errUnzip := zip.Unzip(fileName, "goCoreMaster")

	if errUnzip != nil {
		fmt.Println("Failed to Unzip goCoreMaster.zip:  " + errUnzip.Error())
		return
	}

	fmt.Println("Unzipped & Installed GoCore Master successfully.")
}
