// Debug functions.
package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/DanielRenne/GoCore/core/extensions"
	"github.com/go-errors/errors"
	"github.com/oleiade/reflections"
	"log"
	"os"
	"reflect"
	"runtime"
	"runtime/debug"
)

type core_debug struct{}

var core_logger = log.New(os.Stdout, "", 0)
var Debug = core_debug{}
var Logger = core_logger

// Nop is a dummy function that can be called in source files where
// other debug functions are constantly added and removed.
// That way import "github.com/ungerik/go-start/debug" won't cause an error when
// no other debug function is currently used.
// Arbitrary objects can be passed as arguments to avoid "declared and not used"
// error messages when commenting code out and in.
// The result is a nil interface{} dummy value.
func (self *core_debug) Nop(dummiesIn ...interface{}) (dummyOut interface{}) {
	return nil
}

func (self *core_debug) CallStackInfo(skip int) (info string) {
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

func (self *core_debug) PrintCallStack() {
	debug.PrintStack()
}

func (self *core_debug) LogCallStack() {
	log.Print(self.Stack())
}

func (self *core_debug) Stack() string {
	return string(debug.Stack())
}

func (self *core_debug) formatValue(value interface{}) string {
	return fmt.Sprintf("\n     Type: %T\n    Value: %v\nGo Syntax: %#v", value, value, value)
}

func (self *core_debug) formatCallstack(skip int) string {
	return fmt.Sprintf("\nCallstack: %s", self.CallStackInfo(skip+1))
}

func (self *core_debug) FormatSkip(skip int, value interface{}) string {
	return self.formatValue(value) + self.formatCallstack(skip+1)
}

func (self *core_debug) Format(value interface{}) string {
	return self.FormatSkip(2, value)
}

func (self *core_debug) DumpQuiet(values ...interface{}) {
	// uncomment below to find your callers to quiet
	self.Print("Silently not dumping " + extensions.IntToString(len(values)) + " values")
	//Logger.Println("DumpQuiet has " + extensions.IntToString(len(values)) + " parameters called")
	//Logger.Println("")
	//self.ThrowAndPrintError()
}

func (self *core_debug) Dump(values ...interface{}) {
	Logger.Println("!!!!!!!!!!!!!DEBUG!!!!!!!!!!!!!")
	Logger.Println("")
	Logger.Println("")
	Logger.Println("")

	var err error
	var jsonString string
	isAllJSON := true
	var structKeys []string
	if Logger != nil {
		for _, value := range values {
			Logger.Println("Instance Type:" + reflect.TypeOf(value).Name())

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
			}
			if isAllJSON {
				var rawBytes []byte
				rawBytes, err = json.MarshalIndent(value, "", "\t")
				if err == nil {
					value = string(rawBytes[:])
				}
				Logger.Println(fmt.Printf("%+v\n", value))
			} else {
				Logger.Println(fmt.Printf("%+v\n", value))
			}
		}
	}
	Logger.Println("")
	Logger.Println("")
	self.ThrowAndPrintError()
	Logger.Println("")
	Logger.Println("")
	Logger.Println("!!!!!!!!!!!!!ENDDEBUG!!!!!!!!!!!!!")
}

func (self *core_debug) ThrowAndPrintError() {
	Logger.Println("")
	Logger.Println("Dump Caller:")
	Logger.Println("")
	errorInfo := self.ThrowError()
	Logger.Println(errorInfo.ErrorStack())
}

func (self *core_debug) ThrowError() *errors.Error {
	return errors.Errorf("Debug Dump")
}

func (self *core_debug) GetDump(values ...interface{}) string {
	var buffer bytes.Buffer
	for _, value := range values {
		buffer.WriteString("(" + reflect.TypeOf(value).Name() + ")" + fmt.Sprintf("%+v\n", value))
	}
	return buffer.String()
}

func (self *core_debug) Print(values ...interface{}) {
	if Logger != nil {
		Logger.Print(values...)
	}
}

func (self *core_debug) Printf(format string, values ...interface{}) {
	if Logger != nil {
		Logger.Printf(format, values...)
	}
}
