// Package core is the main package contents of the GoCore application collection of packages and utilities
// Also the root contains some debugging/dumping variable functions
package core

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/DanielRenne/GoCore/core/extensions"
	"github.com/DanielRenne/GoCore/core/serverSettings"
	"github.com/davidrenne/reflections"
	"github.com/go-errors/errors"
)

type core_debug struct{}

var core_logger = log.New(os.Stdout, "", 0)

// No one will ever use this, but TransactionLog provides a thread-safe buffer/string if you have debug.EnableTransactionLog set to true you can call something like core.Debug.GetDump in many places and then eventually read the core.TransactionLog value when you need to (note you must manually clear it as it will just increase your memory usage the more logs are sent).  And always remember to use TransactionLogMutex.RLock and unlock when you are trying to read it in a busy system logging lots of stuff
var TransactionLog string

// See TransactionLog documentation
var EnableTransactionLog bool

// Debug is a base struct for all debug functions.
var Debug = core_debug{}

// LogStackTraces will show you where your logs are coming from in the stack.  Alternatively if you are using a GoCore full web app with a webConfig.json you can set CoreDebugStackTrace to true if you dont override this to false or true.  The default is true.
var LogStackTraces bool

// Logger can be overridden with log.New(os.Stdout, "", 0) to log to stdout or some other writer
var Logger = core_logger

// TransactionLogMutex is a mutex for the TransactionLog which should be used on your end to clear the value safely
var TransactionLogMutex *sync.RWMutex

func init() {
	LogStackTraces = false
	TransactionLogMutex = &sync.RWMutex{}
}

// Nop is a dummy function that can be called in source files where
// other debug functions are constantly added and removed.
// That way import "github.com/ungerik/go-start/debug" won't cause an error when
// no other debug function is currently used.
// Arbitrary objects can be passed as arguments to avoid "declared and not used"
// error messages when commenting code out and in.
// The result is a nil interface{} dummy value.
func (obj *core_debug) Nop(dummiesIn ...interface{}) (dummyOut interface{}) {
	return nil
}

// CallStackInfo returns a string with the call stack info.
func (obj *core_debug) CallStackInfo(skip int) (info string) {
	pc, file, line, ok := runtime.Caller(skip)
	if ok {
		funcName := runtime.FuncForPC(pc).Name()
		info += fmt.Sprintf("In function %s()", funcName)
	}
	for i := 0; ok; i++ {
		info += fmt.Sprintf("\n%s:%d", file, line)
		_, file, line, ok = runtime.Caller(skip + i)
	}
	return info
}

// PrintCallStack prints the call stack info.
func (obj *core_debug) PrintCallStack() {
	debug.PrintStack()
}

// LogCallStack logs the call stack info.
func (obj *core_debug) LogCallStack() {
	log.Print(obj.Stack())
}

func (obj *core_debug) Stack() string {
	return string(debug.Stack())
}

func (obj *core_debug) formatValue(value interface{}) string {
	return fmt.Sprintf("\n     Type: %T\n    Value: %v\nGo Syntax: %#v", value, value, value)
}

func (obj *core_debug) formatCallstack(skip int) string {
	return fmt.Sprintf("\nCallstack: %s", obj.CallStackInfo(skip+1))
}

// FormatSkip formats a value with callstack info.
func (obj *core_debug) FormatSkip(skip int, value interface{}) string {
	return obj.formatValue(value) + obj.formatCallstack(skip+1)
}

// Format formats a value with callstack info.
func (obj *core_debug) Format(value interface{}) string {
	return obj.FormatSkip(2, value)
}

// IsZeroOfUnderlyingType returns true if the value is the zero value (nil) for its type.
func IsZeroOfUnderlyingType(x interface{}) bool {
	return reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}

// IsZeroOfUnderlyingType2 returns true if the value is the zero value (nil) for its type.
func IsZeroOfUnderlyingType2(x interface{}) bool {
	return x == reflect.Zero(reflect.TypeOf(x)).Interface()
}

// IsNil returns true if the value is the zero value (nil) for its type.
func IsNil(c interface{}) bool {
	return c == nil || (reflect.ValueOf(c).Kind() == reflect.Ptr && reflect.ValueOf(c).IsNil())
}

// HandleError is a helper function that will log an error and return it with the callers line and file.
func (obj *core_debug) HandleError(err error) (s string) {
	if err != nil {
		// notice that we're using 1, so it will actually log the where
		// the error happened, 0 = this function, we don't want that.
		_, fn, line, _ := runtime.Caller(1)
		fileNameParts := strings.Split(fn, "/")
		return fmt.Sprintf("  Error Info: %s Line %d. ErrorType: %v", fileNameParts[len(fileNameParts)-1], line, err)
	}
	return ""
}

// ErrLineAndFile returns the line and file of the error.
func (obj *core_debug) ErrLineAndFile(err error) (s string) {
	if err != nil {
		// notice that we're using 1, so it will actually log the where
		// the error happened, 0 = this function, we don't want that.
		_, fn, line, _ := runtime.Caller(1)
		fileNameParts := strings.Split(fn, "/")
		return fmt.Sprintf("%s Line %d", fileNameParts[len(fileNameParts)-1], line)
	}
	return ""
}

// Dump is a helper function that will log unlimited values to print to stdout or however you have log setup if you overload core/Logger
func (obj *core_debug) Dump(valuesOriginal ...interface{}) {
	t := time.Now()
	l := "!!!!!!!!!!!!! DEBUG " + t.Format("2006-01-02 15:04:05.000000") + "!!!!!!!!!!!!!\n\n"
	Logger.Println(l)

	if EnableTransactionLog {
		TransactionLogMutex.Lock()
		TransactionLog += l
		TransactionLogMutex.Unlock()
	}
	for _, value := range valuesOriginal {
		l := obj.dumpBase(value)
		Logger.Print(l)
		if EnableTransactionLog {
			TransactionLogMutex.Lock()
			TransactionLog += l
			TransactionLogMutex.Unlock()
		}
	}
	l = obj.ThrowAndPrintError()
	Logger.Print(l)

	if EnableTransactionLog {
		TransactionLogMutex.Lock()
		TransactionLog += l
		TransactionLogMutex.Unlock()
	}
	l = "!!!!!!!!!!!!! ENDDEBUG " + t.Format("2006-01-02 15:04:05.000000") + "!!!!!!!!!!!!!\n\n"
	Logger.Println(l)
	if EnableTransactionLog {
		TransactionLogMutex.Lock()
		TransactionLog += l
		TransactionLogMutex.Unlock()
	}
}

// GetDump is a helper function that will log unlimited values which will return a string representation of what was logged
func (obj *core_debug) GetDump(valuesOriginal ...interface{}) (output string) {
	for _, value := range valuesOriginal {
		output += obj.dumpBase(value)
	}
	//output += obj.ThrowAndPrintError()
	return output
}

func (obj *core_debug) GetDumpWithInfo(valuesOriginal ...interface{}) (output string) {
	t := time.Now()
	return obj.GetDumpWithInfoAndTimeString(t.String(), valuesOriginal...)
}

// GetDumpWithInfoAndTimeString is a helper function that will log unlimited values which will return a string representation of what was logged but allows you to pass your own time string in a case of timezone offsets
func (obj *core_debug) GetDumpWithInfoAndTimeString(timeStr string, valuesOriginal ...interface{}) (output string) {
	l := "\n!!!!!!!!!!!!! DEBUG " + timeStr + "!!!!!!!!!!!!!\n\n"
	output += l

	if EnableTransactionLog {
		TransactionLogMutex.Lock()
		TransactionLog += l
		TransactionLogMutex.Unlock()
	}

	for _, value := range valuesOriginal {
		output += obj.dumpBase(value) + "\n"
	}

	if EnableTransactionLog {
		TransactionLogMutex.Lock()
		TransactionLog += output
		TransactionLogMutex.Unlock()
	}

	l = obj.ThrowAndPrintError()
	output += l

	if EnableTransactionLog {
		TransactionLogMutex.Lock()
		TransactionLog += l
		TransactionLogMutex.Unlock()
	}

	l = "!!!!!!!!!!!!! ENDDEBUG " + timeStr + "!!!!!!!!!!!!!\n"
	output += l
	if EnableTransactionLog {
		TransactionLogMutex.Lock()
		TransactionLog += l
		TransactionLogMutex.Unlock()
	}
	return output
}

func (obj *core_debug) dumpBase(values ...interface{}) (output string) {
	var jsonString string
	var err error
	var structKeys []string
	if Logger != nil {
		for _, value := range values {
			isAllJSON := true
			var kind string
			kind = strings.TrimSpace(fmt.Sprintf("%T", value))
			kindFormatted := strings.TrimSpace(fmt.Sprintf("%T", value))
			// var pieces = strings.Split(kind, " ")
			// if pieces[0] == "struct" || strings.Index(pieces[0], "model.") != -1 || strings.Index(pieces[0], "viewModel.") != -1 {
			if !IsNil(value) {
				kind = reflections.ReflectKind(value)
				structKeys, err = reflections.FieldsDeep(value)
				if err == nil {
					for _, field := range structKeys {
						jsonString, err = reflections.GetFieldTag(value, field, "json")
						if err != nil {
							isAllJSON = false
						}
						if jsonString == "" {
							isAllJSON = false
						}
					}
				} else {
					isAllJSON = false
				}
			} else {
				isAllJSON = false
			}

			if kindFormatted != "[]uint8" && (isAllJSON || kind == "map" || kind == "bson.M" || kind == "slice") {
				var rawBytes []byte
				rawBytes, err = json.MarshalIndent(value, "", "\t")
				if err == nil {
					if kind == "slice" || kind[:2] == "[]" {
						valReflected := reflect.ValueOf(value)
						output += fmt.Sprintf("#### %-39s [len:%s]####\n%+v\n", kindFormatted, extensions.IntToString(valReflected.Len()), string(rawBytes[:]))
						output += "\n\n" + fmt.Sprintf("%#v", value)
					} else {
						output += fmt.Sprintf("#### %-39s ####\n%+v\n", kindFormatted, string(rawBytes[:]))
						output += "\n\n" + fmt.Sprintf("%#v", value)
					}
				}
			} else {
				if strings.TrimSpace(kind) == "string" {
					var stringVal = fmt.Sprintf("%+v", value)
					position := strings.Index(stringVal, "Desc->")
					if position == -1 {
						if !extensions.IsPrintable(stringVal) {
							kind += " (non printables -> dump hex)"
							stringVal = hex.Dump([]byte(stringVal))
						}
						valReflected := reflect.ValueOf(value)
						output += fmt.Sprintf("#### %-39s [len:%s]####\n%s\n", kindFormatted, extensions.IntToString(valReflected.Len()), stringVal)
					} else {
						output += stringVal[6:] + " --> "
					}
				} else if kindFormatted == "[]uint8" || kind[:2] == "[]" || strings.TrimSpace(kind) == "array" {
					if kindFormatted == "[]uint8" {
						val, ok := value.([]uint8)
						if ok {
							if !extensions.IsPrintable(string(val)) {
								kind += " (non printables -> dump hex of a []byte)"
								output = hex.Dump(value.([]uint8))
							} else {
								valReflected := reflect.ValueOf(value)
								output += fmt.Sprintf("#### %-39s [len:%s]####\n%s\n", kindFormatted, extensions.IntToString(valReflected.Len()), string(val))
							}
						}
					} else {
						valReflected := reflect.ValueOf(value)
						output += fmt.Sprintf("#### %-39s [len:%s]####\n%#v\n", kindFormatted, extensions.IntToString(valReflected.Len()), value)
					}
				} else {
					output += fmt.Sprintf("#### %-39s ####\n%#v\n", kindFormatted, value)
				}
			}
		}
	}
	return output
}

// ThrowAndPrintError is a helper function that will throw a fake error and get the callstack and return it as a string (you probably shouldnt use this)
func (obj *core_debug) ThrowAndPrintError() (output string) {

	serverSettings.WebConfigMutex.RLock()
	ok := serverSettings.WebConfig.Application.CoreDebugStackTrace
	serverSettings.WebConfigMutex.RUnlock()
	if ok || LogStackTraces {
		output += "\n"
		errorInfo := obj.ThrowError()
		stack := strings.Split(errorInfo.ErrorStack(), "\n")
		if len(stack) >= 8 {
			output += "\nDump Caller:"
			output += "\n---------------"
			//output += strings.Join(stack, ",")
			if len(stack) >= 7 {
				output += "\n golines ==> " + strings.TrimSpace(stack[6])
			}
			if len(stack) >= 8 {
				output += "\n         ==> " + strings.TrimSpace(stack[7])
			}
			if len(stack) >= 9 {
				output += "\n         ==> " + strings.TrimSpace(stack[8])
			}
			if len(stack) >= 10 {
				output += "\n         ==> " + strings.TrimSpace(stack[9])
			}
			if len(stack) >= 11 {
				output += "\n         ==> " + strings.TrimSpace(stack[10])
			}
			if len(stack) >= 12 {
				output += "\n         ==> " + strings.TrimSpace(stack[11])
			}
			if len(stack) >= 13 {
				output += "\n         ==> " + strings.TrimSpace(stack[12])
			}
			if len(stack) >= 14 {
				output += "\n         ==> " + strings.TrimSpace(stack[13])
			}
			if len(stack) >= 15 {
				output += "\n         ==> " + strings.TrimSpace(stack[14])
			}
			if len(stack) >= 16 {
				output += "\n         ==> " + strings.TrimSpace(stack[15])
			}
			output += "\n---------------"
			output += "\n"
			output += "\n"
		}
	}
	return output
}

// ThrowError is a helper function that will throw a fake error and get the callstack and return it as an error (you probably shouldnt use this)
func (obj *core_debug) ThrowError() *errors.Error {
	return errors.Errorf("Debug Dump")
}

// GetDump is a helper function that will return a string of the dump of the values passed in
func GetDump(valuesOriginal ...interface{}) string {
	return Debug.GetDump(valuesOriginal...)
}

// Dump is a helper function that will dump the values passed to it
func Dump(valuesOriginal ...interface{}) {
	Debug.Dump(valuesOriginal...)
}
