package utils

import (
	"io/ioutil"
	"log"
	"math/rand"
	"os/exec"
	"strings"
	"time"

	"github.com/DanielRenne/GoCore/core/serverSettings"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func ReplaceTokenInFile(file string, find string, replaceWith string) {
	input, err := ioutil.ReadFile(file)
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
	err = ioutil.WriteFile(file, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

func TalkDirtyToMe(sayWhat string) {
	if serverSettings.WebConfig.Application.ReleaseMode == "development" {
		go exec.Command("say", sayWhat).Output()
	}
}

func TalkDirty(sayWhat string) {
	if serverSettings.WebConfig.Application.ReleaseMode == "development" {
		go exec.Command("say", sayWhat).Output()
	}
}

func TalkDirtySlowly(sayWhat string) {
	if serverSettings.WebConfig.Application.ReleaseMode == "development" {
		exec.Command("say", sayWhat).Output()
	}
}

func InArray(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func InIntArray(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func Array(values ...string) []string {
	var out []string
	for _, value := range values {
		out = append(out, value)
	}
	return out
}

func ArrayRemove(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func Dict(k string, v string) (ret map[string]string) {
	ret = make(map[string]string, 0)
	if k != "" && v != "" {
		ret[k] = v
	}
	return ret
}

func InterfaceMap() (ret map[string]interface{}) {
	ret = make(map[string]interface{}, 1)
	return ret
}

func RandomFloat() float32 {
	return rand.Float32() // Intn(max-min) + min
}
