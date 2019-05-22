package core

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

func init() {
	loadDistributionConfig()
	languageConfig := loadLanguageConfig()
	setLanguages(languageConfig)
	setVersions(languageChoices)
	// fmt.Printf("languageChoices: %+v", languageChoices)
}

type (
	Config interface {
		Parse(data []byte) error
	}
	Block struct {
		Description string
		Data        string
		Before      []string
		After       []string
	}

	Context struct {
		From      Block
		Before    []string
		Languages []Block
		After     []string
	}

	AnswerVersion struct {
		Language Language
		Version  Version
	}
	LanguageSanitize struct {
		Pattern *regexp.Regexp
		Replace func(d Distribution) string
	}

	Language struct {
		Name       string   `yaml:"language"`
		Alias      string   `yaml:"alias"`
		Url        string   `yaml:"url"`
		Sort       []string `yaml:"sort"`
		sortDistro map[string]int
	}

	Distribution struct {
		Name            string   `yaml:"name"`
		ReleaseName     string   `yaml:"releaseName"`
		Release         float32  `yaml:"release"`
		Image           string   `yaml:"image"`
		Weight          int      `yaml:"weight"`
		Tags            []string `yaml:"tags"`
		UrlRepository   string   `yaml:"urlRepository"`
		UrlDockerfile   string   `yaml:"urlDockerfile"`
		ImageRepository string   `yaml:"imageRepository"`
		Language        Language
	}

	Version struct {
		Version              string         `yaml:"version"`
		MajorVersion         string         `yaml:"majorVersion"`
		Prerelease           bool           `yaml:"prerelease"`
		Date                 string         `yaml:"date"`
		Current              bool           `yaml:"current"`
		DistributionReleases string         `yaml:"distributionReleases"`
		Distributions        []Distribution `yaml:"distributions"`
	}

	LanguageConfig     []Language
	VersionConfig      []Version
	DistributionConfig []Distribution
)

const (
	DistributionDebian = "DEBIAN"
	DistributionUbuntu = "UBUNTU"
	DistributionALpine = "ALPINE"
)

var (
	languageChoices map[string]Language
	languages       []string
	distributions   []Distribution

	DisableGpgIPV6 = `
# DISABLE GPG IPV6
RUN mkdir ~/.gnupg && echo 'disable-ipv6' >> ~/.gnupg/dirmngr.conf`

	DistributionContext = map[string]map[string][]string{
		DistributionDebian: {
			"Before": []string{},
			"After": []string{
				`RUN apt-get update && apt-get install -y --no-install-recommends bash \
				&& rm -rf /var/lib/apt/lists/*`,
			},
		},
	}

	LanguageContext = map[string]map[string][]string{
		"NODE": {
			"Before": []string{
				"ENV PATH node_modules/.bin:$PATH",
			},
			"After": []string{},
		},
	}

	LanguageSanitizeDockerfile = map[string][]LanguageSanitize{
		"ALL": []LanguageSanitize{
			{Pattern: regexp.MustCompile(`(^(\s+)?FROM(.*)|^(\s+)?CMD(.*))`),
				Replace: func(d Distribution) string {
					return ""
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

func (lc *LanguageConfig) Parse(data []byte) error {
	return yaml.Unmarshal(data, lc)
}

func (vc *VersionConfig) Parse(data []byte) error {
	return yaml.Unmarshal(data, vc)
}

func (dc *DistributionConfig) Parse(data []byte) error {
	return yaml.Unmarshal(data, dc)
}

func (d Distribution) Description() string {
	desc := fmt.Sprintf(`##### %s ######
# Official Docker Image for %s
# repository: %s
# dockerfile: %s
# image: %s
# tag: %s`,
		strings.ToUpper(d.Language.Name), d.Language.Name, d.UrlRepository, d.UrlDockerfile, d.Image, d.ImageRepository)

	return desc
}

func (d Distribution) Sort() int {
	return d.Language.SortDistro(d.Name)
}

func (d Distribution) Hash() string {
	return fmt.Sprintf("%s-%s-%s", d.Language.Name, d.UrlRepository, d.UrlDockerfile)
}

func (l Language) SortDistro(distributionName string) int {
	return l.sortDistro[strings.ToLower(distributionName)]
}

func NewContext(distros []Distribution, distro Distribution) Context {
	distroContext := getDistributionContext(distro.Name)

	before := make([]string, 0)
	if len(distros) > 1 {
		before = append(before, DisableGpgIPV6)
	}
	before = append(before, distroContext["Before"]...)
	return Context{
		From: Block{
			Description: distro.Description(),
			Data:        fmt.Sprintf("FROM %s", distro.ImageRepository),
		},
		Before:    before,
		Languages: make([]Block, 0),
		After:     distroContext["After"],
	}
}

func NewLanguageBlock(distro Distribution, data string, isFrom bool) Block {
	langContext := getLanguageContext(distro.Language.Name)
	description := distro.Description()
	if isFrom {
		description = ""
	}
	block := Block{
		Description: description,
		Data:        data,
		Before:      langContext["Before"],
		After:       langContext["After"],
	}
	return block
}

func getDistributionContext(distributionName string) map[string][]string {
	distroContext, found := DistributionContext[strings.ToUpper(distributionName)]
	if found {
		return distroContext
	}
	distroContext = map[string][]string{
		"Before": []string{},
		"After":  []string{},
	}
	return distroContext
}

func getLanguageContext(languageName string) map[string][]string {
	langContext, found := LanguageContext[strings.ToUpper(languageName)]
	if found {
		return langContext
	}

	langContext = map[string][]string{
		"Before": []string{},
		"After":  []string{},
	}
	return langContext
}

func getLanguageSanitizeDockerfile(languageName string) (data []LanguageSanitize) {
	data = append(data, LanguageSanitizeDockerfile["ALL"]...)
	lsd, found := LanguageSanitizeDockerfile[strings.ToUpper(languageName)]
	if found {
		data = append(data, lsd...)
		return
	}
	return
}

func GetLanguages() []string {
	return languages
}
func GetDistributions() []Distribution {
	return distributions
}

func GetLanguage(alias string) *Language {
	if lang, found := languageChoices[strings.ToLower(alias)]; found {
		return &lang
	}
	return nil
}

func loadLanguageConfig() LanguageConfig {
	var languageConfig LanguageConfig
	loadConfig("languages.yml", &languageConfig)
	return languageConfig
}

func loadConfig(fileName string, config Config) {
	data, err := FindBytes(fileName)
	if err != nil {
		panic(err)
	}

	if err := config.Parse(data); err != nil {
		panic(err)
	}
}

func loadDistributionConfig() {
	var dc DistributionConfig
	loadConfig("distributions.yml", &dc)
	sort.Slice(dc, func(i, j int) bool {
		return dc[i].Sort() < dc[j].Sort()
	})
	distributions = dc
}

func setLanguages(config LanguageConfig) {
	languageChoices = make(map[string]Language)
	for i := 0; i < len(config); i++ {
		lang := config[i]
		if len(lang.Sort) > 0 {
			lang.sortDistro = make(map[string]int)
			for _, sort := range lang.Sort {
				data := strings.Split(sort, ":")
				key := data[0]
				value, _ := strconv.Atoi(data[1])
				lang.sortDistro[key] = value
			}
		}

		languageChoices[strings.ToLower(lang.Alias)] = lang
		languages = append(languages, lang.Alias)
	}

	sort.Strings(languages)
}
