package command

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/zshamrock/vmx/config"
	"gopkg.in/urfave/cli.v1"
)

const (
	versionURL       = "https://raw.githubusercontent.com/zshamrock/vmx/master/version.txt"
	versionCheckFile = ".versioncheck"
)

func CheckUpdate(c *cli.Context) error {
	versionCheckFilePath := filepath.Join(config.DefaultConfig.Dir, versionCheckFile)
	content, err := ioutil.ReadFile(versionCheckFilePath)
	if err == nil {
		data := strings.Split(string(content), "\n")
		checkedAt, _ := time.Parse(time.Stamp, data[0])
		now := time.Now()
		if checkedAt.Month() == now.Month() && checkedAt.Day() == now.Day() {
			// Already checked today, then skip checking the latest version online
			version := data[1]
			if version != "" {
				compareVersionsAndNotify(version, c.App.Version)
			}
			return nil
		}
	}
	stampVersionCheckFile(versionCheckFilePath, "")
	resp, err := http.Get(versionURL)
	if err == nil {
		defer resp.Body.Close()
	}
	if err != nil || resp.StatusCode != http.StatusOK {

		// Ignore and don't bother the user about the problem, i.e. failure to check the latest version should not
		// prevent the app from functioning. As checking for the latest version is a nice to have feature, but not the
		// core functionality of the app. So the app should continue to work in that specific case as expected.
		return nil
	}
	version, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// Same as above
		return nil
	}
	latestVersion := string(version)
	compareVersionsAndNotify(latestVersion, c.App.Version)
	stampVersionCheckFile(versionCheckFilePath, latestVersion)
	return nil
}

func compareVersionsAndNotify(latestVersion, currentVersion string) {
	stale := isStale(currentVersion, latestVersion)
	if stale {
		fmt.Printf("===\n")
		fmt.Printf("* There is the latest version of the app available: current version %s, latest version: %s\n",
			currentVersion, latestVersion)
		fmt.Printf("===\n\n")
	}
}

func isStale(currentVersion, latestVersion string) bool {
	currentMajor, currentMinor, currentPatch := splitVersion(currentVersion)
	latestMajor, latestMinor, latestPatch := splitVersion(latestVersion)
	if latestMajor > currentMajor {
		return true
	}
	if latestMinor > currentMinor {
		return true
	}
	if latestPatch > currentPatch {
		return true
	}
	return false
}

func stampVersionCheckFile(versionCheckFilePath, version string) {
	stamp := time.Now().Format(time.Stamp)
	ioutil.WriteFile(versionCheckFilePath, []byte(stamp+"\n"+version), 0644)
}

func splitVersion(version string) (int, int, int) {
	labels := strings.Split(version, ".")
	major, _ := strconv.Atoi(labels[0])
	minor, _ := strconv.Atoi(labels[1])
	patch, _ := strconv.Atoi(labels[2])
	return major, minor, patch
}
