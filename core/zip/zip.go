package zip

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Unzip(archive, target string, excludedFiles []string) error {
	reader, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}

	defer reader.Close()

	if err := os.MkdirAll(target, 0777); err != nil {
		return err
	}

	for _, file := range reader.File {

		excludeFile := false
		for _, excludedFile := range excludedFiles {
			if excludedFile == file.Name {
				excludeFile = true
				break
			}
		}

		if excludeFile == true {
			fmt.Println("Excluding File:  " + file.Name)
			continue
		}

		path := filepath.Join(target, file.Name)
		if file.FileInfo().IsDir() {
			fmt.Println("Creating Directory:  " + file.Name)
			os.MkdirAll(path, file.Mode())
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			if fileReader != nil {
                fileReader.Close()
            }
			return err
		}

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			if fileReader != nil {
                fileReader.Close()
            }
			return err
		}

		fmt.Println("Extracting File:  " + file.Name)
        if _, err := io.Copy(targetFile, fileReader); err != nil {
			fileReader.Close()
			targetFile.Close()
			return err
		}
		fileReader.Close()
		targetFile.Close()
	}

	return nil
}

func Zipit(source, target string) error {
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	info, err := os.Stat(source)
	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})

	return err
}
