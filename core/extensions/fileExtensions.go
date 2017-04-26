package extensions

import (
	"archive/tar"
	"compress/gzip"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

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

func ReadFileAndParse(path string, v interface{}) (err error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &v)
	return
}

func ParseAndWriteFile(path string, v interface{}, perm os.FileMode) (err error) {
	data, err := json.Marshal(v)
	if err != nil {
		return
	}

	err = WriteToFile(string(data), path, perm)
	return
}

func WriteToFile(value string, path string, perm os.FileMode) error {
	err := ioutil.WriteFile(path, []byte(value), perm)
	if err != nil {
		log.Println("Error writing file " + path + ":  " + err.Error())
		return err
	}
	return nil
}

func DoesFileExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func DoesFileNotExist(path string) bool {
	if _, err := os.Stat(path); err == nil {
		// path/to/whatever exists
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

func UnGzipfunc(source, target string) error {
	reader, err := os.Open(source)
	if err != nil {
		return err
	}
	defer reader.Close()

	archive, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}
	defer archive.Close()

	target = filepath.Join(target, archive.Name)
	writer, err := os.Create(target)
	if err != nil {
		return err
	}
	defer writer.Close()

	_, err = io.Copy(writer, archive)
	return err
}

func GetFileSize(path string) (size int64, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	fi, err := file.Stat()
	if err != nil {
		return
	}
	size = fi.Size()
	return
}

func UnTar(tarball, target string) error {
	reader, err := os.Open(tarball)
	if err != nil {
		return err
	}
	defer reader.Close()
	tarReader := tar.NewReader(reader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		path := filepath.Join(target, header.Name)
		info := header.FileInfo()
		if info.IsDir() {
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				return err
			}
			continue
		}

		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(file, tarReader)
		if err != nil {
			return err
		}
	}
	return err
}
