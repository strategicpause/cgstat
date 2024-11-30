package v2

import (
	"fmt"
	"io"
	"time"

	"github.com/rodaine/table"
	"github.com/strategicpause/cgstat/stats/common"
)

func NewCollection(stats []*CgroupStats) common.CgroupStatsCollection {
	return common.Collection[*CgroupStats]{
		Stats:                    stats,
		CsvHeadersProvider:       getCSVHeaders,
		CsvRowTransformer:        toCSVRow,
		DisplayHeadersProvider:   getDisplayHeaders,
		DisplayRowTransformer:    toDisplayRow,
		VerboseOutputTransformer: toVerboseOutput,
	}
}

func getCSVHeaders() []string {
	return []string{
		"Name", "Timestamp", "Throttled Periods", "Runnable Periods", "Current PIDs", "PID Limit", "Anon Memory Usage",
		"Kernel Memory", "Page Cache", "OOM Events", "OOM Kill Events", "TCP Sockets", "UDP Sockets", "Open Files",
	}
}

func toCSVRow(c *CgroupStats) []string {
	t, _ := time.Now().UTC().MarshalText()
	return []string{
		string(t),
		c.Name,
		fmt.Sprintf("%f", c.CPU.Utilization),
		fmt.Sprintf("%d", c.CPU.NumThrottledPeriods),
		fmt.Sprintf("%d", c.CPU.NumRunnablePeriods),
		fmt.Sprintf("%d", c.PID.Current),
		fmt.Sprintf("%d", c.PID.Limit),
		fmt.Sprintf("%d", c.Memory.Anon.Total),
		fmt.Sprintf("%d", c.Memory.Kernel.Slab+c.Memory.Kernel.Stack),
		fmt.Sprintf("%d", c.Memory.Filesystem.Active),
		fmt.Sprintf("%d", c.MemoryEvent.NumOomEvents),
		fmt.Sprintf("%d", c.MemoryEvent.NumOomKillEvents),
		fmt.Sprintf("%d", c.Network.TCPStats.Sockets),
		fmt.Sprintf("%d", c.Network.UDPStats.Sockets),
		fmt.Sprintf("%d", c.ProcStats.NumFD),
	}
}

func getDisplayHeaders() []interface{} {
	return []interface{}{
		"Name", "CPU Usage", "Throttled Periods", "PIDs", "Mem Usage", "Anon Mem", "Swap Mem", "File Mem", "Kernel Mem",
		"OOM Events / Kills", "TCP Sockets", "UDP Sockets", "Open Files",
	}
}

func toDisplayRow(c *CgroupStats) []interface{} {
	cgroupName := common.Shorten(c.Name, 32)
	cpuUsage := fmt.Sprintf("%.2f%%", c.CPU.Utilization)
	throttledPeriods := common.DisplayRatio(c.CPU.NumThrottledPeriods, c.CPU.NumRunnablePeriods)
	pids := common.DisplayRatio(c.PID.Current, c.PID.Limit)
	memUsage := common.DisplayRatio(c.Memory.Usage-c.Memory.Filesystem.Inactive, c.Memory.UsageLimit, common.WithBytes())
	anonMemUsage := common.FormatBytes(c.Memory.Anon.Total)
	swapMemUsage := common.FormatBytes(c.Memory.Swap.Usage)
	fileMemUsage := common.FormatBytes(c.Memory.Filesystem.Active + c.Memory.Filesystem.Inactive)
	kernelMemUsage := common.FormatBytes(c.Memory.Kernel.Slab + c.Memory.Kernel.Stack)
	numOomEvents := fmt.Sprintf("%d / %d", c.MemoryEvent.NumOomEvents, c.MemoryEvent.NumOomKillEvents)
	tcpSockets := fmt.Sprintf("%d (%s)", c.Network.TCPStats.Sockets, common.FormatBytes(c.Network.TCPStats.SocketMemory))
	udpSockets := fmt.Sprintf("%d (%s)", c.Network.UDPStats.Sockets, common.FormatBytes(c.Network.UDPStats.SocketMemory))
	numFDs := fmt.Sprintf("%d", c.ProcStats.NumFD)

	return []interface{}{
		cgroupName,
		cpuUsage,
		throttledPeriods,
		pids,
		memUsage,
		anonMemUsage,
		swapMemUsage,
		fileMemUsage,
		kernelMemUsage,
		numOomEvents,
		tcpSockets,
		udpSockets,
		numFDs,
	}
}

func toVerboseOutput(w io.Writer, c []*CgroupStats) {
	tbl := table.New()
	tbl.WithWriter(w)
	for _, cgroupStats := range c {
		tbl.AddRow("Name:", cgroupStats.Name)
		tbl.AddRow("PIDs:", common.DisplayRatio(cgroupStats.PID.Current, cgroupStats.PID.Limit, common.WithTotal()))
		tbl.AddRow("CPU Usage:", cgroupStats.CPU.Utilization)
		tbl.AddRow("Throttled Periods:", common.DisplayRatio(cgroupStats.CPU.NumThrottledPeriods, cgroupStats.CPU.NumRunnablePeriods, common.WithTotal()))
		tbl.AddRow("Throttled Time:", cgroupStats.CPU.ThrottledTimeInUsec)
		tbl.AddRow("System Usage", common.DisplayRatio(cgroupStats.CPU.SystemTimeInUsec, cgroupStats.CPU.UsageInUsec, common.WithTotal()))
		tbl.AddRow("User Usage", common.DisplayRatio(cgroupStats.CPU.UserTimeInUsec, cgroupStats.CPU.UsageInUsec, common.WithTotal()))
	}

	tbl.Print()
}
