package dbServices

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/DanielRenne/GoCore/core/atomicTypes"
	"github.com/DanielRenne/GoCore/core/serverSettings"
	"github.com/globalsign/mgo/bson"
)

const (
	collectionCacheMaxCount       = 10000
	collectionCacheManagementTime = 5
)

//CollectionCache provides DB object Caching functions.
type CollectionCache struct {
}

//CacheKey is the key lookup for the collectionCache
type CacheKey struct {
	collection string
	id         string
}

//CacheValue is the value for the collectionCache
type CacheValue struct {
	lastUpdate atomicTypes.AtomicTime
	value      []byte
}

var collectionCache sync.Map
var collectionCacheCount atomicTypes.AtomicInt

func init() {
	go CollectionCache{}.manage()
}

//Fetch will get the collection entity
func (cc CollectionCache) Fetch(collection string, id string, value interface{}) (ok bool) {

	ck := CacheKey{collection: collection, id: id}

	cvObj, found := collectionCache.Load(ck)
	if found {
		cv, parsed := cvObj.(*CacheValue)
		if parsed && cv != nil {
			ok = true
			if serverSettings.WebConfig.DbConnection.Driver == DATABASE_DRIVER_BOLTDB {
				json.Unmarshal(cv.value, value)
			} else {
				bson.Unmarshal(cv.value, value)
			}
			cv.lastUpdate.Set(time.Now())
			return
		}
	}

	return
}

//Store will store the collection object.
func (cc CollectionCache) Store(collection string, id string, value interface{}) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Panic Recovered at dbServices.collectionCache.Store():  ", r)
			return
		}
	}()
	ck := CacheKey{collection: collection, id: id}

	data := []byte{}

	if serverSettings.WebConfig.DbConnection.Driver == DATABASE_DRIVER_BOLTDB {
		data, _ = json.Marshal(value)
	} else {
		data, _ = bson.Marshal(value)
	}

	cv := CacheValue{value: data}
	cv.lastUpdate.Set(time.Now())

	_, ok := collectionCache.Load(ck)
	if !ok {
		collectionCache.Store(ck, &cv)
		count := collectionCacheCount.Get()
		count++
		collectionCacheCount.Set(count)
	}
}

//Count returns the length of the cache.
func (cc CollectionCache) Count() (value int) {
	value = collectionCacheCount.Get()
	return
}

//removeByKey will remove from the collection cache.
func (cc CollectionCache) removeByKey(key interface{}) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Panic Recovered at dbServices.collectionCache.removeByKey():  ", r)
			return
		}
	}()

	_, ok := collectionCache.Load(key)
	if ok {
		collectionCache.Delete(key)
		count := collectionCacheCount.Get()
		count--
		collectionCacheCount.Set(count)
	}
}

//Remove will remove from the collection cache.
func (cc CollectionCache) Remove(collection string, id string) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Panic Recovered at dbServices.collectionCache.Remove():  ", r)
			return
		}
	}()
	ck := CacheKey{collection: collection, id: id}
	_, ok := collectionCache.Load(ck)
	if ok {
		collectionCache.Delete(ck)

		count := collectionCacheCount.Get()
		count--
		collectionCacheCount.Set(count)
	}

}

//reduce will reduce the collection cache by a config time limit or document limit.
func (cc CollectionCache) reduce() {

	defer func() {
		if r := recover(); r != nil {
			log.Println("Panic Recovered at dbServices.collectionCache.reduce():  ", r)
			return
		}
	}()

	count := collectionCacheCount.Get()

	if count >= collectionCacheMaxCount { //Clear to 50%

		countDown := count

		collectionCache.Range(func(key interface{}, value interface{}) bool {

			go cc.removeByKey(key)

			if countDown == 0 {
				return false
			}

			countDown--

			return true
		})
	} else {

		staleTime := time.Now().Add(-1 * time.Hour)

		collectionCache.Range(func(key interface{}, value interface{}) bool {

			cv, parsed := value.(*CacheValue)
			if parsed && cv.lastUpdate.Get().Before(staleTime) {
				go cc.removeByKey(key)
			}

			return true
		})

	}

}

//manage checks the collectionCache every hour to reduce the cache by date or doc count
func (cc CollectionCache) manage() {

	defer func() {
		if r := recover(); r != nil {
			log.Println("Panic Recovered at dbServices.collectionCache.manage():  ", r)
			return
		}
	}()

	for {
		go cc.reduce()
		time.Sleep(time.Minute * collectionCacheManagementTime)
	}

}
