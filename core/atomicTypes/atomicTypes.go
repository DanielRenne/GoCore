//Package atomicTypes provides object locking / unlocking for setting and getting.
package atomicTypes

import "sync"

//AtomicString provides a string object that is lock safe.
type AtomicString struct {
	valueSync sync.RWMutex
	value     string
}

//Get returns the string value
func (obj *AtomicString) Get() (value string) {
	obj.valueSync.RLock()
	value = obj.value
	obj.valueSync.RUnlock()
	return
}

//Set sets the string value
func (obj *AtomicString) Set(value string) {
	obj.valueSync.Lock()
	obj.value = value
	obj.valueSync.Unlock()
}

//AtomicInt provides an int object that is lock safe.
type AtomicInt struct {
	valueSync sync.RWMutex
	value     int
}

//Get returns the int value
func (obj *AtomicInt) Get() (value int) {
	obj.valueSync.RLock()
	value = obj.value
	obj.valueSync.RUnlock()
	return
}

//Set sets the int value
func (obj *AtomicInt) Set(value int) {
	obj.valueSync.Lock()
	obj.value = value
	obj.valueSync.Unlock()
}

//AtomicBool provides an bool object that is lock safe.
type AtomicBool struct {
	valueSync sync.RWMutex
	value     bool
}

//Get returns the bool value
func (obj *AtomicBool) Get() (value bool) {
	obj.valueSync.RLock()
	value = obj.value
	obj.valueSync.RUnlock()
	return
}

//Set sets the bool value
func (obj *AtomicBool) Set(value bool) {
	obj.valueSync.Lock()
	obj.value = value
	obj.valueSync.Unlock()
}
