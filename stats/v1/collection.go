package v1

import "github.com/strategicpause/cgstat/stats/common"

type CgroupStatsCollection []common.CgroupStats

func (_ CgroupStatsCollection) GetCSVHeaders() []string {
	return []string{
		"Time", "Name", "UserCPU", "CurrentUsage", "MaxUsage", "UsageLimit", "RSS",
		"Cache", "Dirty", "WriteBack", "UnderOom", "OomKill",
	}
}
func (_ CgroupStatsCollection) GetDisplayHeaders() []interface{} {
	return []interface{}{
		"Name", "CPU", "NumProcesses", "CurrentUsage", "MaxUsage", "UsageLimit",
		"RSS", "Cache", "Dirty", "WriteBack", "UnderOom", "OomKill",
	}
}

func (c CgroupStatsCollection) GetCgroupStats() []common.CgroupStats {
	return c
}
