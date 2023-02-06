package v2

import "github.com/strategicpause/cgstat/stats/common"

type CgroupStatsCollection []common.CgroupStats

func (_ CgroupStatsCollection) GetCSVHeaders() []string {
	return []string{}
}
func (_ CgroupStatsCollection) GetDisplayHeaders() []interface{} {
	return []interface{}{
		"Name", "CPU Usage", "Throttled Periods", "PIDs", "Memory Usage", "Anon Memory",
		"File Memory", "Memory High", "Memory Max", "Memory Low", "OOM Kills",
	}
}

func (c CgroupStatsCollection) GetCgroupStats() []common.CgroupStats {
	return c
}
