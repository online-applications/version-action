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
	arg1 := "--sort='-*authordate' | head -n1"
	out, err := exec.Command("sh", "-c", "git tag -l --merged", environment, arg1).Output()
	if err != nil {
		log.Println("Error was found while getting the latest commit message", err)
		log.Println("Tried to execute the command:", "sh -c git tag -l --merged", environment, arg1)
	}
	log.Println("Fetched tag:", string(out))
	return string(out), err
}

func TrimTag(latestTagRaw string) string {
	log.Println("Trimming tag")
	latest_tag := strings.Trim(latestTagRaw, "\n")
	latest := strings.Trim(latest_tag, " ")
	latest_tag_no_v :=  RemoveV(latest)
	return latest_tag_no_v
}

func GetVersionType(input string, words [3]string) (bool, string) {
    
	v1 := strings.Count(input, words[0])
	v2 := strings.Count(input, words[1])
	v3 := strings.Count(input, words[2])

	switch true {
	case v1 >= 1:
		return true, words[0]
	case v2 >= 1:
		return true, words[1]
	case v3 >= 1:
		return true, words[2]
	}
	return false, ""
}

func SemVerToString(semVer *semver.Version) string {
	return semVer.String()
}

func RemoveV(tag string) string {
	return strings.Trim(tag, "v")
}

func AddV(tag string) string {
	return "v" + string(tag)
}

func MakeSemVer(tag string) *semver.Version {
	log.Printf("Coverting tag: %s to SemVer", tag)
	return semver.New(tag)
}

func Removerc(tag string) string {
	splitted := strings.Split(tag, ".rc-")
	return splitted[0]
}

func IncreaseRc(tag string) string {
	// Extract rc
	splitted := strings.Split(tag, ".rc-")
	// Convert to int
	intV, err := strconv.Atoi(splitted[1])
	if err != nil {
		log.Fatalf("error convertiong rc version to int")
	}
	// Increase by 1
	intVIncreased := intV +1
	// Convert to string
	strVIncreased := strconv.Itoa(intVIncreased)
	// return
	return splitted[0] + ".rc-" + strVIncreased
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

func RcVersionHandler(tag string, rc bool) (string) {
	// Check if rc exists
	if rc {
		log.Println("Increasing rc version")
		return IncreaseRc(tag)
	} else {
		log.Println("Adding .rc-1 to version")
		return tag + ".rc-1"
	}
}