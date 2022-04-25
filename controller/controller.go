package controller

import (
	"github.com/strategicpause/cgstat/stats"
	"github.com/strategicpause/cgstat/writer"
	"fmt"
	"time"
)

// Controller manages the different interactions between components.
type Controller interface {
	Start() error
}

// CgroupStatsProviderFn controls which set of CgroupStats are returned for a user request.
type CgroupStatsProviderFn func() ([]*stats.CgroupStats, error)

type CgroupStatsController struct {
	displayWriter writer.StatsWriter
	csvWriter writer.StatsWriter
	statsProviderFn CgroupStatsProviderFn
	followMode bool
	ticker *time.Ticker
}

func NewCgroupStatsController(args *stats.CgstatArgs) (Controller, error) {
	csvWriter, err := getCsvWriter(args)
	if err != nil {
		return nil, err
	}

	return &CgroupStatsController{
		displayWriter: getDisplayWriter(args),
		csvWriter: csvWriter,
		statsProviderFn: getStatsProvider(args),
		followMode: args.FollowMode,
		ticker: time.NewTicker(args.GetRefreshInterval()),
	}, nil
}

func getDisplayWriter(args *stats.CgstatArgs) writer.StatsWriter {
	if args.VerboseOutput {
		return writer.NewCgroupStatsVerboseWriter()
	}
	return writer.NewCgStatsDisplayWriter()
}

func getCsvWriter(args *stats.CgstatArgs) (writer.StatsWriter, error) {
	if args.HasOutputFile() {
		return writer.NewCgroupStatsCsvWriter(args)
	}
	return writer.NewNullWriter(), nil
}

func getStatsProvider(args *stats.CgstatArgs) CgroupStatsProviderFn {
	provider := stats.NewCgroupStatsProvider()

	if args.HasPrefix() {
		return func() ([]*stats.CgroupStats, error) {
			return provider.GetCgroupStatsByPrefix(args.CgroupPrefix)
		}
	}
	return func() ([]*stats.CgroupStats, error) {
		return provider.GetCgroupStatsByName(args.CgroupName)
	}
}

func (c *CgroupStatsController) Start() error {
	// Clear Screen
	fmt.Print("\033[H\033[2J")

	for range c.ticker.C {
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

func (c *CgroupStatsController) writeStats() error {
	cgroupStats, err := c.statsProviderFn()
	if err != nil {
		return err
	}
	err = c.displayWriter.Write(cgroupStats)
	if err != nil {
		return err
	}
	return c.csvWriter.Write(cgroupStats)
}
