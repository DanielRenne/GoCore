// Package path provides basic path functions for OS specific path folder names and a function to get the base path of where the binary resides
package path

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

// GetBinaryPath will return the location of the binary or the project in go run mode.
// Note, if you compiled a go Binary and exported it into an OSX App Package folder, the base path would be the root of where the app lives
// If you are running a go binary from the command line using go run, the base path returned of where the binary is located so long as you are in the directory of the main and you pass $(cwd) to the first parameter of the go run
// Else the path when compiled is easily found by using filepath.Dir(os.Executable())
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
		if strings.Contains(exPath, "Contents/MacOS") {
			dirs := strings.Split(exPath, "/")
			dirs = dirs[:len(dirs)-3]
			path = strings.Join(dirs, "/")
		} else {
			path = exPath
		}
	}
	return
}
