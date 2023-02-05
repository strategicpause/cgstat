package list

import (
	"fmt"

	"github.com/strategicpause/cgstat/stats"

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
	provider := stats.NewCgroupStatsProvider()

	prefix := cCtx.String(ArgsPrefix)
	cgroups := provider.ListCgroupsByPrefix(prefix)

	for _, cgroup := range cgroups {
		fmt.Println(cgroup)
	}
	return nil
}
