// Package extensions provides a set of functions to work with files, directories, and converting data types
package extensions

import (
	"bytes"
	cryptoRand "crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/forPelevin/gomoji"
)

// Version is a type that represents a version number.
type Version struct {
	Major          int
	Minor          int
	Revision       int
	MajorString    string
	MinorString    string
	RevisionString string
	Value          string
}

// FilePath is a type that represents a file path.
type FilePath struct {
	Name string `json:"Name"`
	Path string `json:"Path"`
	Type string `json:"Type"`
}

// ToString converts a FilePath to a string.
func (obj *FilePath) ToString() (str string) {
	return obj.Path + string(os.PathSeparator) + obj.Name + " | " + obj.Type
}

// GexToInt converts a string to an int.
func HexToInt(hexStr string) int {
	// remove 0x suffix if found in the input string
	cleaned := strings.Replace(hexStr, "0x", "", -1)

	// base 16 for hexadecimal
	result, _ := strconv.ParseUint(cleaned, 16, 64)
	return int(result)
}

// GetDecimalAndStringFromHex takes a hex string and returns the decimal value and the string value.
func GetDecimalAndStringFromHex(twoDigitHexCode string) (decimalValue int, asciiString string, err error) {
	decoded, err := hex.DecodeString(twoDigitHexCode)
	if err != nil {
		return
	}
	decimalValue = int(decoded[0])
	asciiString = string(string(decoded))
	return
}

// LeftPad pads a string with a specified character to a specified length.
func LeftPad(s string, padStr string, pLen int) string {
	return strings.Repeat(padStr, pLen) + s
}

// RightPad pads a string with a specified character to a specified length.
func RightPad(s string, padStr string, pLen int) string {
	return s + strings.Repeat(padStr, pLen)
}

/* RightPad2Len functions are generally assumed to be padded with short sequences of strings
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
	var padCountInt int = 1 + ((overallLen - len(padStr)) / len(padStr))
	var retStr = s + strings.Repeat(padStr, padCountInt)
	return retStr[:overallLen]
}

// LeftPad2Len functions are generally assumed to be padded with short sequences of strings
func LeftPad2Len(s string, padStr string, overallLen int) string {
	var padCountInt int = 1 + ((overallLen - len(padStr)) / len(padStr))
	var retStr = strings.Repeat(padStr, padCountInt) + s
	return retStr[(len(retStr) - overallLen):]
}

// TrimSuffix removes a suffix from a string.
func TrimSuffix(s, suffix string) string {
	return strings.TrimSuffix(s, suffix)
}

// PrintKiloBytes prints a number of bytes in kilobytes as a string
func PrintKiloBytes(bytes int64) string {

	var kilobytes float64 = float64(bytes / 1024)

	return fmt.Sprint(FloatToString(kilobytes, 2), " kB")
}

// PrintMegaBytes prints a number of bytes in megabytes as a human readable string
func PrintMegaBytes(bytes int64) string {

	var kilobytes float64 = float64(bytes / 1024)

	var megabytes float64 = kilobytes / 1024 // cast to type float64

	return fmt.Sprint(FloatToString(megabytes, 2), " MB")
}

// IsPrintable checks if a string is printable and returns true or false.
// If it is an emoji, it returns true even if non printables are also in the string
func IsPrintable(s string) bool {
	// if an emoji is in the string just return true.
	if gomoji.ContainsEmoji(s) {
		return true
	}

	for _, c := range s {
		notTabLineFeed := c != 10 && c != 13 && c != 9
		if notTabLineFeed && (c < 32 || c == 127 || c == 44014 || c == 5760 || c == 8192 || c == 8193 || c == 8194 || c == 8195 || c == 8196 || c == 8197 || c == 8198 || c == 8199 || c == 8200 || c == 8201 || c == 8202 || c == 8232 || c == 8233 || c == 8239 || c == 8287 || c == 8006 || c == 8007 || c == 8024 || c == 8026 || c == 8028 || c == 8030 || c == 8062 || c == 8063 || c == 8117 || c == 8148 || c == 8149 || c == 8156 || c == 8176 || c == 8177 || c == 8181 || c == 888 || c == 889 || c == 907 || c == 909 || c == 930 || c == 42125 || c == 7958 || c == 7959 || c == 7966 || c == 7967 || c == 42126 || c == 42127 || (896 <= c && c <= 899) || (8191 <= c && c <= 8207) || (8233 <= c && c <= 8239) || (8287 <= c && c <= 8303) || (8384 <= c && c <= 8399) || (6170 <= c && c <= 6175) || (6264 <= c && c <= 6271) || (6315 <= c && c <= 6319) || (1480 <= c && c <= 1487) || (1515 <= c && c <= 1519) || (1525 <= c && c <= 1535) || (42936 <= c && c <= 42998) || c >= 55204) {
			return false
		}
	}
	return true
}

// PrintZettaBytes prints a number of bytes in zettabytes as a human readable string
func PrintZettaBytes(bytes int64) string {

	var kilobytes float64 = float64(bytes / 1024)

	var megabytes float64 = (kilobytes / 1024) // cast to type float64

	var gigabytes float64 = (megabytes / 1024)

	var terabytes float64 = (gigabytes / 1024)

	var petabytes float64 = (terabytes / 1024)

	var exabytes float64 = (petabytes / 1024)

	var zettabytes float64 = (exabytes / 1024)

	return fmt.Sprint(FloatToString(zettabytes, 2), " ZB")
}

// FloatToString converts a float to a string
func FloatToString(input_num float64, precision int) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', precision, 64)
}

// Round rounds a float to a specified number of decimal places
func Round(x, unit float64) float64 {
	return float64(int64(x/unit+0.5)) * unit
}

// StringToInteger converts a string to an integer
func StringToInt(val string) int {

	r, err := strconv.Atoi(val)
	if err != nil {
		return 0
	}
	return r
}

// StringToUInt16 converts a string to an uin16
func StringToUInt16(val string) uint16 {

	r, err := strconv.Atoi(val)
	if err != nil {
		return 0
	}
	return uint16(r)
}

// StringToFloat converts a string to a float
func StringToFloat(val string, precision int) (r float64) {
	r, err := strconv.ParseFloat(val, precision)
	if err != nil {
		return 0
	}
	return r
}

// StringToUInt64 converts a string to an uin64
func StringToUInt64(val string) uint64 {
	i, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0
	}
	return i
}

// StringToUInt8 converts a string to an uint8
func StringToUInt8(val string) uint8 {
	i, err := strconv.ParseUint(val, 10, 8)
	if err != nil {
		return 0
	}
	return uint8(i)
}

// IntToString converts an integer to a string
func IntToString(val int) string {
	return strconv.Itoa(val)
}

// IntToBool converts an integer to a bool
func IntToBool(val int) bool {
	return val != 0
}

// Int32ToString converts a 32 bit integer to a string
func Int32ToString(val int32) string {
	return strconv.Itoa(Int32ToInt(val))
}

// Int64ToString converts a 64 bit integer to a string
func Int64ToString(val int64) string {
	return strconv.FormatInt(val, 10)
}

// Int64ToInt32 converts a 64 bit integer to a 32 bit integer
func Int64ToInt32(val int64) (ret int) {
	tempLong := ((val >> 32) << 32) //shift it right then left 32 bits, which zeroes the lower half of the long
	ret = (int)(val - tempLong)
	return ret
}

// Int32ToInt converts a 32 bit integer to a 32 bit integer
func Int32ToInt(val int32) (ret int) {
	return int(val)
}

// Int32ToDWORD converts an int32 to DWORD byte array.
func Int32ToDWORD(val int32) (value []byte) {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, val)
	if err != nil {
		return
	}
	value = buff.Bytes()
	return
}

// DWordToInt converts an 4 byte to a int32
func DWordToInt(data []byte) (num int32) {
	buff := bytes.NewReader(data)
	err := binary.Read(buff, binary.BigEndian, &num)
	if err != nil {
		return
	}

	return
}

// BoolToString converts a bool to a string
func BoolToString(val bool) string {
	return strconv.FormatBool(val)
}

// StringToBool converts a string to a bool
func StringToBool(val string) bool {
	r, err := strconv.ParseBool(val)
	if err != nil {
		return false
	}
	return r
}

// Init initilizes the version
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

// GenPackageImport a utility for initilizing a go package file
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

// MakeFirstLowerCase makes the first letter of a string lowercase
func MakeFirstLowerCase(s string) string {

	if len(s) < 2 {
		return strings.ToLower(s)
	}

	bts := []byte(s)

	lc := bytes.ToLower([]byte{bts[0]})
	rest := bts[1:]

	return string(bytes.Join([][]byte{lc, rest}, nil))
}

// ExtractArgsWithinBrackets extracts the arguments within brackets
func ExtractArgsWithinBrackets(str string) (res []string) {

	brackets := &unicode.RangeTable{
		R16: []unicode.Range16{
			// {0x0028, 0x0029, 1}, // ( )
			// {0x005b, 0x005d, 1}, // [ ]
			{0x007b, 0x007d, 1}, // { }
		},
	}

	isBracket := func(r rune) bool {
		return unicode.In(r, brackets)
	}

	res = strings.FieldsFunc(str, isBracket)
	return
}

// Random returns a random number between min and max
func Random(min, max int) int {
	// Create a new random source with the current time as the seed
	seed := time.Now().UTC().UnixNano()
	source := rand.NewSource(seed)
	random := rand.New(source)

	// Generate a random number between min and max
	return random.Intn(max-min+1) + min
}

// Between will return substring between two strings
func Between(value string, a string, b string) string {
	// Get substring between two strings.
	posFirst := strings.Index(value, a)
	if posFirst == -1 {
		return ""
	}
	posLast := strings.Index(value, b)
	if posLast == -1 {
		return ""
	}
	posFirstAdjusted := posFirst + len(a)
	if posFirstAdjusted >= posLast {
		return ""
	}
	return value[posFirstAdjusted:posLast]
}

// SyncMapAny will return true if there are any items in the sync.Map
func SyncMapAny(x *sync.Map) (ok bool) {
	x.Range(func(key interface{}, value interface{}) bool {
		ok = true
		return false
	})
	return
}

// SyncMapLength will return true the length of items in the sync.Map
func SyncMapLength(x *sync.Map) (length int) {
	x.Range(func(key interface{}, value interface{}) bool {
		length = length + 1
		return true
	})
	return
}

// NewUUID generates a random UUID according to RFC 4122
func NewUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(cryptoRand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

// RandomString returns a random string of length
func RandomString(strlen int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

// BigEndianUInt16 will convert a uint16 to 2 bytes
func BigEndianUInt16(value uint16) (value1 byte, value2 byte) {
	inject := make([]byte, 2)
	binary.BigEndian.PutUint16(inject, value)
	value1 = inject[0]
	value2 = inject[1]
	return
}

// BigEndianUInt32 will convert a uint32 to 4 bytes
func BigEndianUInt32(value uint32) (value1 byte, value2 byte, value3 byte, value4 byte) {
	inject := make([]byte, 4)
	binary.BigEndian.PutUint32(inject, value)
	value1 = inject[0]
	value2 = inject[1]
	value3 = inject[2]
	value4 = inject[3]
	return
}

// IpAddressToUint32 will convert a string ip address to UInt32
func IpAddressToUint32(ip string) (value uint32) {
	binary.Read(bytes.NewBuffer(net.ParseIP(ip).To4()), binary.BigEndian, &value)
	return
}

// GetUserAgentName will return the name of the user agent.
func GetUserAgentName(req *http.Request) (name string) {
	ua := req.Header.Get("User-Agent")
	if ua != "" {
		values := strings.Split(ua, " ")
		lastBlock := values[len(values)-1]
		details := strings.Split(lastBlock, "/")
		if len(details) > 0 {
			name = details[0]
		}
	}
	return
}
