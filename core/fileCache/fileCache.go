//  Package fileCache will provide simple file content caching tools for in-Memory access to files.
//  It uses golang/groupcache to cache your data into memory on multiple HTTP Pool servers.
package fileCache

import (
	"log"
	"os"
	"sync"

	"encoding/json"
	"io/ioutil"

	"github.com/DanielRenne/GoCore/core/extensions"
	"github.com/DanielRenne/GoCore/core/serverSettings"
	"github.com/DanielRenne/GoCore/core/utils"
	"github.com/golang/groupcache"
)

const (
	CACHE_STORAGE_PATH = "/usr/local/goCore/caches"
	CACHE_JOBS         = "/usr/local/goCore/jobs"
)

const (
	CACHE_BOOTSTRAP_STORAGE_PATH = CACHE_STORAGE_PATH + "/bootstrap"
	CACHE_MANIFEST_STORAGE_PATH  = CACHE_STORAGE_PATH + "/manifests"
)

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
var Model model
var Jobs job
var ByteManifest byteManifest
var peers *groupcache.HTTPPool
var htmlFileCache *groupcache.Group
var stringCache *groupcache.Group

// contains the temporary string cache used to cache large strings.
var tempStringCacheSynced = struct {
	sync.RWMutex
	cache map[string]string
}{cache: make(map[string]string)}

func init() {
	os.MkdirAll(CACHE_BOOTSTRAP_STORAGE_PATH, 0777)
	os.MkdirAll(CACHE_JOBS, 0777)
	os.MkdirAll(CACHE_MANIFEST_STORAGE_PATH, 0777)
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

//Call Initialize in main before any calls to this package are performed.  serverSettings package must be initialized before fileCache.
func Initialize() {
	if serverSettings.WebConfig.Application.Domain != "" {
		initializeGroupCache(serverSettings.WebConfig.Application.Domain)
	}
	LoadJobsFile()
}

func WriteJobCacheFile() (err error) {
	strjson, err := json.Marshal(Jobs.Jobs)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(CACHE_JOBS+"/jobs.json", []byte(strjson), 0777)
	if err != nil {
		return err
	}
	return nil
}

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

func WriteBootStrapCacheFile(key string) (err error) {
	Model.RLock()
	caches, ok := Model.BootstrapCache[key]
	Model.RUnlock()
	if ok {
		strjson, err := json.Marshal(caches)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(CACHE_BOOTSTRAP_STORAGE_PATH+"/"+key+".json", []byte(strjson), 0777)
		if err != nil {
			return err
		}
	}
	return nil
}

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

func WriteManifestCacheFile(key string) (err error) {
	ByteManifest.RLock()
	caches, ok := ByteManifest.Cache[key]
	ByteManifest.RUnlock()
	if ok {
		strjson, err := json.Marshal(caches)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(CACHE_MANIFEST_STORAGE_PATH+"/"+key+".json", []byte(strjson), 0777)
		if err != nil {
			return err
		}
	}
	return nil
}

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

// Returns the html by path (key) from group cache
func GetHTMLFile(path string) (string, error) {
	var ctx groupcache.Context
	var data []byte
	err := htmlFileCache.Get(ctx, path, groupcache.AllocatingByteSliceSink(&data))
	if err != nil {
		return "", err
	}

	return string(data[:]), err
}

//Returns binary data by path(key) from group cache
func GetFile(path string) ([]byte, error) {
	var ctx groupcache.Context
	var data []byte
	err := htmlFileCache.Get(ctx, path, groupcache.AllocatingByteSliceSink(&data))
	if err != nil {
		return data, err
	}

	return data, err
}

// Gets a value by Key from group cache
func GetString(key string) (string, error) {
	var ctx groupcache.Context
	var data []byte
	err := stringCache.Get(ctx, key, groupcache.AllocatingByteSliceSink(&data))
	if err != nil {
		return "", err
	}

	return string(data[:]), err
}

// Sets a Key value pair in group cache
func SetString(key string, value string) error {

	var ctx groupcache.Context
	setTempStringCache(key, value)
	var data []byte
	return stringCache.Get(ctx, key, groupcache.AllocatingByteSliceSink(&data))
}

// Will update the group cache http pool.  Use for dynamic systems that update at runtime.
func SetGroupCache(servers []string) {
	peers.Set(servers...)
}

// Creates the Peers for group cache and creates caches for multiple types.
func initializeGroupCache(domain string) {

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
func handleHtmlFileCache(ctx groupcache.Context, key string, dest groupcache.Sink) error {
	fileName := key
	data, err := extensions.ReadFile(fileName)
	if err != nil {
		return err
	}

	dest.SetBytes(data)
	return nil
}

// Handles group cache callback on getting a string key value pair.
func handleStringCache(ctx groupcache.Context, key string, dest groupcache.Sink) error {

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
