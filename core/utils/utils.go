// Package utils provides a set of helper functions for your application
package utils

import (
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/DanielRenne/GoCore/core/serverSettings"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// RandStringRunes returns a random string of length n
func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// ReplaceTokenInFile replaces a token in a file with a find/replace value
func ReplaceTokenInFile(file string, find string, replaceWith string) {
	input, err := os.ReadFile(file)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, find) {
			lines[i] = strings.Replace(lines[i], find, replaceWith, -1)
		}
	}
	output := strings.Join(lines, "\n")
	err = os.WriteFile(file, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

// TalkDirtyToMe is a function that uses say command (linux or mac only if installed) to talk to you in development while your code is running
func TalkDirtyToMe(sayWhat string) {
	if serverSettings.WebConfig.Application.ReleaseMode == "development" {
		go exec.Command("say", sayWhat).Output()
	}
}

// TalkDirty is a function that uses say command (linux or mac only if installed) to talk to you in development while your code is running
func TalkDirty(sayWhat string) {
	if serverSettings.WebConfig.Application.ReleaseMode == "development" {
		go exec.Command("say", sayWhat).Output()
	}
}

// TalkDirtySlowly is a function that uses say (linux or mac only if installed) command to talk to you in development while your code is running, but will wait for the previous say to finish before starting the next line of code
func TalkDirtySlowly(sayWhat string) {
	if serverSettings.WebConfig.Application.ReleaseMode == "development" {
		exec.Command("say", sayWhat).Output()
	}
}

// InArray returns true if the string is in the slice of strings
func InArray(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// InIntArray returns true if the int is in the slice of ints
func InIntArray(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// Array passes in unlimited amount of strings and returns a slice of strings
func Array(values ...string) []string {
	var out []string
	out = append(out, values...)
	return out
}

// ArrayRemove removes an item from a slice of strings
func ArrayRemove(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

// Dict is a function that returns a map of strings
func Dict(k string, v string) (ret map[string]string) {
	ret = make(map[string]string, 0)
	if k != "" && v != "" {
		ret[k] = v
	}
	return ret
}

// InterfaceMap is a function that returns a map of interfaces
func InterfaceMap() (ret map[string]interface{}) {
	ret = make(map[string]interface{}, 1)
	return ret
}

// RandomFloat returns a random float
func RandomFloat() float32 {
	return rand.Float32()
}
