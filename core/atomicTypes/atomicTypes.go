// Package atomicTypes provides object locking / unlocking for setting and getting.
package atomicTypes

import (
	"sync"
	"time"
)

// AtomicString provides a string object that is lock safe.
type AtomicString struct {
	valueSync sync.RWMutex
	value     string
}

// Get returns the string value
func (obj *AtomicString) Get() (value string) {
	obj.valueSync.RLock()
	value = obj.value
	obj.valueSync.RUnlock()
	return
}

// Set sets the string value
func (obj *AtomicString) Set(value string) {
	obj.valueSync.Lock()
	obj.value = value
	obj.valueSync.Unlock()
}

// AtomicInt provides an int object that is lock safe.
type AtomicInt struct {
	valueSync sync.RWMutex
	value     int
}

// Get returns the int value
func (obj *AtomicInt) Get() (value int) {
	obj.valueSync.RLock()
	value = obj.value
	obj.valueSync.RUnlock()
	return
}

// Set sets the int value
func (obj *AtomicInt) Set(value int) {
	obj.valueSync.Lock()
	obj.value = value
	obj.valueSync.Unlock()
}

// Increment allows for the value to safely be increased by 1
func (obj *AtomicInt) Increment() int {
	obj.valueSync.Lock()
	defer obj.valueSync.Unlock()

	obj.value++
	return obj.value
}

// Decrement allows for the value to safely be decreased by 1
func (obj *AtomicInt) Decrement() int {
	obj.valueSync.Lock()
	defer obj.valueSync.Unlock()

	obj.value--
	return obj.value
}

// Add allows for the value to safely be add by n
func (obj *AtomicInt) Add(n int) int {
	obj.valueSync.Lock()
	defer obj.valueSync.Unlock()

	obj.value += n
	return obj.value
}

// AtomicUInt16 provides an uint16 object that is lock safe.
type AtomicUInt16 struct {
	valueSync sync.RWMutex
	value     uint16
}

// Get returns the int value
func (obj *AtomicUInt16) Get() (value uint16) {
	obj.valueSync.RLock()
	value = obj.value
	obj.valueSync.RUnlock()
	return
}

// Set sets the int value
func (obj *AtomicUInt16) Set(value uint16) {
	obj.valueSync.Lock()
	obj.value = value
	obj.valueSync.Unlock()
}

// AtomicUInt32 provides an uint16 object that is lock safe.
type AtomicUInt32 struct {
	valueSync sync.RWMutex
	value     uint32
}

// Get returns the int value
func (obj *AtomicUInt32) Get() (value uint32) {
	obj.valueSync.RLock()
	value = obj.value
	obj.valueSync.RUnlock()
	return
}

// Set sets the int value
func (obj *AtomicUInt32) Set(value uint32) {
	obj.valueSync.Lock()
	obj.value = value
	obj.valueSync.Unlock()
}

// AtomicBool provides an bool object that is lock safe.
type AtomicBool struct {
	valueSync sync.RWMutex
	value     bool
}

// Get returns the bool value
func (obj *AtomicBool) Get() (value bool) {
	obj.valueSync.RLock()
	value = obj.value
	obj.valueSync.RUnlock()
	return
}

// Set sets the bool value
func (obj *AtomicBool) Set(value bool) {
	obj.valueSync.Lock()
	obj.value = value
	obj.valueSync.Unlock()
}

// ToggleTrue sets the bool value only if the check back true; return if the value was set
func (obj *AtomicBool) ToggleTrue() (changed bool) {
	obj.valueSync.Lock()
	defer obj.valueSync.Unlock()

	changed = false
	if obj.value == false {
		obj.value = true

		changed = true
	}

	return
}

// AtomicTime provides a time.Time object that is lock safe.
type AtomicTime struct {
	valueSync sync.RWMutex
	value     time.Time
}

// Get returns the time.Time value
func (obj *AtomicTime) Get() (value time.Time) {
	obj.valueSync.RLock()
	value = obj.value
	obj.valueSync.RUnlock()
	return
}

// Set sets the time.Time value
func (obj *AtomicTime) Set(value time.Time) {
	obj.valueSync.Lock()
	obj.value = value
	obj.valueSync.Unlock()
}

// AtomicByteArray provides a []byte object that is lock safe.
type AtomicByteArray struct {
	valueSync sync.RWMutex
	value     []byte
}

// Get returns the []byte value
func (obj *AtomicByteArray) Get() (value []byte) {
	obj.valueSync.RLock()
	value = obj.value
	obj.valueSync.RUnlock()
	return
}

// Set sets the []byte value
func (obj *AtomicByteArray) Set(value []byte) {
	obj.valueSync.Lock()
	obj.value = value
	obj.valueSync.Unlock()
}

// AtomicBoolArray provides a []bool object that is lock safe.
type AtomicBoolArray struct {
	valueSync sync.RWMutex
	value     []bool
}

// Get returns the []bool value
func (obj *AtomicBoolArray) Get() (value []bool) {
	obj.valueSync.RLock()
	value = obj.value
	obj.valueSync.RUnlock()
	return
}

// Set sets the []bool value
func (obj *AtomicBoolArray) Set(value []bool) {
	obj.valueSync.Lock()
	obj.value = value
	obj.valueSync.Unlock()
}

// SetByIndex sets a bool by index value
func (obj *AtomicBoolArray) SetByIndex(index int, value bool) {
	obj.valueSync.Lock()
	if len(obj.value) > index {
		obj.value[index] = value
	}
	obj.valueSync.Unlock()
}

// AtomicFloat64 provides a float64 object that is lock safe.
type AtomicFloat64 struct {
	valueSync sync.RWMutex
	value     float64
}

// Get returns the float64 value
func (obj *AtomicFloat64) Get() (value float64) {
	obj.valueSync.RLock()
	value = obj.value
	obj.valueSync.RUnlock()
	return
}

// Set sets the float64 value
func (obj *AtomicFloat64) Set(value float64) {
	obj.valueSync.Lock()
	obj.value = value
	obj.valueSync.Unlock()
}
