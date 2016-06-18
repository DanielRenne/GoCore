//  Package fileCache will provide simple file content caching tools for in-Memory access to files.
//  It uses golang/groupcache to cache your data into memory on multiple HTTP Pool servers.
package fileCache

import (
	"github.com/DanielRenne/GoCore/core/extensions"
	"github.com/DanielRenne/GoCore/core/serverSettings"
	"github.com/golang/groupcache"
	"log"
	"sync"
)

var peers *groupcache.HTTPPool
var htmlFileCache *groupcache.Group
var stringCache *groupcache.Group

// contains the temporary string cache used to cache large strings.
var tempStringCacheSynced = struct {
	sync.RWMutex
	cache map[string]string
}{cache: make(map[string]string)}

// func init() {
// 	if serverSettings.WebConfig.Application.Domain != "" {
// 		initializeGroupCache(serverSettings.WebConfig.Application.Domain)
// 	}
// }

//Call Initialize in main before any calls to this package are performed.  serverSettings package must be initialized before fileCache.
func Initialize() {
	if serverSettings.WebConfig.Application.Domain != "" {
		initializeGroupCache(serverSettings.WebConfig.Application.Domain)
	}
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
	peers = groupcache.NewHTTPPool(domain)
	htmlFileCache = groupcache.NewGroup("htmlFileCache", 64<<20, groupcache.GetterFunc(handleHtmlFileCache))
	stringCache = groupcache.NewGroup("stringCache", 64<<20, groupcache.GetterFunc(handleStringCache))

	log.Println("Initialized Group Cache Succesfully.")

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
