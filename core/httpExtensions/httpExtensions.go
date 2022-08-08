package httpExtensions

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/DanielRenne/GoCore/core/extensions"
)

func GetJSONRequest(urlRequest string, x interface{}) (err error) {
	var myClient = &http.Client{Timeout: 25 * time.Second}
	var r *http.Response
	r, err = myClient.Get(urlRequest)
	if err != nil {
		return
	}
	defer r.Body.Close()

	if r.StatusCode <= 400 {
		err = json.NewDecoder(r.Body).Decode(x)
	} else {
		err = errors.New("Non 200 status code: " + extensions.IntToString(r.StatusCode))
	}
	return
}

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
