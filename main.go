package main

import (
	"log"
	"version-action/pkgs/utils"
	"version-action/pkgs/version"
)

var commitTypes = [...]string{"breaking", "feature", "bugfix"}
var bumps = map[string]string{"breaking": "major", "feature": "minor", "bugfix": "patch"}

type Commit struct {
	Tag  string
	Type string
}

func prepareTagCommit(commitMessage, environment, bump string) Commit {
	// PATCH - Untill version 2.35.2 is supported on alpine
	exportOut, err := version.AddSafeDirectory()
	if err != nil {
		log.Println("AddSafeDirectory - Error was found while getting the latest tag")
	}
	log.Println("AddSafeDirectory - export_out:", exportOut)

	commit := Commit{}
	// Get latest tag
	latestTagRaw, err := version.GetLatestTag()
	if err != nil {
		log.Println("prepareTagCommit - Error was found while getting the latest tag")
	}
	// Trim tag & Remove 'v'
	commit.Tag = version.TrimTag(latestTagRaw)

	// Get commit message version level (breaking, feature, bugfix)
	versionType := version.GetVersionType(commitMessage, commitTypes, bump)
	if versionType == "" && environment == "staging" {
		log.Fatalln("Commit message must contain one of the following: [breaking, feature, bugfix]")
	}
	commit.Type = versionType
	return commit
}

func stagingVersion(commit Commit, rc bool) string {
	log.Println("Building staging version...")
	// increase rc version by 1
	if rc {
		rcTag, err := version.IncreaseRc(commit.Tag)
		if err != nil {
			log.Fatalln("Error converting rc version to int:", err)
		}
		finalTag := version.AddV(rcTag)
		return finalTag
	}
	// Bump version and add rc
	semVer := version.MakeSemVer(commit.Tag)
	log.Println("succesfully made semver")
	bumped := version.Bump(bumps, commit.Type, semVer)
	// Add rc
	strSemver := version.SemVerToString(bumped)
	rcTag := version.AddRc(strSemver)
	// Restore v
	finalTag := version.AddV(rcTag)
	return finalTag

}

func productionVersion(commit Commit, rc bool) string {
	log.Println("Building production version for tag:", commit.Tag)
	if rc {
		tagNoRc := version.RemoveSuffix(commit.Tag, "-rc.")
		return version.AddV(tagNoRc)
	}
	// Bump version
	semVer := version.MakeSemVer(commit.Tag)
	log.Println("succesfully made semver")
	bumped := version.Bump(bumps, "bugfix", semVer)
	strSemver := version.SemVerToString(bumped)
	// Restore v
	finalTag := version.AddV(strSemver)
	return finalTag
}

func sdkVersion(commit Commit) string {
	log.Println("Building sdk version for tag:", commit.Tag)
	// Bump version
	semVer := version.MakeSemVer(commit.Tag)
	log.Println("succesfully made semver")
	bumped := version.Bump(bumps, commit.Type, semVer)
	strSemver := version.SemVerToString(bumped)
	// Restore v
	finalTag := version.AddV(strSemver)
	return finalTag
}

func main() {
	// Getting os variables
	environment := utils.GetEnv("ENVIRONMENT")
	commitMessage := utils.GetEnv("COMMIT_MESSAGE")

	log.Println("Commit message:", commitMessage)

	// Get CLI arguments
	suffix := utils.GetEnv("INPUT_SUFFIX")
	bump := utils.GetEnv("INPUT_BUMP")

	if bump != "" && !utils.SliceContains([]string{"major", "minor", "patch"}, bump) {
		log.Fatalln("Error bump must be on of: major, minor, patch")
	}

	log.Println("suffix is:", suffix)

	// Preparing tag & commit
	commit := prepareTagCommit(commitMessage, environment, bump)

	log.Println("Trimmed tag:", commit.Tag)
	log.Println("versionType:", commit.Type)

	// Check if PreRelease exists
	rc := version.CheckRc(commit.Tag)
	// Calculate staging or production version
	switch environment {
	case "main", "master", "production":
		log.Println("Environment is production")
		finalTag := productionVersion(commit, rc)
		// Set repo & ecr tag
		utils.SetTagOutputName(finalTag)
	default:
		log.Printf("Branch is %s", environment)
		if suffix == "none" {
			finalTag := sdkVersion(commit)
			// Set repo & ecr tag
			utils.SetTagOutputName(finalTag)
		} else {
			finalTag := stagingVersion(commit, rc)
			// Set repo & ecr tag
			utils.SetTagOutputName(finalTag)
		}

	}
}
