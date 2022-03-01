package main

import (
	"log"
	"version-action/pkgs/version"
	"version-action/pkgs/utils"
)

var commitTypes = [...]string {"breaking", "feature", "bugfix"}
var bumps = map[string]string{"breaking": "major", "feature": "minor", "bugfix": "patch"}

type Commit struct {
	Tag string
	Commit string
}

func prepareTagCommit(commitMessage, environment string) Commit {
	commit := Commit{}
	// Get latest tag
	latestTagRaw, err := version.GetLatestTag(environment)
	if err != nil {
		log.Fatalln("Error was found while getting the latest commit message", err)
	}
	// Trim tag & Remove 'v'
	commit.Tag = version.TrimTag(latestTagRaw)

	// Get commit message version level (breaking, feature, bugfix)
	versionType := version.GetVersionType(commitMessage, commitTypes)
	if !versionType.Contains {
		log.Fatalln("Version tag is not valid. Commit message must contain one of the following: [breaking, feature, bugfix]")
	}
	commit.Commit = versionType.VersionType
	return commit
}

func stagingVersion(tag, verionType string, rc bool) string {
	log.Println("Building staging version...")
	// increase rc version by 1
	if rc {
		rcTag, err := version.IncreaseRc(tag)
		if err != nil {
			log.Fatalln("Error converting rc version to int:", err)
		}
		finalTag := version.AddV(rcTag)
		return finalTag
	// Bump version and add rc
	} else {
		semVer := version.MakeSemVer(tag)
		log.Println("succesfully made semver")
		bumped := version.Bump(bumps, verionType, semVer)
		// Add rc
		strSemver := version.SemVerToString(bumped)
		rcTag := version.AddRc(strSemver)
		// Restore v
		finalTag := version.AddV(rcTag)
		return finalTag
	}
}

func productionVersion(tag string, rc bool) string{
	log.Println("Building production version for tag:", tag)

	if rc {
		// Remove rc
		tagNoRc := version.RemoveSuffix(tag, ".rc-")
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
	commit := prepareTagCommit(commitMessage, environment)
	
	log.Println("Trimmed tag:", commit.Tag)
	log.Println("versionType:", commit.Commit)
	
	// Check if PreRelease exists
	rc := version.CheckRc(commit.Tag)
	
	// Calculate staging or production version
	switch environment {
	case "staging":
		finalTag := stagingVersion(commit.Tag, commit.Commit, rc)
		log.Println(finalTag)
	case "main", "master":
		finalTag := productionVersion(commit.Tag, rc)
		log.Println(finalTag)
	}
}