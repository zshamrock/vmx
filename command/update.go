package command

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"gopkg.in/urfave/cli.v1"
)

const versionURL = "https://raw.githubusercontent.com/zshamrock/vmx/master/version.txt"

func CheckUpdate(c *cli.Context) error {
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
		fmt.Println(err)
		// Same as above
		return nil
	}
	latestVersion := string(version)
	currentVersion := c.App.Version
	stale := isStale(currentVersion, latestVersion)
	if stale {
		fmt.Printf("===\n")
		fmt.Printf("* There is the latest version of the app available: current version %s, latest version: %s\n",
			currentVersion, latestVersion)
		fmt.Printf("===\n\n")
	}
	return nil
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

func splitVersion(version string) (int, int, int) {
	labels := strings.Split(version, ".")
	major, _ := strconv.Atoi(labels[0])
	minor, _ := strconv.Atoi(labels[1])
	patch, _ := strconv.Atoi(labels[2])
	return major, minor, patch
}
