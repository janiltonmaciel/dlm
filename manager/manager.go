package manager

import (
	"bufio"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	reSlug = regexp.MustCompile("[^a-z0-9]+")
)

func loadConfig(fileName string, o interface{}) {
	data, err := FindBytes(fileName)
	if err != nil {
		panic(err)
	}

	if err := yaml.Unmarshal(data, o); err != nil {
		panic(err)
	}
}

func HasDockerfile() bool {
	if _, err := os.Stat("Dockerfile"); err == nil {
		return true
	} else {
		return false
	}
}

func GetUrl(url string) ([]string, error) {
	if val, err := cache.get(url); err == nil {
		// println("in cache...")
		return val, nil
	}

	content, err := request(url)
	if err != nil {
		return nil, err
	}

	// set in cache
	// println("not in cache...")
	_ = cache.set(content, url)

	return content, nil
}

func request(url string) ([]string, error) {
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

func slug(s string) string {
	return strings.Trim(reSlug.ReplaceAllString(strings.ToLower(s), "-"), "-")
}
