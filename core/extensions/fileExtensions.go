package extensions

import (
	"archive/tar"
	"compress/gzip"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/DanielRenne/GoCore/core/utils"
)

// ByOldestFile is a type for sorting files by oldest first
type ByOldestFile []os.FileInfo

func (a ByOldestFile) Len() int           { return len(a) }
func (a ByOldestFile) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByOldestFile) Less(i, j int) bool { return a[i].ModTime().Unix() < a[j].ModTime().Unix() }

// GetAllFilesSortedBy returns all files in a directory sorted
func GetAllFilesSortedBy(path string, fileSearch string) (files []os.FileInfo, err error) {
	files, err = GetAllFilesWithSearch(path, fileSearch)
	if err == nil {
		sort.Sort(ByOldestFile(files))
	}
	return files, err
}

// GetAllFiles returns all files in a directory
func GetAllFiles(path string) (files []os.FileInfo, err error) {
	return GetAllFilesWithSearch(path, "")
}

// GetAllFilesWithSearch returns all files in a directory with a search string
func GetAllFilesWithSearch(path string, fileSearch string) ([]fs.FileInfo, error) {
	files := []fs.FileInfo{}
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.Contains(entry.Name(), fileSearch) {
			// Convert fs.DirEntry to fs.FileInfo
			info, err := entry.Info()
			if err != nil {
				return nil, err
			}
			files = append(files, info)
		}
	}

	return files, nil
}

// GetAllFolders returns all folders in a directory
func GetAllFolders(path string) (files []os.FileInfo, err error) {
	return GetAllFoldersWithSearch(path, "")
}

// GetAllFoldersWithSearch returns all folders in a directory with a search string
func GetAllFoldersWithSearch(path string, fileSearch string) ([]fs.FileInfo, error) {
	// Initialize a slice to store the results
	files := []fs.FileInfo{}

	// Read all entries in the directory
	filesAll, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	// Loop through each directory entry
	for _, file := range filesAll {
		if file.IsDir() { // Check if the entry is a directory
			if fileSearch == "" || strings.Contains(file.Name(), fileSearch) {
				// Convert DirEntry to FileInfo
				info, err := file.Info()
				if err != nil {
					return nil, err
				}
				// Append the directory's FileInfo to the slice
				files = append(files, info)
			}
		}
	}

	return files, nil
}

// GetAllFilesDeepWithSearch recursively returns all files in a directory with a search string
func GetAllFilesDeepWithSearch(path string, fileSearch string) ([]fs.FileInfo, error) {
	// Initialize a slice to store the results
	files := []fs.FileInfo{}

	// Read all entries in the directory
	filesAll, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	// Loop through each directory entry
	for _, file := range filesAll {
		if !file.IsDir() { // Check if the entry is not a directory
			if fileSearch == "" || strings.Contains(file.Name(), fileSearch) {
				// Convert DirEntry to FileInfo
				info, err := file.Info()
				if err != nil {
					return nil, err
				}
				// Append the file's FileInfo to the slice
				files = append(files, info)
			}
		} else { // If the entry is a directory, recursively search it
			subPath := filepath.Join(path, file.Name())
			subFiles, err := GetAllFilesDeepWithSearch(subPath, fileSearch)
			if err != nil {
				return nil, err
			}
			files = append(files, subFiles...)
		}
	}

	return files, nil
}

// GetAllFilesSearchWithPath returns all files in a directory with a search string
func GetAllFilesSearchWithPath(path string, fileSearch string) (files []FilePath, err error) {
	files = make([]FilePath, 0)
	filesAll, err := os.ReadDir(path)
	if err == nil {
		for _, file := range filesAll {
			if !file.IsDir() {
				if fileSearch == "" || strings.Contains(file.Name(), fileSearch) {
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

// GetAllDirs returns all directories in a directory
func GetAllDirs(path string) (files []os.FileInfo, err error) {
	return GetAllDirWithExclude(path, "")
}

// GetAllDirWithExclude returns all directories in a directory excluding the exclude string
func GetAllDirWithExclude(path string, except string) ([]fs.FileInfo, error) {
	// Initialize a slice to store the results
	files := []fs.FileInfo{}

	// Read all entries in the directory
	filesAll, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	// Loop through each directory entry
	for _, file := range filesAll {
		if file.IsDir() { // Check if the entry is a directory
			if except == "" || !strings.Contains(file.Name(), except) {
				// Convert DirEntry to FileInfo
				info, err := file.Info()
				if err != nil {
					return nil, err
				}
				// Append the directory's FileInfo to the slice
				files = append(files, info)
			}
		}
	}

	return files, nil
}

// DirSize returns the size of a directory
func DirSize(path string) (int64, error) {
	return DirSizeWithSearch(path, "")
}

// DirSizeWithSearch returns the size of a directory with a search string
func DirSizeWithSearch(path string, fileSearch string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			if fileSearch == "" || strings.Contains(info.Name(), fileSearch) {
				size += info.Size()
			}
		}
		return err
	})
	return size, err
}

// RemoveDirectoryShell removes a directory using the shell
func RemoveDirectoryShell(dir string) (err error) {
	if dir == "/" {
		return errors.New("cannot remove root directory")
	}
	cmd := exec.Command("rm", "-rf", dir)
	err = cmd.Run()
	return
}

// RemoveDirectory removes a directory
func RemoveDirectory(dir string) error {
	if dir == "/" {
		return errors.New("cannot remove root directory")
	}
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

// CopyFolder copies a folder to another folder
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
	defer directory.Close()
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

// CopyFile copies a file to another file
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
			if err != nil {
				return err
			}
		}

	}

	return
}

// WriteAndGoFormat writes a file and formats it with go
func WriteAndGoFormat(value string, path string) error {
	return WriteAndGoFmt(value, path, false, 0777)
}

// WriteAndGoFmt writes a file and formats it with go
func WriteAndGoFmt(value string, path string, quiet bool, perm os.FileMode) error {

	err := os.WriteFile(path, []byte(value), perm)
	if err != nil {
		if !quiet {
			log.Println("Error writing file " + path + ":  " + err.Error())
		}
		return err
	}

	cmd := exec.Command("gofmt", "-w", path)
	err = cmd.Start()
	if err != nil {
		if !quiet {
			log.Println("Failed to gofmt on file " + path + ":  " + err.Error())
		}
		return err
	}
	if !quiet {
		log.Println("Saved file " + path + " successfully.")
	}
	return nil
}

// ReadFile reads a file
func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// ReadFileAndParse reads a file and parses it with json.Unmarshal
func ReadFileAndParse(path string, v interface{}) (err error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &v)
	return
}

// ParseAndWriteFile parses a file and writes it with json.Marshal
func ParseAndWriteFile(path string, v interface{}, perm os.FileMode) (err error) {
	data, err := json.Marshal(v)
	if err != nil {
		return
	}

	err = WriteToFile(string(data), path, perm)
	return
}

// GetAllFilesRecursively recurses a directory and returns the path and file names of all the files.
func GetAllFilesRecursively(fileDir string) (files []string, err error) {
	var satisfy []string
	return recurseFiles(fileDir, satisfy)
}

// GetAllFilesRecursivelyByExtension recurses a directory and returns the path and file names of all the files.  Pass []string of fileExtensions as the second parameter (no need for period and cannot have multiple extensions like tar.gz)
func GetAllFilesRecursivelyByExtension(fileDir string, fileExtensionFilter []string) (files []string, err error) {
	return recurseFiles(fileDir, fileExtensionFilter)
}

func recurseFiles(fileDir string, fileExtensionFilter []string) (files []string, err error) {
	path := fileDir

	if !DoesFileExist(path) {
		err = errors.New(fileDir + " does not exist")
		return
	}
	var newFilterExtensions []string
	for i := range fileExtensionFilter {
		newFilterExtensions = append(newFilterExtensions, strings.ToUpper(strings.ReplaceAll(fileExtensionFilter[i], ".", "")))
	}
	err = filepath.Walk(path, func(path string, f os.FileInfo, errWalk error) (err error) {

		if errWalk != nil {
			err = errWalk
			return
		}

		if !f.IsDir() {
			addFile := true
			if len(fileExtensionFilter) > 0 {
				pieces := strings.Split(path, ".")
				ext := strings.ToUpper(pieces[len(pieces)-1:][0])
				if !utils.InArray(ext, newFilterExtensions) {
					addFile = false
				}
			}
			if addFile {
				files = append(files, path)
			}
		}

		return
	})

	return
}

// InvokeMethodOnAllFilesRecursively recurses a directory and returns the path and file names of all the files.  The second parameter is a callback functon which will be called in go routines for each file where you will read the os.File and do something with it
func InvokeMethodOnAllFilesRecursively(fileDir string, function func(*os.File) error) (files []string, fileErrors map[string]error, err error) {
	var satisfy []string
	return recurseFilesAsync(fileDir, satisfy, function)
}

// InvokeMethodOnAllFilesRecursivelyByExtension recurses a directory and returns the path and file names of all the files.  Pass []string of fileExtensions as the second parameter (no need for period and cannot have multiple extensions like tar.gz). The third parameter is a callback functon which will be called in go routines for each file where you will read the os.File and do something with it
func InvokeMethodOnAllFilesRecursivelyByExtension(fileDir string, fileExtensionFilter []string, function func(*os.File) error) (files []string, fileErrors map[string]error, err error) {
	return recurseFilesAsync(fileDir, fileExtensionFilter, function)
}

func recurseFilesAsync(fileDir string, fileExtensionFilter []string, function func(*os.File) error) (files []string, fileErrors map[string]error, err error) {
	fileErrors = make(map[string]error, 0)
	type itemSync struct {
		sync.Mutex
		Items []string
	}
	var syncedItems itemSync
	path := fileDir
	var wg sync.WaitGroup
	if !DoesFileExist(path) {
		err = errors.New(fileDir + " does not exist")
		return
	}
	var newFilterExtensions []string
	for i := range fileExtensionFilter {
		newFilterExtensions = append(newFilterExtensions, strings.ToUpper(strings.ReplaceAll(fileExtensionFilter[i], ".", "")))
	}
	totalFiles := 0
	err = filepath.Walk(path, func(path string, f os.FileInfo, errWalk error) (err error) {

		if errWalk != nil {
			err = errWalk
			return
		}

		if !f.IsDir() {
			addFile := true
			if len(fileExtensionFilter) > 0 {
				pieces := strings.Split(path, ".")
				ext := strings.ToUpper(pieces[len(pieces)-1:][0])
				if !utils.InArray(ext, newFilterExtensions) {
					addFile = false
				}
			}
			if addFile {
				wg.Add(1)
				totalFiles++
				go func() {
					defer func() {
						if r := recover(); r != nil {
							return
						}
					}()

					defer wg.Done()
					f, err := os.Open(path)
					if err != nil {
						fileErrors[path] = errors.New("File Open Error:" + err.Error())
						return
					}
					defer f.Close()
					err = function(f)
					if err != nil {
						fileErrors[path] = errors.New("Callback function error:" + err.Error())
						return
					}
					syncedItems.Lock()
					syncedItems.Items = append(syncedItems.Items, path)
					syncedItems.Unlock()

				}()
				files = append(files, path)
			}
		}

		return
	})
	wg.Wait()
	files = syncedItems.Items

	return
}

// WriteToFile writes a file
func WriteToFile(value string, path string, perm os.FileMode) error {
	err := os.WriteFile(path, []byte(value), perm)
	if err != nil {
		log.Println("Error writing file " + path + ":  " + err.Error())
		return err
	}
	return nil
}

// Write to a file
func Write(value string, path string) error {
	return WriteToFile(value, path, UNIX_COMMON_FILE)
}

// MkDir creates a directory
func MkDir(path string) error {
	return os.MkdirAll(path, UNIX_COMMON_DIR)
}

// MkDirRWAll creates a directory with 0777 permissions
func MkDirRWAll(path string) error {
	return os.MkdirAll(path, UNIX_DIR_RWALL)
}

// DoesFileExist checks if a file exists
func DoesFileExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// DoesFileNotExist checks if a file does not exist
func DoesFileNotExist(path string) bool {
	if _, err := os.Stat(path); err == nil {
		// path/to/whatever exists
		return false
	}
	return true
}

// MD returns the md5 hash of a string
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

// UnGzipfunc ungzips a file
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

// Gzipfunc gzips a file
func Gzipfunc(source string, target string) (err error) {
	var data []byte
	data, err = ReadFile(source)
	if err != nil {
		return
	}
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
	defer file.Close()
	fi, err := file.Stat()
	if err != nil {
		return
	}
	size = fi.Size()
	return
}

// UnTar unzips a tar file
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
