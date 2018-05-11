package command

import (
	"fmt"
	"sort"

	"gopkg.in/urfave/cli.v1"

	"github.com/zshamrock/vmx/config"
)

// CmdHosts lists available hosts
func CmdHosts(c *cli.Context) {
	CheckUpdate(c)
	hostsGroups := config.GetHostsGroups()
	groupNames := make([]string, 0, len(hostsGroups))
	for groupName := range hostsGroups {
		groupNames = append(groupNames, groupName)
	}
	sort.Strings(groupNames)
	for _, groupName := range groupNames {
		fmt.Printf("[%s]\n", groupName)
		hosts := hostsGroups[groupName]
		for _, host := range hosts {
			fmt.Println(host)
		}
		fmt.Println()
	}
}
