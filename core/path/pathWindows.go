//go:build windows
// +build windows

package path

// PathSeparator provides windows folder path separator.  Compiled only for windows
const PathSeparator = "\\"

// IsWindows is a flag to determine if the current OS is windows.
const IsWindows = true
