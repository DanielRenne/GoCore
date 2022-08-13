package path

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

//GetBinaryPath will return the location of the binary or the project in go run mode.
func GetBinaryPath() (path string) {
	ex, err := os.Executable()
	if err != nil {
		log.Println("path.GetBinaryPath() failed to get Executable: " + err.Error())
		return
	}
	exPath := filepath.Dir(ex)
	if strings.Contains(exPath, "go-build") { //This means we are running in go run and need to use the goPath
		if len(os.Args) >= 2 {
			path = os.Args[1]
		} else {
			log.Println("Using GoCore in go run, you must pass go run yourMain.go $(pwd) so it can get your path of your project")
			os.Exit(1)
		}
	} else {
		if strings.Index(exPath, "Contents/MacOS") != -1 {
			dirs := strings.Split(exPath, "/")
			dirs = dirs[:len(dirs)-3]
			path = strings.Join(dirs, "/")
		} else {
			path = exPath
		}
	}
	return
}
