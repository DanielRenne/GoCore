package httpExtensions

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func DownloadFromUrl(url string, path string) error {
	fmt.Println("Downloading", url, "to", path)

	// TODO: check file existence first with io.IsExist
	output, err := os.Create(path)
	if err != nil {
		fmt.Println("Error while creating", path, "-", err)
		return err
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return err
	}
	defer response.Body.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return err
	}

	fmt.Println(n, "bytes downloaded.")

	return nil
}
