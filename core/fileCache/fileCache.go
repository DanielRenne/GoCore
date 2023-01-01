// Package fileCache will provide simple file content caching tools for in-Memory access to files.
// It uses golang/groupcache to cache your data into memory on multiple HTTP Pool servers.
package fileCache

import (
	"context"
	"log"
	"os"
	"sync"

	"encoding/json"

	"github.com/DanielRenne/GoCore/core/extensions"
	"github.com/DanielRenne/GoCore/core/path"
	"github.com/DanielRenne/GoCore/core/serverSettings"
	"github.com/DanielRenne/GoCore/core/utils"
	"github.com/golang/groupcache"
)

// CACHE_STORAGE_PATH is the base GoCore path where all caches are stored
var CACHE_STORAGE_PATH string

// CACHE_JOBS is the directory where jobs are stored
var CACHE_JOBS string

// CACHE_BOOTSTRAP_STORAGE_PATH is the directory where bootstrap caches are stored
var CACHE_BOOTSTRAP_STORAGE_PATH string

// CACHE_MANIFEST_STORAGE_PATH is the directory where bootstrap caches are stored
var CACHE_MANIFEST_STORAGE_PATH string

var hasInitialized bool

// SetGoCoreStoragePath set a directory with a trailing slash of where you want goCore to set and make directories for caching files needed to keep track of your application cron and bootstrap caches.
func SetGoCoreStoragePath(directory string) {
	CACHE_STORAGE_PATH = directory + "caches"
	CACHE_JOBS = directory + "jobs"
	CACHE_BOOTSTRAP_STORAGE_PATH = CACHE_STORAGE_PATH + path.PathSeparator + "bootstrap"
	CACHE_MANIFEST_STORAGE_PATH = CACHE_STORAGE_PATH + path.PathSeparator + "manifests"
	os.MkdirAll(CACHE_BOOTSTRAP_STORAGE_PATH, 0777)
	os.MkdirAll(CACHE_JOBS, 0777)
	os.MkdirAll(CACHE_MANIFEST_STORAGE_PATH, 0777)
}

type model struct {
	sync.RWMutex
	BootstrapCache map[string][]string
}
type job struct {
	sync.RWMutex
	Jobs map[string]bool
}

type byteManifest struct {
	sync.RWMutex
	Cache map[string]map[string]int
}

var allGroupCacheDomains []string
var peers *groupcache.HTTPPool
var htmlFileCache *groupcache.Group
var stringCache *groupcache.Group

// Model is the in memory model for bootstrap caches
var Model model

// Jobs is the in memory model for jobs
var Jobs job

// ByteManifest is the in memory model for byte manifest caches
var ByteManifest byteManifest

// contains the temporary string cache used to cache large strings.
var tempStringCacheSynced = struct {
	sync.RWMutex
	cache map[string]string
}{cache: make(map[string]string)}

func init() {
	Model = model{
		BootstrapCache: make(map[string][]string, 0),
	}
	Jobs = job{
		Jobs: make(map[string]bool, 0),
	}
	ByteManifest = byteManifest{
		Cache: make(map[string]map[string]int, 0),
	}
}

// Init will initilize a groupCache if you pass a non empty string and create necessary folders for internal file caching of GoCore
func Init(groupCache string) {
	if !hasInitialized {
		if groupCache != "" {
			InitializeGroupCache(groupCache)
		}
		Initialize()
		hasInitialized = true
	}
}

// Initialize in main before any calls to this package are performed.  serverSettings package must be initialized before fileCache.
// Developers can call SetGoCoreStoragePath() with a path of their choice for storage of where bootstrap caches and jobs files (for one time cron jobs are stored
func Initialize() {
	if !hasInitialized {

		if !path.IsWindows && CACHE_STORAGE_PATH == "" {
			SetGoCoreStoragePath("/usr/local/goCore/")
		} else if path.IsWindows && CACHE_STORAGE_PATH == "" {
			SetGoCoreStoragePath("C:\\goCore\\")
		}

		if serverSettings.WebConfig.Application.Domain != "" {
			InitializeGroupCache(serverSettings.WebConfig.Application.Domain)
		}
		LoadJobsFile()
		hasInitialized = true
	}
}

// WriteJobCacheFile is exported internally to share between GoCore packages and should not be called directly by you.
func WriteJobCacheFile() (err error) {
	strjson, err := json.Marshal(Jobs.Jobs)
	if err != nil {
		return err
	}
	err = os.WriteFile(CACHE_JOBS+"/jobs.json", []byte(strjson), 0777)
	if err != nil {
		return err
	}
	return nil
}

// LoadJobsFile is exported internally to share between GoCore packages and should not be called directly by you.
func LoadJobsFile() (err error) {
	fname := CACHE_JOBS + "/jobs.json"
	if extensions.DoesFileExist(fname) {
		var size int64
		size, err = extensions.GetFileSize(fname)
		if err != nil {
			return
		}
		if size > 0 {
			var data map[string]bool
			jsonData, err := extensions.ReadFile(fname)
			if err != nil {
				log.Println("Cache failed to read for " + fname + " deleting file and starting fresh.")
				if extensions.DoesFileExist(fname) {
					err = os.Remove(fname)
					if err != nil {
						return err
					}
				}
				return err
			}
			err = json.Unmarshal(jsonData, &data)
			if err != nil {
				return err
			}
			Jobs.Lock()
			Jobs.Jobs = data
			Jobs.Unlock()
		}
	}
	return
}

// WriteBootstrapCacheFile is exported internally to share between GoCore packages and should not be called directly by you.
func WriteBootStrapCacheFile(key string) (err error) {
	Model.RLock()
	caches, ok := Model.BootstrapCache[key]
	Model.RUnlock()
	if ok {
		strjson, err := json.Marshal(caches)
		if err != nil {
			return err
		}
		err = os.WriteFile(CACHE_BOOTSTRAP_STORAGE_PATH+"/"+key+".json", []byte(strjson), 0777)
		if err != nil {
			return err
		}
	}
	return nil
}

// UpdateBootstrapMemoryCache is exported internally to share between GoCore packages and should not be called directly by you.
func UpdateBootStrapMemoryCache(key string, value string) {
	Model.RLock()
	_, ok := Model.BootstrapCache[key]
	Model.RUnlock()
	if !ok {
		Model.Lock()
		Model.BootstrapCache[key] = utils.Array(value)
		Model.Unlock()
	} else {
		if !DoesHashExistInCache(key, value) {
			Model.Lock()
			Model.BootstrapCache[key] = append(Model.BootstrapCache[key], value)
			Model.Unlock()
		}
	}
	return
}

// DeleteBootstrapCache is exported internally to share between GoCore packages and should not be called directly by you.
func DeleteBootStrapFileCache(key string) (err error) {
	fname := CACHE_BOOTSTRAP_STORAGE_PATH + "/" + key + ".json"
	if extensions.DoesFileExist(fname) {
		err = os.Remove(fname)
		if err != nil {
			return err
		}
	}
	return
}

// DeleteAllBootstrapFileCache is exported internally to share between GoCore packages and should not be called directly by you.
func DeleteAllBootStrapFileCache() (err error) {
	if extensions.DoesFileExist(CACHE_BOOTSTRAP_STORAGE_PATH) {
		err = extensions.RemoveDirectory(CACHE_BOOTSTRAP_STORAGE_PATH)
		if err != nil {
			log.Println("DeleteAllBootStrapFileCache1", err)
			return err
		}
		os.MkdirAll(CACHE_BOOTSTRAP_STORAGE_PATH, 0777)
	}
	if extensions.DoesFileExist(CACHE_MANIFEST_STORAGE_PATH) {
		err = extensions.RemoveDirectory(CACHE_MANIFEST_STORAGE_PATH)
		if err != nil {
			log.Println("DeleteAllBootStrapFileCache2", err)
			return err
		}
		os.MkdirAll(CACHE_MANIFEST_STORAGE_PATH, 0777)
	}
	return
}

// LoadCachedBootStrapFileFromKeyIntoMemory is exported internally to share between GoCore packages and should not be called directly by you.
func LoadCachedBootStrapFromKeyIntoMemory(key string) (err error) {
	fname := CACHE_BOOTSTRAP_STORAGE_PATH + "/" + key + ".json"
	if extensions.DoesFileExist(fname) {
		var size int64
		size, err = extensions.GetFileSize(fname)
		if err != nil {
			return
		}
		if size > 0 {
			UpdateBootStrapMemoryCache(key, "")
			var data []string
			jsonData, err := extensions.ReadFile(fname)
			if err != nil {
				log.Println("Cache failed to read for " + fname + " deleting file and starting fresh.")
				DeleteBootStrapFileCache(key)
				return err
			}
			err = json.Unmarshal(jsonData, &data)
			if err != nil {
				return err
			}

			Model.RLock()
			_, ok := Model.BootstrapCache[key]
			Model.RUnlock()
			if ok {
				Model.Lock()
				Model.BootstrapCache[key] = data
				Model.Unlock()
			}
		}
	}
	return
}

// DoesHashExistInCache is exported internally to share between GoCore packages and should not be called directly by you.
func DoesHashExistInCache(key string, value string) (exists bool) {
	Model.RLock()
	caches, ok := Model.BootstrapCache[key]
	Model.RUnlock()
	if !ok {
		return exists
	} else {
		return utils.InArray(value, caches)
	}
}

// WriteManifestCacheFile is exported internally to share between GoCore packages and should not be called directly by you.
func WriteManifestCacheFile(key string) (err error) {
	ByteManifest.RLock()
	caches, ok := ByteManifest.Cache[key]
	ByteManifest.RUnlock()
	if ok {
		strjson, err := json.Marshal(caches)
		if err != nil {
			return err
		}
		err = os.WriteFile(CACHE_MANIFEST_STORAGE_PATH+"/"+key+".json", []byte(strjson), 0777)
		if err != nil {
			return err
		}
	}
	return nil
}

// UpdateManifestMemoryCache is exported internally to share between GoCore packages and should not be called directly by you.
func UpdateManifestMemoryCache(key string, value string, byteSize int) {
	ByteManifest.Lock()
	_, ok := ByteManifest.Cache[key]
	if !ok {
		ByteManifest.Cache[key] = make(map[string]int, 0)
		ByteManifest.Cache[key][value] = byteSize
	} else {
		ByteManifest.Cache[key][value] = byteSize
	}
	ByteManifest.Unlock()
	return
}

// DeleteManifestFileCache is exported internally to share between GoCore packages and should not be called directly by you.
func DeleteManifestFileCache(key string) (err error) {
	fname := CACHE_MANIFEST_STORAGE_PATH + "/" + key + ".json"
	if extensions.DoesFileExist(fname) {
		err = os.Remove(fname)
		if err != nil {
			return err
		}
	}
	return
}

// LoadCachedManifestFileFromKeyIntoMemory is exported internally to share between GoCore packages and should not be called directly by you.
func LoadCachedManifestFromKeyIntoMemory(key string) (err error) {
	fname := CACHE_MANIFEST_STORAGE_PATH + "/" + key + ".json"
	ByteManifest.RLock()
	_, ok := ByteManifest.Cache[key]
	ByteManifest.RUnlock()
	if extensions.DoesFileExist(fname) && !ok {
		var data map[string]int
		jsonData, err := extensions.ReadFile(fname)
		if err != nil {
			log.Println("Cache failed to read for " + fname + " deleting file and starting fresh.")
			DeleteManifestFileCache(key)
			return err
		}
		err = json.Unmarshal(jsonData, &data)
		if err != nil {
			return err
		}
		ByteManifest.Lock()
		ByteManifest.Cache[key] = data
		ByteManifest.Unlock()
	} else if !extensions.DoesFileExist(fname) {
		UpdateManifestMemoryCache(key, "", 0)
	}
	return
}

// DoesHashExistInManifestCache is exported internally to share between GoCore packages and should not be called directly by you.
func DoesHashExistInManifestCache(key string, value string) (exists bool) {
	ByteManifest.RLock()
	_, ok := ByteManifest.Cache[key]
	ByteManifest.RUnlock()
	if !ok {
		return exists
	} else {
		ByteManifest.RLock()
		_, ok = ByteManifest.Cache[key][value]
		ByteManifest.RUnlock()
		if !ok {
			return exists
		}
		return true
	}
}

// DeleteAllManifestFileCache is exported internally to share between GoCore packages and should not be called directly by you.
func DeleteAllManifestFileCache() (err error) {
	if extensions.DoesFileExist(CACHE_MANIFEST_STORAGE_PATH) {
		err = extensions.RemoveDirectory(CACHE_MANIFEST_STORAGE_PATH)
		if err != nil {
			log.Println("DeleteAllManifestFileCache", err)
			return err
		}
		os.MkdirAll(CACHE_MANIFEST_STORAGE_PATH, 0777)
	}
	return
}

// GetHTMLFile returns the html by path (key) from group cache
func GetHTMLFile(path string) (string, error) {
	var ctx context.Context
	var data []byte
	err := htmlFileCache.Get(ctx, path, groupcache.AllocatingByteSliceSink(&data))
	if err != nil {
		return "", err
	}

	return string(data[:]), err
}

// GetFile returns binary data by path(key) from group cache
func GetFile(path string) ([]byte, error) {
	var ctx context.Context
	var data []byte
	err := htmlFileCache.Get(ctx, path, groupcache.AllocatingByteSliceSink(&data))
	if err != nil {
		return data, err
	}

	return data, err
}

// GetString gets a value by Key from group cache
func GetString(key string) (string, error) {
	var ctx context.Context
	var data []byte
	err := stringCache.Get(ctx, key, groupcache.AllocatingByteSliceSink(&data))
	if err != nil {
		return "", err
	}

	return string(data[:]), err
}

// SetString sets a Key value pair in group cache
func SetString(key string, value string) error {

	var ctx context.Context
	setTempStringCache(key, value)
	var data []byte
	return stringCache.Get(ctx, key, groupcache.AllocatingByteSliceSink(&data))
}

// SetGroupCache will update the group cache http pool.  Use for dynamic systems that update at runtime.
func SetGroupCache(servers []string) {
	peers.Set(servers...)
}

// InitializeGroupCache creates the Peers for group cache and creates caches for multiple types.
func InitializeGroupCache(domain string) {

	//For now use the app domain, later we will read from a list of domains.
	if !utils.InArray(domain, allGroupCacheDomains) {
		// just in case recover happens on main app, we cannot initialize the same cache twice.
		peers = groupcache.NewHTTPPool(domain)
		allGroupCacheDomains = append(allGroupCacheDomains, domain)
		htmlFileCache = groupcache.NewGroup("htmlFileCache", 64<<20, groupcache.GetterFunc(handleHtmlFileCache))
		stringCache = groupcache.NewGroup("stringCache", 64<<20, groupcache.GetterFunc(handleStringCache))

		log.Println("Initialized Group Cache Succesfully.")
	}
}

// Handles group cache callback on getting http file cache requests.
func handleHtmlFileCache(ctx context.Context, key string, dest groupcache.Sink) error {
	fileName := key
	data, err := extensions.ReadFile(fileName)
	if err != nil {
		return err
	}

	dest.SetBytes(data)
	return nil
}

// Handles group cache callback on getting a string key value pair.
func handleStringCache(ctx context.Context, key string, dest groupcache.Sink) error {

	stringKey := key
	value := getTempStringCache(stringKey)
	dest.SetBytes([]byte(value))
	deleteTempStringCache(stringKey)

	return nil
}

// Safely locks a cache map and gets the value
func getTempStringCache(key string) (value string) {
	tempStringCacheSynced.RLock()
	value = tempStringCacheSynced.cache[key]
	tempStringCacheSynced.RUnlock()
	return
}

// Safely locks a cache map and sets the value
func setTempStringCache(key string, value string) {
	tempStringCacheSynced.Lock()
	tempStringCacheSynced.cache[key] = value
	tempStringCacheSynced.Unlock()
}

// Safely locks a cache map and deletes the value
func deleteTempStringCache(key string) {
	tempStringCacheSynced.Lock()
	delete(tempStringCacheSynced.cache, key)
	tempStringCacheSynced.Unlock()
}
