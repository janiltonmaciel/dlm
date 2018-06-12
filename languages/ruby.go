package languages

import "strings"

type Ruby struct {
}

func (r Ruby) Name() string {
	return "Ruby"
}

func (r Ruby) Help() string {
	return "https://www.ruby-lang.org/en/downloads/releases/"
}

func (r Ruby) Default() string {
	return "2.5.1"
}

func (r Ruby) GetDockerfile(version string) string {
	rubyVersion := r.rightPadVersion(version, ".0", 5)
	context := map[string]string{
		"RUBY_MAJOR":   rubyVersion[0:3],
		"RUBY_VERSION": rubyVersion,
	}

	return ParseTemplate(context, r.Name())
}

func (r Ruby) rightPadVersion(version string, padStr string, overallLen int) string {
	var padCountInt int
	padCountInt = 1 + ((overallLen - len(padStr)) / len(padStr))
	var retStr = version + strings.Repeat(padStr, padCountInt)
	return retStr[:overallLen]
}
