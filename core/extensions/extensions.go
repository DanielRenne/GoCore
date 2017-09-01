package extensions

import (
	"bytes"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type Version struct {
	Major          int
	Minor          int
	Revision       int
	MajorString    string
	MinorString    string
	RevisionString string
	Value          string
}

/*
* leftPad and rightPad just repoeat the padStr the indicated
* number of times
*
 */
func LeftPad(s string, padStr string, pLen int) string {
	return strings.Repeat(padStr, pLen) + s
}
func RightPad(s string, padStr string, pLen int) string {
	return s + strings.Repeat(padStr, pLen)
}

/* the Pad2Len functions are generally assumed to be padded with short sequences of strings
* in many cases with a single character sequence
*
* so we assume we can build the string out as if the char seq is 1 char and then
* just substr the string if it is longer than needed
*
* this means we are wasting some cpu and memory work
* but this always get us to want we want it to be
*
* in short not optimized to for massive string work
*
* If the overallLen is shorter than the original string length
* the string will be shortened to this length (substr)
*
 */
func RightPad2Len(s string, padStr string, overallLen int) string {
	var padCountInt int
	padCountInt = 1 + ((overallLen - len(padStr)) / len(padStr))
	var retStr = s + strings.Repeat(padStr, padCountInt)
	return retStr[:overallLen]
}

func LeftPad2Len(s string, padStr string, overallLen int) string {
	var padCountInt int
	padCountInt = 1 + ((overallLen - len(padStr)) / len(padStr))
	var retStr = strings.Repeat(padStr, padCountInt) + s
	return retStr[(len(retStr) - overallLen):]
}

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

func IsPrintable(s string) bool {
	for _, c := range s {
		if (c < 32 || c > 126) && c != 10 && c != 13 {
			return false
		}
	}
	return true
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

func StringToInt(val string) int {

	r, err := strconv.Atoi(val)
	if err != nil {
		return 0
	}
	return r
}

func StringToUInt64(val string) uint64 {
	i, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0
	}
	return i
}

func IntToString(val int) string {
	return strconv.Itoa(val)
}

func Int32ToString(val int32) string {
	return strconv.Itoa(Int32ToInt(val))
}

func Int64ToString(val int64) string {
	return strconv.FormatInt(val, 10)
}

func Int64ToInt32(val int64) (ret int) {
	tempLong := ((val >> 32) << 32) //shift it right then left 32 bits, which zeroes the lower half of the long
	ret = (int)(val - tempLong)
	return ret
}

func Int32ToInt(val int32) (ret int) {
	tempLong := ((val >> 32) << 32) //shift it right then left 32 bits, which zeroes the lower half of the long
	ret = (int)(val - tempLong)
	return ret
}

func BoolToString(val bool) string {
	return strconv.FormatBool(val)
}

func StringToBool(val string) bool {
	r, err := strconv.ParseBool(val)
	if err != nil {
		return false
	}
	return r
}

func (obj *Version) Init(value string) {
	versionInfo := strings.Split(value, ".")

	obj.MajorString = versionInfo[0]
	obj.MinorString = versionInfo[1]
	obj.RevisionString = versionInfo[2]
	obj.Value = value

	if val, err := strconv.Atoi(versionInfo[0]); err == nil {
		obj.Major = val
	}

	if val, err := strconv.Atoi(versionInfo[1]); err == nil {
		obj.Minor = val
	}

	if val, err := strconv.Atoi(versionInfo[2]); err == nil {
		obj.Revision = val
	}
}

func GenPackageImport(name string, imports []string) string {

	val := "package " + name + "\n\n"
	val += "import(\n"
	for _, imp := range imports {
		if imp == "" {
			continue
		}
		val += "\t\"" + imp + "\"\n"
	}
	val += ")\n\n"

	return val
}

func MakeFirstLowerCase(s string) string {

	if len(s) < 2 {
		return strings.ToLower(s)
	}

	bts := []byte(s)

	lc := bytes.ToLower([]byte{bts[0]})
	rest := bts[1:]

	return string(bytes.Join([][]byte{lc, rest}, nil))
}

func ExtractArgsWithinBrackets(str string) (res []string) {

	brackets := &unicode.RangeTable{
		R16: []unicode.Range16{
			// {0x0028, 0x0029, 1}, // ( )
			// {0x005b, 0x005d, 1}, // [ ]
			{0x007b, 0x007d, 1}, // { }
		},
	}

	isBracket := func(r rune) bool {
		if unicode.In(r, brackets) {
			return true
		}
		return false
	}

	res = strings.FieldsFunc(str, isBracket)
	return
}

func Random(min, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(max-min) + min
}
