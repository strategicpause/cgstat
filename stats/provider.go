package stats

import (
	"github.com/strategicpause/cgstat/stats/common"
	v1 "github.com/strategicpause/cgstat/stats/v1"
	v2 "github.com/strategicpause/cgstat/stats/v2"
)

func NewCgroupStatsProvider() common.CgroupStatsProvider {
	if isCgroupsV2Enabled() {
		return v2.NewCgroupStatsProvider()
	}
	return v1.NewCgroupStatsProvider()
}

func isCgroupsV2Enabled() bool {
	return true
}
