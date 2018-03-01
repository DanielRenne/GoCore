package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func main() {
	log.Println(GetGoPath())
	err := CopyFolder(GetGoPath()+"/src/github.com/DanielRenne/goCoreAppTemplate/vendorPackages", GetGoPath()+"/src")

	if err != nil {
		log.Println("Failed to Copy vendor packages.  \n" + err.Error())
		return
	}

	log.Println("Successfully installed go core app dependencies.")
}

func RunCommand(cmd string) error {

	var subProcess = exec.Command(cmd)

	if runtime.GOOS != "windows" {
		subProcess = exec.Command("/bin/sh", "-c", cmd)
	}

	fmt.Println("Attempting to run command:  " + cmd + "\n")

	subProcess.Stdout = os.Stdout
	subProcess.Stderr = os.Stdout

	err := subProcess.Run()
	if err != nil {
		fmt.Println("Failed to start command:  " + cmd)
		fmt.Println("Error:  " + err.Error() + "\n")
		return err
	}

	return nil
}

func GetGoPath() string {
	return os.Getenv("GOPATH")
}

func RemoveDirectory(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	os.Remove(dir)
	return nil
}

func CopyFolder(source string, dest string) (err error) {

	sourceinfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	if dest != "" {
		fmt.Println(dest)

		err = os.MkdirAll(dest, sourceinfo.Mode())
		if err != nil {
			fmt.Println("Failed to create Directory: " + dest + ",  " + err.Error())
			return err
		}
	}

	directory, _ := os.Open(source)

	objects, err := directory.Readdir(-1)

	for _, obj := range objects {

		sourcefilepointer := source + "/" + obj.Name()

		dstFilePath := dest + "/" + obj.Name()
		if dest == "" {
			dstFilePath = obj.Name()
		}

		destinationfilepointer := dstFilePath

		if obj.IsDir() {
			err = CopyFolder(sourcefilepointer, destinationfilepointer)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			err = CopyFile(sourcefilepointer, destinationfilepointer)
			if err != nil {
				fmt.Println(err)
			}
		}

	}
	return
}

func CopyFile(source string, dest string) (err error) {
	sourcefile, err := os.Open(source)
	if err != nil {
		return err
	}

	defer sourcefile.Close()

	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer destfile.Close()

	_, err = io.Copy(destfile, sourcefile)
	if err == nil {
		sourceinfo, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dest, sourceinfo.Mode())
		}

	}

	return
}

func WriteAndGoFormat(value string, path string) error {

	err := ioutil.WriteFile(path, []byte(value), 0777)
	if err != nil {
		log.Println("Error writing file " + path + ":  " + err.Error())
		return err
	}

	cmd := exec.Command("gofmt", "-w", path)
	err = cmd.Start()
	if err != nil {
		log.Println("Failed to gofmt on file " + path + ":  " + err.Error())
		return err
	}

	log.Println("Saved file " + path + " successfully.")
	return nil
}

func ReadFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func DoesFileExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func MD5(path string) (string, error) {
	hasher := md5.New()

	fileData, err := ReadFile(path)
	if err != nil {
		return "", err
	}

	hasher.Write([]byte(fileData))
	val := hex.EncodeToString(hasher.Sum(nil))

	return val, nil
}
