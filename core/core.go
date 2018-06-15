package core

import "os"

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
