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

func SetTagOutputName(value string){
	// Set repo tag
	repo_tag := value
	log.Println("Setting repo tag as:", repo_tag)
	log.Printf("Running: ::set-output name=repo_tag::%s`", repo_tag)
	fmt.Printf(`::set-output name=repo_tag::%s`, repo_tag)
	
	// Set ecr tag
	ecr_tag := version.RemovePrefix(value, "v")
	log.Println("Setting ecr tag as:", ecr_tag)
	log.Printf("Running: ::set-output name=ecr_tag::%s`", ecr_tag)
	fmt.Printf(`::set-output name=ecr_tag::%s`, ecr_tag)


}
