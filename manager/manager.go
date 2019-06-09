package manager

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

var (
	reSlug = regexp.MustCompile("[^a-z0-9]+")
	ttlYML = 24 * time.Hour
)

func loadConfigUrl(url string, o interface{}) {
	if data, err := GetUrl(url, ttlYML); err == nil {
		_ = yaml.Unmarshal(data, o)
	}
}

func loadConfig(fileName string, o interface{}) {
	data, err := FindBytes(fileName)
	if err != nil {
		panic(err)
	}

	if err := yaml.Unmarshal(data, o); err != nil {
		panic(err)
	}
}

func GetUrl(url string, ttl ...time.Duration) ([]byte, error) {
	if val, err := cache.get(url, ttl...); err == nil {
		return val, nil
	}

	data, err := request(url)
	if err != nil {
		return nil, err
	}

	_ = cache.set(data, url)

	return data, nil
}

func request(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request error status code: %d", resp.StatusCode)
	}

	return ioutil.ReadAll(resp.Body)
}

func LinesFromReader(b []byte) ([]string, error) {
	var lines []string

	scanner := bufio.NewScanner(strings.NewReader(string(b)))
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

func HasDockerfile() bool {
	if _, err := os.Stat("Dockerfile"); err == nil {
		return true
	} else {
		return false
	}
}
