package manager

import (
	"strings"
)

type (
	AnswerVersion struct {
		Language Language
		Version  Version
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
)

var versions map[string][]Version

func FindVersion(languageName string, versionTarget string) *Version {
	if versionConfig, found := versions[languageName]; found {
		for _, version := range versionConfig {
			if version.Version == versionTarget {
				return &version
			}
		}
	}
	return nil
}

func FindVersions(languageName string, withPrerelease bool, versionTarget string) []Version {
	data := make([]Version, 0)
	if versionConfig, found := versions[strings.ToLower(languageName)]; found {
		for _, version := range versionConfig {
			if versionTarget != "" && !strings.HasPrefix(version.Version, versionTarget) {
				continue
			}

			if withPrerelease {
				data = append(data, version)
			} else {
				if !version.Prerelease {
					data = append(data, version)
				}
			}
		}
	}
	return reverse(data)
}

func GetValueDefault(lang Language) string {
	valueDefault := ""
	switch strings.ToLower(lang.Name) {
	case "node":
		valueDefault = "12.2.0"
		// valueDefault = "11.4.0"
	case "python":
		valueDefault = "3.8.0a4"
	case "ruby":
		valueDefault = "2.6.0"
	case "golang":
		valueDefault = "1.12.0"
	case "swift":
		valueDefault = "5.0"
	}

	return valueDefault
}

func reverse(versions []Version) []Version {
	for i, j := 0, len(versions)-1; i < j; i, j = i+1, j-1 {
		versions[i], versions[j] = versions[j], versions[i]
	}
	return versions
}

func setVersions(lang Language, v []Version) {
	if versions == nil {
		versions = make(map[string][]Version)
	}
	for _, version := range v {
		for i := 0; i < len(version.Distributions); i++ {
			version.Distributions[i].Language = lang
		}
	}

	versions[strings.ToLower(lang.Name)] = v
}
