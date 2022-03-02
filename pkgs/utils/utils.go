package utils

import (
	"os"
	"log"
	"fmt"
)

func GetEnv(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return ""
}

func SetTagOutputName(value string){
	log.Println("Setting tag as:", value)
	fmt.Printf(`::set-output name=tag::%s`, value)
}
