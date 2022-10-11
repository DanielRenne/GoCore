//go:build !windows
// +build !windows

package path

// PathSeparator provides windows or non-windows folder path separator.
const PathSeparator = "/"

// IsWindows is a flag to determine if the current OS is windows.
const IsWindows = false
