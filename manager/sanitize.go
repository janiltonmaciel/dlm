package manager

import (
	"fmt"
	"regexp"
	"strings"
)

type (
	SanitizeFunc func(line string, distroFrom Distribution, distro Distribution) string
	Sanitize     struct {
		distro Distribution
		funcs  []SanitizeFunc
	}
)

var (
	allRemoveRe    = regexp.MustCompile(`(^(\s+)?FROM(.*)|^(\s+)?CMD(.*)|^(\s+)?ENTRYPOINT(.*))`)
	allOpenSSLRe   = regexp.MustCompile(`\s+openssl-dev\s+`)
	allNoNetworkRe = regexp.MustCompile(`\s+--no-network\s+`)
	allAwkSystemRe = regexp.MustCompile(`(awk\s+\'system.*next.*so.*\')`)

	phpGpgRe        = regexp.MustCompile(`\s*gpg.*--keyserver.*--recv-keys.*\\`)
	phpCopySourceRe = regexp.MustCompile(`(\s*COPY.*docker-php-source.*/usr/local/bin/)`)
	phpCopyExtRe    = regexp.MustCompile(`(\s*COPY.*docker-php-ext.*/usr/local/bin/)`)

	pythonGpGRe = regexp.MustCompile(`.*&&.*gpg\s+--keyserver.*--recv-keys`)

	LanguageSanitizeDockerfile = map[string][]SanitizeFunc{
		"ALL": []SanitizeFunc{
			func(line string, _ Distribution, distro Distribution) string {
				if strings.ToUpper(distro.Name) == DistributionALpine {
					return allNoNetworkRe.ReplaceAllString(line, " ")
				}
				return line
			},

			func(line string, _ Distribution, _ Distribution) string {
				if match := allRemoveRe.MatchString(line); match {
					return ""
				}
				return line
			},

			func(line string, _ Distribution, distro Distribution) string {
				if strings.ToUpper(distro.Name) == DistributionALpine {
					return allAwkSystemRe.ReplaceAllString(line, " $1 | xargs -r apk info --installed | sort -u")
				}
				return line
			},

			// http://lists.alpinelinux.org/alpine-devel/5463.html
			func(line string, distroFrom Distribution, distro Distribution) string {
				if strings.ToUpper(distroFrom.Name) == DistributionALpine && distroFrom.Release >= float32(3.5) && distro.Release < float32(3.5) {
					return allOpenSSLRe.ReplaceAllString(line, " libressl-dev ")
				}
				return line
			},
		},

		"PHP": []SanitizeFunc{
			func(line string, _ Distribution, _ Distribution) string {
				if match := phpGpgRe.MatchString(line); match {
					return `		  gpg --batch --keyserver p80.pool.sks-keyservers.net --recv-keys "$key" || \
		  gpg --batch --keyserver ha.pool.sks-keyservers.net --recv-keys "$key" || \
		  gpg --batch --keyserver ipv4.pool.sks-keyservers.net --recv-keys "$key" || \
		  gpg --batch --keyserver pgp.mit.edu --recv-keys "$key" || \
		  gpg --batch --keyserver keyserver.pgp.com --recv-keys "$key"; \`
				}
				return line
			},

			func(line string, _ Distribution, distro Distribution) string {
				if match := phpCopySourceRe.MatchString(line); match {
					respository := strings.TrimSuffix(distro.UrlDockerfile, "/Dockerfile")
					runCmd := `RUN %s && \
						curl -fs -o /usr/local/bin/docker-php-source %s/docker-php-source && \
						chmod +x /usr/local/bin/docker-php-source %s`

					cmds := make(map[string]string)
					switch strings.ToUpper(distro.Name) {
					case DistributionALpine:
						cmds["INSTALL"] = `apk add --no-cache --virtual .fetch-deps curl`
						cmds["DELETE"] = ` && apk del .fetch-deps`
					default:
						cmds["INSTALL"] = `apt-get update && apt-get install -y --no-install-recommends curl`
						cmds["DELETE"] = ``
					}
					return fmt.Sprintf(runCmd, cmds["INSTALL"], respository, cmds["DELETE"])
				}
				return line
			},

			func(line string, _ Distribution, distro Distribution) string {
				if match := phpCopyExtRe.MatchString(line); match {
					respository := strings.TrimSuffix(distro.UrlDockerfile, "/Dockerfile")
					runCmd := `RUN %s && \
					curl -fs -o /usr/local/bin/docker-php-entrypoint %s/docker-php-entrypoint && \
					curl -fs -o /usr/local/bin/docker-php-ext-configure %s/docker-php-ext-configure && \
					curl -fs -o /usr/local/bin/docker-php-ext-enable %s/docker-php-ext-enable && \
					curl -fs -o /usr/local/bin/docker-php-ext-install %s/docker-php-ext-install && \
					chmod +x /usr/local/bin/docker-php-* %s`

					cmds := make(map[string]string)
					switch strings.ToUpper(distro.Name) {
					case DistributionALpine:
						cmds["INSTALL"] = `apk add --no-cache --virtual .fetch-deps curl`
						cmds["DELETE"] = `&& apk del .fetch-deps`
					default:
						cmds["INSTALL"] = `apt-get update && apt-get install -y --no-install-recommends curl`
						cmds["DELETE"] = ``
					}
					return fmt.Sprintf(runCmd, cmds["INSTALL"], respository, respository, respository, respository, cmds["DELETE"])
				}
				return line
			},
		},

		"PYTHON": []SanitizeFunc{
			func(line string, _ Distribution, _ Distribution) string {
				if match := pythonGpGRe.MatchString(line); match {
					return `  && gpg --keyserver p80.pool.sks-keyservers.net --recv-keys "$GPG_KEY" || \
		  gpg --keyserver ha.pool.sks-keyservers.net --recv-keys "$GPG_KEY" || \
		  gpg --keyserver ipv4.pool.sks-keyservers.net --recv-keys "$GPG_KEY" || \
		  gpg --keyserver pgp.mit.edu --recv-keys "$GPG_KEY" || \
		  gpg --keyserver keyserver.pgp.com --recv-keys "$GPG_KEY" \`
				}
				return line
			},
		},
	}
)

func NewSanitize(distro Distribution) Sanitize {
	funcs := getLanguageSanitizeDockerfile(distro.Language.Name)
	return Sanitize{
		distro: distro,
		funcs:  funcs,
	}
}

func (s Sanitize) Do(distributionFrom Distribution) (string, error) {
	content, err := GetUrl(s.distro.UrlDockerfile)
	if err != nil {
		return "", err
	}

	var data []string
	data, err = LinesFromReader(content)
	if err != nil {
		return "", err
	}

	newData := make([]string, 0)
	for _, line := range data {
		for _, sanitize := range s.funcs {
			line = sanitize(line, distributionFrom, s.distro)
		}
		newData = append(newData, line)
	}

	return strings.Join(newData, "\n"), nil
}

func getLanguageSanitizeDockerfile(languageName string) (data []SanitizeFunc) {
	data = append(data, LanguageSanitizeDockerfile["ALL"]...)
	lsd, found := LanguageSanitizeDockerfile[strings.ToUpper(languageName)]
	if found {
		data = append(data, lsd...)
		return
	}
	return
}
