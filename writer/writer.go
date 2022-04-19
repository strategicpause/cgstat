package writer

import (
	"cgstat/stats"
)

type StatsWriter interface {
	Write(cgroupStats []*stats.CgroupStats) error
}