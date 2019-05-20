package core

import (
	"fmt"
	"strings"
)

var (
	versions map[string]VersionConfig
)

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

func reverse(versions []Version) []Version {
	for i, j := 0, len(versions)-1; i < j; i, j = i+1, j-1 {
		versions[i], versions[j] = versions[j], versions[i]
	}
	return versions
}

func setVersions(languageChoices map[string]Language) {
	versions = make(map[string]VersionConfig)
	for _, lang := range languageChoices {
		versionConfig := VersionConfig{}
		fileNameVersions := fmt.Sprintf("%s-versions.yml", strings.ToLower(lang.Name))
		loadConfig(fileNameVersions, &versionConfig)

		for _, version := range versionConfig {
			for i := 0; i < len(version.Distributions); i++ {
				version.Distributions[i].Language = lang
			}
		}

		versions[strings.ToLower(lang.Name)] = versionConfig
	}
}
