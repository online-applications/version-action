package version

import (
	"fmt"
	"os/exec"
	"testing"
	"github.com/coreos/go-semver/semver"
)

func TestCheckRc(t *testing.T) {
	got := CheckRc("v1.0.0.rc-2")
	want := true

	if got != want {
		t.Errorf("got %t, wanted %t", got, want)
	}

	got_two := CheckRc("v1.0.0")
	want_two := false

	if got_two != want_two {
		t.Errorf("got %t, wanted %t", got, want)
	}

}

func TestGetLatestTag(t *testing.T) {
	got, err := GetLatestTag("main")
	want_byte, err_two := exec.Command("sh", "-c", "git tag  -l --merged main --sort='-*authordate' | head -n1").Output()
	want := string(want_byte)
	if err != nil || err_two != nil {
		t.Error("Error was found while getting the latest commit message", err)
	}

	if got != want {
		t.Errorf("got %s, wanted %s", got, want)
	}

	got_three, err_three := GetLatestTag("main")
	want_byte_two, err_four := exec.Command("sh", "-c", "git tag  -l --merged main --sort='-*authordate' | head -n1").Output()
	want_two := string(want_byte_two)
	if err_three != nil || err_four != nil {
		t.Error("Error was found while getting the latest commit message", err)
	}

	if got_three != want_two {
		t.Errorf("got %s, wanted %s", got, want)
	}
}

func TestTrimTag(t *testing.T) {
	got := TrimTag("v1.0.0.rc-2")
	want := "1.0.0.rc-2"

	if got != want {
		t.Errorf("got %s, wanted %s", got, want)
	}

	got_two := TrimTag("v1.0.0")
	want_two := "1.0.0"

	if got_two != want_two {
		t.Errorf("got %s, wanted %s", got, want)
	}
}

func TestGetVersionType(t *testing.T) {
	var commitTypes = [...]string{"breaking", "feature", "bugfix"}
	// breaking test
	commitMessage := "breaking - this is a new feature"

	contains, versionType := GetVersionType(commitMessage, commitTypes)
	contains_want := true
	versionType_want := "breaking"

	if contains != contains_want {
		t.Errorf("got %t, wanted %t", contains, contains_want)
	}

	if versionType != versionType_want {
		t.Errorf("got %s, wanted %s", versionType, versionType_want)
	}
	// feature test

	commitMessage_two := "feature this is a bugfix"

	contains_two, versionType_two := GetVersionType(commitMessage_two, commitTypes)
	contains_want_two := true
	versionType_want_two := "feature"

	if contains_two != contains_want_two {
		t.Errorf("got %t, wanted %t", contains_two, contains_want_two)
	}

	if versionType_two != versionType_want_two {
		t.Errorf("got %s, wanted %s", versionType_two, versionType_want_two)
	}
	// bugfix test
	commitMessage_three := "this is a bugfix"

	contains_three, versionType_three := GetVersionType(commitMessage_three, commitTypes)
	contains_want_three := true
	versionType_want_three := "bugfix"

	if contains_three != contains_want_three {
		t.Errorf("got %t, wanted %t", contains_three, contains_want_three)
	}

	if versionType_three != versionType_want_three {
		t.Errorf("got %s, wanted %s", versionType_three, versionType_want_three)
	}

}

func TestSemVerToString(t *testing.T) {
	semver := semver.New("1.2.0")
	got := SemVerToString(semver)
	want := "1.2.0"

	if got != want {
		t.Errorf("got %s, wanted %s", got, want)
	}
}

func TestRemoveV(t *testing.T) {
	got := RemoveV("v5.2.1")
	want := "5.2.1"

	if got != want {
		t.Errorf("got %s, wanted %s", got, want)
	}
}

func TestAddV(t *testing.T) {
	got := AddV("5.2.1")
	want := "v5.2.1"

	if got != want {
		t.Errorf("got %s, wanted %s", got, want)
	}
}

func TestMakeSemVer(t *testing.T) {
	got := MakeSemVer("5.2.1")
	var s interface{} = got
	if _, ok := s.(*semver.Version); ok {
		fmt.Println("This is an int")
	} else {
		t.Errorf("got %s, wanted *semver.Version", s)
	}
}
func TestRemoverc(t *testing.T) {
	got := Removerc("v5.2.1.rc-3")
	want := "v5.2.1"

	if got != want {
		t.Errorf("got %s, wanted %s", got, want)
	}
	got_two := Removerc("5.2.1.rc-3")
	want_two := "5.2.1"

	if got_two != want_two {
		t.Errorf("got %s, wanted %s", got, want)
	}
}

func TestIncreaseRc(t *testing.T) {
	got := IncreaseRc("v5.2.1.rc-3")
	want := "v5.2.1.rc-4"

	if got != want {
		t.Errorf("got %s, wanted %s", got, want)
	}
}

func TestBump(t *testing.T) {
	var bumps = map[string]string{"breaking": "major", "feature": "minor", "bugfix": "patch"}

	// major
	bumped := Bump(bumps, "breaking", MakeSemVer("2.5.0"))
	got := SemVerToString(bumped)

	want := "3.0.0"

	if got != want {
		t.Errorf("got %s, wanted %s", got, want)
	}

	// minor
	bumped_two := Bump(bumps, "feature", MakeSemVer("2.5.0"))
	got_two := SemVerToString(bumped_two)

	want_two := "2.6.0"

	if got_two != want_two {
		t.Errorf("got %s, wanted %s", got_two, want_two)
	}

	// patch
	bumped_three := Bump(bumps, "bugfix", MakeSemVer("2.5.0"))
	got_three := SemVerToString(bumped_three)

	want_three := "2.5.1"

	if got_three != want_three {
		t.Errorf("got %s, wanted %s", got_three, want_three)
	}

}

func TestRcVersionHandler(t *testing.T) {
	got := RcVersionHandler("v5.2.1.rc-3", true)
	want := "v5.2.1.rc-4"

	if got != want {
		t.Errorf("got %s, wanted %s", got, want)
	}

	got_two := RcVersionHandler("v5.2.1", false)
	want_two := "v5.2.1.rc-1"

	if got_two != want_two {
		t.Errorf("got %s, wanted %s", got_two, want_two)
	}
}