// Debug functions.
package core

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/debug"
)

// Set Logger to nil to suppress debug output
var logger = log.New(os.Stdout, "", 0)

type Debug interface{}

// Nop is a dummy function that can be called in source files where
// other debug functions are constantly added and removed.
// That way import "github.com/ungerik/go-start/debug" won't cause an error when
// no other debug function is currently used.
// Arbitrary objects can be passed as arguments to avoid "declared and not used"
// error messages when commenting code out and in.
// The result is a nil interface{} dummy value.
func (self *Debug) Nop(dummiesIn ...interface{}) (dummyOut interface{}) {
	return nil
}

func (self *Debug) Log() log {
	return logger
}

//func (self *Debug) CallStackInfo(skip int) (info string) {
func CallStackInfo(skip int) (info string) {
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

func (self *Debug) PrintCallStack() {
	debug.PrintStack()
}

func (self *Debug) LogCallStack() {
	log.Print(Stack())
}

func Stack() string {
	return string(debug.Stack())
}

func formatValue(value interface{}) string {
	return fmt.Sprintf("\n     Type: %T\n    Value: %v\nGo Syntax: %#v", value, value, value)
}

func formatCallstack(skip int) string {
	//return fmt.Sprintf("\nCallstack: %s", self.CallStackInfo(skip+1))
	//return fmt.Sprintf("\nCallstack: %s", Debug.CallStackInfo(skip+1))
	return fmt.Sprintf("\nCallstack: %s", CallStackInfo(skip+1))
}

func FormatSkip(skip int, value interface{}) string {
	return formatValue(value) + formatCallstack(skip+1)
}

func (self *Debug) Format(value interface{}) string {
	return FormatSkip(2, value)
}

func (self *Debug) Dump(values ...interface{}) {
	if self.Log() != nil {
		for _, value := range values {
			//self.Log().Println(fmt.Printf("%+v\n", value))
			logger.Println(fmt.Printf("%+v\n", value))
		}
	}
}

func (self *Debug) GetDump(values ...interface{}) string {
	var buffer bytes.Buffer
	for _, value := range values {
		buffer.WriteString(fmt.Sprintf("%+v\n", value))
	}
	return buffer.String()
}

func (self *Debug)  Print(values ...interface{}) {
	if self.Log() != nil {
		//self.Log().Print(values...)
		logger.Print(values...)
	}
}

func (self *Debug)  Printf(format string, values ...interface{}) {
	if self.Log() != nil {
		//self.Log().Printf(format, values...)
		logger.Printf(format, values...)
	}
}