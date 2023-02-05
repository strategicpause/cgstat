package view

import (
	"fmt"
	"github.com/strategicpause/cgstat/stats"
	"github.com/strategicpause/cgstat/stats/common"
	"github.com/strategicpause/cgstat/writer"
	"time"

	"github.com/urfave/cli"
)

// CgroupStatsProviderFn controls which set of CgroupStats are returned for a user request.
type CgroupStatsProviderFn func() (common.CgroupStatsCollection, error)

type Command struct {
	writers         []writer.StatsWriter
	statsProviderFn CgroupStatsProviderFn
	followMode      bool
	ticker          *time.Ticker
}

func Register() cli.Command {
	return cli.Command{
		Name:        "view",
		Aliases:     []string{"v"},
		Usage:       "usage of view",
		UsageText:   "Usage text",
		Description: "Description",
		Action:      action,
		Flags:       flags(),
	}
}

func action(cCtx *cli.Context) error {
	viewArgs, err := parseArgs(cCtx)
	if err != nil {
		return err
	}

	cmd := Command{
		writers:         getWriters(viewArgs),
		statsProviderFn: getStatsProvider(viewArgs),
		followMode:      viewArgs.FollowMode,
		ticker:          time.NewTicker(viewArgs.GetRefreshInterval()),
	}
	return cmd.Run()
}

func getWriters(args *Args) []writer.StatsWriter {
	var options []writer.ViewWriterOptions

	if args.HasOutputFile() {
		options = append(options, writer.WithCSVWriter(args.OutputFile))
	}

	displayVerbosity := writer.Normal
	if args.VerboseOutput {
		displayVerbosity = writer.Verbose
	}
	options = append(options, writer.WithDisplayWriter(displayVerbosity))

	return writer.NewViewWriters(options)
}

func getStatsProvider(args *Args) CgroupStatsProviderFn {
	provider := stats.NewCgroupStatsProvider()

	if args.HasPrefix() {
		return func() (common.CgroupStatsCollection, error) {
			return provider.GetCgroupStatsByPrefix(args.CgroupPrefix)
		}
	}
	return func() (common.CgroupStatsCollection, error) {
		return provider.GetCgroupStatsByName(args.CgroupName)
	}
}

func (c *Command) Run() error {
	for range c.ticker.C {
		// Clear Screen
		fmt.Print("\033[H\033[2J")
		err := c.writeStats()
		if err != nil {
			return err
		}
		if !c.followMode {
			return nil
		}
	}
	return nil
}

func (c *Command) writeStats() error {
	cgroupStats, err := c.statsProviderFn()
	if err != nil {
		return err
	}
	for _, w := range c.writers {
		if err = w.Write(cgroupStats); err != nil {
			return err
		}
	}
	return nil
}
