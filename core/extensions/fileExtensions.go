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
	"sort"
	"strings"
)

type ByOldestFile []os.FileInfo

func (a ByOldestFile) Len() int           { return len(a) }
func (a ByOldestFile) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByOldestFile) Less(i, j int) bool { return a[i].ModTime().Unix() < a[j].ModTime().Unix() }

func GetAllFilesSortedBy(path string, fileSearch string) (files []os.FileInfo, err error) {
	files, err = GetAllFilesWithSearch(path, fileSearch)
	if err == nil {
		sort.Sort(ByOldestFile(files))
	}
	return files, err
}

func GetAllFiles(path string) (files []os.FileInfo, err error) {
	return GetAllFilesWithSearch(path, "")
}

func GetAllFilesWithSearch(path string, fileSearch string) (files []os.FileInfo, err error) {
	files = make([]os.FileInfo, 0)
	filesAll, err := ioutil.ReadDir(path)
	if err == nil {
		for _, file := range filesAll {
			if !file.IsDir() {
				if fileSearch == "" || strings.Index(file.Name(), fileSearch) != -1 {
					files = append(files, file)
				}
			}
		}
	}
	return files, err
}

func GetAllFolders(path string) (files []os.FileInfo, err error) {
	return GetAllFoldersWithSearch(path, "")
}

func GetAllFoldersWithSearch(path string, fileSearch string) (files []os.FileInfo, err error) {
	files = make([]os.FileInfo, 0)
	filesAll, err := ioutil.ReadDir(path)
	if err == nil {
		for _, file := range filesAll {
			if file.IsDir() {
				if fileSearch == "" || strings.Index(file.Name(), fileSearch) != -1 {
					files = append(files, file)
				}
			}
		}
	}
	return files, err
}

func GetAllFilesDeepWithSearch(path string, fileSearch string) (files []os.FileInfo, err error) {
	files = make([]os.FileInfo, 0)
	filesAll, err := ioutil.ReadDir(path)
	if err == nil {
		for _, file := range filesAll {
			if !file.IsDir() {
				if fileSearch == "" || strings.Index(file.Name(), fileSearch) != -1 {
					files = append(files, file)
				}
			} else {
				subFiles, err := GetAllFilesDeepWithSearch(path+string(os.PathSeparator)+file.Name(), fileSearch)
				files = append(files, subFiles...)
				if err != nil {
					return files, err
				}
			}
		}
	}
	return files, err
}

func GetAllFilesSearchWithPath(path string, fileSearch string) (files []FilePath, err error) {
	files = make([]FilePath, 0)
	filesAll, err := ioutil.ReadDir(path)
	if err == nil {
		for _, file := range filesAll {
			if !file.IsDir() {
				if fileSearch == "" || strings.Index(file.Name(), fileSearch) != -1 {
					split := strings.Split(file.Name(), ".")
					fileType := ""
					if len(split) > 1 {
						fileType = split[len(split)-1]
					}

					fp := FilePath{
						Name: file.Name(),
						Path: path,
						Type: fileType,
					}
					files = append(files, fp)
				}
			} else {
				subFiles, err := GetAllFilesSearchWithPath(path+string(os.PathSeparator)+file.Name(), fileSearch)
				files = append(files, subFiles...)
				if err != nil {
					return files, err
				}
			}
		}
	}
	return files, err
}

func GetAllDirs(path string) (files []os.FileInfo, err error) {
	return GetAllDirWithExclude(path, "")
}

func GetAllDirWithExclude(path string, except string) (files []os.FileInfo, err error) {
	files = make([]os.FileInfo, 0)
	filesAll, err := ioutil.ReadDir(path)
	if err == nil {
		for _, file := range filesAll {
			if file.IsDir() {
				if except == "" || strings.Index(file.Name(), except) == -1 {
					files = append(files, file)
				}
			}
		}
	}
	return files, err
}

func DirSize(path string) (int64, error) {
	return DirSizeWithSearch(path, "")
}

func DirSizeWithSearch(path string, fileSearch string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			if fileSearch == "" || strings.Index(info.Name(), fileSearch) != -1 {
				size += info.Size()
			}
		}
		return err
	})
	return size, err
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

func Gzipfunc(source string, target string) (err error) {

	data, err := ReadFile(source)

	// Open a file for writing.
	f, _ := os.Create(target)

	// Create gzip writer.
	w := gzip.NewWriter(f)

	// Write bytes in compressed form to the file.
	_, err = w.Write(data)

	// Close the file.
	w.Close()
	return
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

// Tar takes a source and variable writers and walks 'source' writing each file
// found to the tar writer; the purpose for accepting multiple writers is to allow
// for multiple outputs (for example a file, or md5 hash)
func Tar(source string, target string) error {

	tarfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer tarfile.Close()

	tarball := tar.NewWriter(tarfile)
	defer tarball.Close()

	info, err := os.Stat(source)
	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	return filepath.Walk(source,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			header, err := tar.FileInfoHeader(info, info.Name())
			if err != nil {
				return err
			}

			if baseDir != "" {
				header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
			}

			if err := tarball.WriteHeader(header); err != nil {
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
			_, err = io.Copy(tarball, file)
			return err
		})
}
