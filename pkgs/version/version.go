package version

import (
	"log"
	"os/exec"
	"strings"
	"github.com/coreos/go-semver/semver"
	"strconv"
)

func CheckRc(s string) bool {
	isContains := strings.Contains(s, ".rc-")
	if !isContains {
		log.Println("Version doesnt contain rc")
		return false
	}
	log.Printf("Version: %s contain rc", s)
	return true
}

func GetLatestTag(environment string) (string, error){
	log.Println("Fetching latest tag version...")
	out, err := exec.Command("sh", "-c", "git describe --tags $(git rev-list --tags --max-count=1)").Output()
	if err != nil {
		log.Println("Error was found while getting the latest commit message", err)
	}
	log.Println("Fetched tag:", string(out))
	return string(out), err
}

func TrimTag(latestTagRaw string) string {
	log.Println("Trimming tag:", latestTagRaw)
	latest_tag := strings.Trim(latestTagRaw, "\n")
	latest := strings.Trim(latest_tag, " ")
	// After trimming - check if no previous tag exists, and return fallout tag
	if latest == "" {
		return "0.0.1"
	}
	latest_tag_no_v :=  RemovePrefix(latest, "v")
	return latest_tag_no_v
}

func GetVersionType(input string, words [3]string) string {
	v1 := strings.Contains(input, words[0])
	v2 := strings.Contains(input, words[1])
	v3 := strings.Contains(input, words[2])

	switch true {
	case v1: 
		return words[0]
	case v2:
		return words[1]
	case v3:
		return  words[2]
	}

	return ""
}

func SemVerToString(semVer *semver.Version) string {
	return semVer.String()
}

func RemovePrefix(tag, prefix string) string { 
	return strings.Trim(tag, prefix)
}

func AddV(tag string) string {
	return "v" + string(tag)
}

func MakeSemVer(tag string) *semver.Version {
	log.Printf("Coverting tag: %s to SemVer", tag)
	return semver.New(tag)
}

func RemoveSuffix(tag, suffix string) string {
	splitted := strings.Split(tag, suffix)
	return splitted[0]
}

func IncreaseRc(tag string) (string, error) {
	log.Println("Increasing rc version")
	// Extract rc
	splitted := strings.Split(tag, ".rc-")
	// Convert to int
	intV, err := strconv.Atoi(splitted[1])
	if err != nil {
		log.Println("Error converting rc version to int!")
	}
	// Increase by 1
	intVIncreased := intV +1
	// Convert to string
	strVIncreased := strconv.Itoa(intVIncreased)
	// return
	return splitted[0] + ".rc-" + strVIncreased, err
}

func Bump(bumps map[string]string, versionType string, semVer *semver.Version ) *semver.Version {
	bump, found := bumps[versionType]
	log.Println("Bumping", bump)
	switch found {
		case bump == "major":
			semVer.BumpMajor()
			return semVer
		case bump == "minor":
			 semVer.BumpMinor()
			 return semVer
		case bump == "patch":
			semVer.BumpPatch()
			return semVer
	}
	return semVer
}

func AddRc(tag string) string {
	log.Println("Adding .rc-1 to version")
	return tag + ".rc-1"
}