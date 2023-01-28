package list

import (
	"fmt"

	v1 "github.com/strategicpause/cgstat/stats/v1"
	"github.com/urfave/cli"
)

func Register() cli.Command {
	return cli.Command{
		Name:    "list",
		Aliases: []string{"l"},
		Usage:   "View a list of cgroups on the system.",
		Action:  action,
		Flags:   flags(),
	}
}

func action(cCtx *cli.Context) error {
	provider := v1.NewCgroupStatsProvider()

	prefix := cCtx.String(ArgsPrefix)
	cgroups := provider.ListCgroupsByPrefix(prefix)

	for _, cgroup := range cgroups {
		fmt.Println(cgroup)
	}
	return nil
}
