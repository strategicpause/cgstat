package writer

import (
	"github.com/strategicpause/cgstat/stats"
)

type StatsWriter interface {
	Write(cgroupStats []*stats.CgroupStats) error
}