package writer

import "github.com/strategicpause/cgstat/stats/common"

type StatsWriter interface {
	// Write will output the given CgroupStatsCollection to either the display or to a file.
	Write(cgroupStats common.CgroupStatsCollection) error
}

type ViewWriterOptions func() (StatsWriter, error)
