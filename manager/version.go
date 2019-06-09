package manager

import (
	"reflect"
	"strings"

	"gopkg.in/AlecAivazis/survey.v1"
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
	if versionConfig, found := versions[strings.ToLower(languageName)]; found {
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

func FindMathVersion(languageName string, versionTarget string) string {
	versions := FindVersions(languageName, false, versionTarget)
	qtdVersions := len(versions)

	switch qtdVersions {
	case 0:
		return versionTarget
	case 1:
		return versions[0].Version
	default:
		return versions[qtdVersions-1].Version
	}
}

func TransformFindVersion(languageName string) survey.Transformer {
	return func(ans interface{}) interface{} {
		if isZero(reflect.ValueOf(ans)) {
			return nil
		}

		s, ok := ans.(string)
		if !ok {
			return nil
		}

		return FindMathVersion(languageName, s)
	}
}

func GetValueDefault(lang Language) string {
	valueDefault := ""
	switch strings.ToLower(lang.Name) {
	case "php":
		valueDefault = "7.3.5-cli"
	case "node":
		valueDefault = "12.2.0"
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

func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Slice, reflect.Map:
		return v.Len() == 0
	}

	// compare the types directly with more general coverage
	return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
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
