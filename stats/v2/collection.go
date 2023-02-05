package v2

import "github.com/strategicpause/cgstat/stats/common"

type CgroupStatsCollection []common.CgroupStats

func (_ CgroupStatsCollection) GetCSVHeaders() []string {
	return []string{}
}
func (_ CgroupStatsCollection) GetDisplayHeaders() []string {
	return []string{}
}

func (c CgroupStatsCollection) GetCgroupStats() []common.CgroupStats {
	return c
}
