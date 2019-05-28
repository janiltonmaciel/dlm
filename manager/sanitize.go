package manager

import (
	"fmt"
	"regexp"
	"strings"
)

type (
	LanguageSanitize struct {
		Pattern *regexp.Regexp
		Replace func(d Distribution) string
	}
)

var (
	LanguageSanitizeDockerfile = map[string][]LanguageSanitize{
		"ALL": []LanguageSanitize{
			{Pattern: regexp.MustCompile(`(^(\s+)?FROM(.*)|^(\s+)?CMD(.*))`),
				Replace: func(d Distribution) string {
					return ""
				},
			},

			// http://lists.alpinelinux.org/alpine-devel/5463.html
			{Pattern: regexp.MustCompile(`\s+openssl-dev\s+`),
				Replace: func(d Distribution) string {
					if strings.ToUpper(d.Name) == DistributionALpine && d.Release >= float32(3.5) {
						return "libressl-dev"
					}
					return "openssl-dev"
				},
			},
		},

		"PHP": []LanguageSanitize{
			{Pattern: regexp.MustCompile(`\s*gpg.*--keyserver.*--recv-keys.*\\`),
				Replace: func(d Distribution) string {
					return `  gpg --batch --keyserver p80.pool.sks-keyservers.net --recv-keys "$key" || \
		  gpg --batch --keyserver ha.pool.sks-keyservers.net --recv-keys "$key" || \
		  gpg --batch --keyserver ipv4.pool.sks-keyservers.net --recv-keys "$key" || \
		  gpg --batch --keyserver pgp.mit.edu --recv-keys "$key" || \
		  gpg --batch --keyserver keyserver.pgp.com --recv-keys "$key"; \`
				},
			},

			{Pattern: regexp.MustCompile(`(\s*COPY.*docker-php-source.*/usr/local/bin/)`),
				Replace: func(d Distribution) string {
					respository := strings.TrimSuffix(d.UrlDockerfile, "/Dockerfile")
					runCmd := `RUN %s && \
						curl -fs -o /usr/local/bin/docker-php-source %s/docker-php-source && \
						chmod +x /usr/local/bin/docker-php-source && \
						%s`
					cmds := make(map[string]string)
					switch strings.ToUpper(d.Name) {
					case DistributionALpine:
						cmds["ADD"] = `apk add --no-cache --virtual .fetch-deps curl`
						cmds["DEL"] = `apk del .fetch-deps`
					default:
						cmds["ADD"] = `apt-get update && apt-get install -y --no-install-recommends curl`
						cmds["DEL"] = ``
					}
					return fmt.Sprintf(runCmd, cmds["ADD"], respository, cmds["DEL"])
				},
			},

			{Pattern: regexp.MustCompile(`(\s*COPY.*docker-php-ext.*/usr/local/bin/)`),
				Replace: func(d Distribution) string {
					respository := strings.TrimSuffix(d.UrlDockerfile, "/Dockerfile")
					runCmd := `RUN %s && \
					curl -fs -o /usr/local/bin/docker-php-entrypoint %s/docker-php-entrypoint && \
					curl -fs -o /usr/local/bin/docker-php-ext-configure %s/docker-php-ext-configure && \
					curl -fs -o /usr/local/bin/docker-php-ext-enable %s/docker-php-ext-enable && \
					curl -fs -o /usr/local/bin/docker-php-ext-install %s/docker-php-ext-install && \
					chmod +x /usr/local/bin/docker-php-* && \
						%s`
					cmds := make(map[string]string)
					switch strings.ToUpper(d.Name) {
					case DistributionALpine:
						cmds["ADD"] = `apk add --no-cache --virtual .fetch-deps curl`
						cmds["DEL"] = `apk del .fetch-deps`
					default:
						cmds["ADD"] = `apt-get update && apt-get install -y --no-install-recommends curl`
						cmds["DEL"] = ``
					}
					return fmt.Sprintf(runCmd, cmds["ADD"], respository, respository, respository, respository, cmds["DEL"])
				},
			},
		},

		"PYTHON": []LanguageSanitize{
			{Pattern: regexp.MustCompile(`.*&&.*gpg\s+--keyserver.*--recv-keys`),
				Replace: func(d Distribution) string {
					return `  && gpg --keyserver p80.pool.sks-keyservers.net --recv-keys "$GPG_KEY" || \
		  gpg --keyserver ha.pool.sks-keyservers.net --recv-keys "$GPG_KEY" || \
		  gpg --keyserver ipv4.pool.sks-keyservers.net --recv-keys "$GPG_KEY" || \
		  gpg --keyserver pgp.mit.edu --recv-keys "$GPG_KEY" || \
		  gpg --keyserver keyserver.pgp.com --recv-keys "$GPG_KEY" \`
				},
			},
		},
	}
)

func GetLanguageSanitizeDockerfile(languageName string) (data []LanguageSanitize) {
	data = append(data, LanguageSanitizeDockerfile["ALL"]...)
	lsd, found := LanguageSanitizeDockerfile[strings.ToUpper(languageName)]
	if found {
		data = append(data, lsd...)
		return
	}
	return
}
