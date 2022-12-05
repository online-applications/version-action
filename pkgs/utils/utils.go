package utils

import (
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

	// Open GITHUB output file
	ghOutputFile, err := os.OpenFile(os.Getenv("GITHUB_OUTPUT"), os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	// Set ecr tag
	ecrTag := version.RemovePrefix(value, "v")
	log.Println("Setting ecr tag as: ", ecrTag)
	ecrTagOutput := "ecr_tag=" + ecrTag + "\n"
	_, err = ghOutputFile.WriteString(ecrTagOutput)
	if err != nil {
		panic(err)
	}

	// Set repo tag
	log.Println("Setting repo tag as: ", value)
	repoTagOutput := "repo_tag=" + value + "\n"
	_, err = ghOutputFile.WriteString(repoTagOutput)
	if err != nil {
		panic(err)
	}
}

func SliceContains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
