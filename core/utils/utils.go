package utils

import (
	"io/ioutil"
	"log"
	"strings"
)

func ReplaceTokenInFile(file string, find string, replaceWith string) {
	input, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, find) {
			lines[i] = replaceWith
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(file, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
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

func Array(values ...string) []string {
	var out []string
	for _, value := range values {
		out = append(out, value)
	}
	return out
}
