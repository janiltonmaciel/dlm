package core

import (
	"bufio"
	"io"
	"net/http"
	"os"
	"strings"
)

func CheckErr(err error) {
	if err != nil {
		println("\n> Good bye!\n")
		os.Exit(2)
	}
}

func HasDockerfile() bool {
	if _, err := os.Stat("Dockerfile"); err == nil {
		return true
	} else {
		return false
	}
}

func SanitizeDockerfile(distribution Distribution) (string, error) {
	data, err := GetUrl(distribution.UrlDockerfile)
	if err != nil {
		return "", err
	}

	sanitizeAll := getLanguageSanitizeDockerfile(distribution.Language.Name)
	newData := make([]string, 0)
	for _, line := range data {
		for _, sanitize := range sanitizeAll {
			if result := sanitize.Pattern.MatchString(line); result {
				line = sanitize.Replace(distribution)
				break
			}
		}
		newData = append(newData, line)
	}
	return strings.Join(newData, "\n"), nil
}

func GetUrl(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil
	}

	return LinesFromReader(resp.Body)
}

func LinesFromReader(r io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func Intersection(a, b []Distribution) (c []Distribution) {
	m := make(map[string]Distribution)
	for _, item := range a {
		m[item.Name] = item
	}

	for _, item := range b {
		if _, ok := m[item.Name]; ok {
			c = append(c, item)
		}
	}
	return
}
