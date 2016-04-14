package extensions

import (
	"fmt"
	"strconv"
	"strings"
)

func TrimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}

func PrintKiloBytes(bytes int64) string {

	var kilobytes float64
	kilobytes = float64(bytes / 1024)

	return fmt.Sprint(FloatToString(kilobytes, 2), " kB")
}

func PrintMegaBytes(bytes int64) string {

	var kilobytes float64
	kilobytes = float64(bytes / 1024)

	var megabytes float64
	megabytes = kilobytes / 1024 // cast to type float64

	return fmt.Sprint(FloatToString(megabytes, 2), " MB")
}

func PrintZettaBytes(bytes int64) string {

	var kilobytes float64
	kilobytes = float64(bytes / 1024)

	var megabytes float64
	megabytes = (kilobytes / 1024) // cast to type float64

	var gigabytes float64
	gigabytes = (megabytes / 1024)

	var terabytes float64
	terabytes = (gigabytes / 1024)

	var petabytes float64
	petabytes = (terabytes / 1024)

	var exabytes float64
	exabytes = (petabytes / 1024)

	var zettabytes float64
	zettabytes = (exabytes / 1024)

	return fmt.Sprint(FloatToString(zettabytes, 2), " ZB")
}

func FloatToString(input_num float64, decimals int) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', decimals, 64)
}
