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
	// Set ecr tag
	ecr_tag := version.RemovePrefix(value, "v")
	log.Println("Setting ecr tag as:", ecr_tag)
	fmt.Printf(`::set-output name=ecr_tag::%s`, ecr_tag)

	// Set repo tag
	log.Println("Setting repo tag as:", value)
	fmt.Printf(`::set-output name=tag::%s`, value)
}
