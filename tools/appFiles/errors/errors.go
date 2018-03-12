//Package errors provides tools to print stack traces
package errors

import (
	"fmt"

	"github.com/go-stack/stack"
)

func PrintStackTrace(err error, traceLevel int) string {

	stackTrace := "Error message:\n\t" + err.Error() + "\n\nStack Trace:\n"

	s := stack.Trace().TrimAbove(stack.Caller(traceLevel))

	for i := len(s) - 1; i > 0; i-- {
		stackTrace += fmt.Sprintf("\t%+v\n", fmt.Sprintf("%+v", s[i]))
	}

	return stackTrace
}
