package utils

import (
	"fmt"
	"log"
	"os"
	"version-action/pkgs/version"
)

func GetEnv(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return ""
}

func SetTagOutputName(value string) {
	// Set ecr tag
	ecrTag := version.RemovePrefix(value, "v")
	log.Println("Setting ecr tag as:", ecrTag)
	fmt.Printf(`::set-output name=ecr_tag::%s`, ecrTag)
	fmt.Print("\n")
	// Set repo tag
	repoTag := value
	log.Println("Setting repo tag as:", repoTag)
	fmt.Printf(`::set-output name=repo_tag::%s`, repoTag)
	fmt.Print("\n")
}

func SliceContains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
