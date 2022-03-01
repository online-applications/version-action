package main

import (
	"log"
	"version-action/pkgs/version"
	"version-action/pkgs/utils"
)

var commitTypes = [...]string {"breaking", "feature", "bugfix"}
var bumps = map[string]string{"breaking": "major", "feature": "minor", "bugfix": "patch"}


func prepareTagCommit(commitMessage, environment string) (string, string){
	// Get latest tag
	latestTagRaw, err := version.GetLatestTag(environment)
	if err != nil {
		log.Fatalln("Error was found while getting the latest commit message", err)
	}
	// Trim tag & Remove 'v'
	latest_tag := version.TrimTag(latestTagRaw)

	// Get commit message version level (breaking, feature, bugfix)
	contains, versionType := version.GetVersionType(commitMessage, commitTypes)
	if !contains {
		log.Fatalln("Version tag is not valid. Commit message must contain one of the following: [breaking, feature, bugfix]")
	}
	
	return latest_tag, versionType
}

func stagingVersion(tag, verionType string, rc bool) string {
	log.Println("Building staging version...")
	// increase rc version by 1
	if rc {
		rcTag := version.RcVersionHandler(tag, rc)
		finalTag := version.AddV(rcTag)
		return finalTag
	// Bump version and add rc
	} else {
		semVer := version.MakeSemVer(tag)
		log.Println("succesfully made semver")
		bumped := version.Bump(bumps, verionType, semVer)
		// Add rc
		strSemver := version.SemVerToString(bumped)
		rcTag := version.RcVersionHandler(strSemver, rc)
		// Restore v
		finalTag := version.AddV(rcTag)
		return finalTag
	}
}

func productionVersion(tag string, rc bool) string{
	log.Println("Building production version for tag:", tag)

	if rc {
		// Remove rc
		tagNoRc := version.Removerc(tag)
		// Restore v
		finalTag := version.AddV(tagNoRc)
		return finalTag
	} else {
		return version.AddV(tag)
	}
}

func main() {
	// Getting os variables
	environment 		:= utils.GetEnv("ENVIRONMENT")
	commitMessage 		:= utils.GetEnv("COMMIT_MESSAGE")
	
	log.Println("Commit message:", commitMessage)
	
	// Preparing tag & commit
	tag, versionType := prepareTagCommit(commitMessage, environment)
	
	log.Println("Trimmed tag:", tag)
	log.Println("versionType:", versionType)
	
	// Check if PreRelease exists
	rc := version.CheckRc(tag)
	
	// Calculate staging or production version
	switch environment {
	case "staging":
		finalTag := stagingVersion(tag, versionType, rc)
		log.Println(finalTag)
	case "main", "master":
		finalTag := productionVersion(tag, rc)
		log.Println(finalTag)
	}
}