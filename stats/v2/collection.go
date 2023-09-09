package v2

import "github.com/strategicpause/cgstat/stats/common"

type CgroupStatsCollection []common.CgroupStats

func (_ CgroupStatsCollection) GetCSVHeaders() []string {
	return []string{}
}
func (_ CgroupStatsCollection) GetDisplayHeaders() []interface{} {
	return []interface{}{
		"Name", "CPU Usage", "Throttled Periods", "PIDs", "Anon Memory Usage", "Kernel Memory", "Page Cache",
		"OOM Kills", "TCP Sockets", "UDP Sockets", "Open Files",
	}
}

func (c CgroupStatsCollection) GetCgroupStats() []common.CgroupStats {
	return c
}
