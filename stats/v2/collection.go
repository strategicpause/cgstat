package v2

import (
	"fmt"
	"github.com/rodaine/table"
	"github.com/strategicpause/cgstat/stats/common"
	"io"
	"time"
)

type CgroupStatsCollection []*CgroupStats

func (c CgroupStatsCollection) ToCsvOutput() *common.CsvOutput {
	csvOutput := common.CsvOutput{
		Headers: c.getCSVHeaders(),
	}

	for _, s := range c {
		csvOutput.Rows = append(csvOutput.Rows, c.toCSVRow(s))
	}

	return &csvOutput
}

func (_ CgroupStatsCollection) getCSVHeaders() []string {
	return []string{
		"Name", "Timestamp", "Throttled Periods", "Runnable Periods", "Current PIDs", "PID Limit", "Anon Memory Usage",
		"Kernel Memory", "Page Cache", "OOM Events", "OOM Kill Events", "TCP Sockets", "UDP Sockets", "Open Files",
	}
}

func (_ CgroupStatsCollection) toCSVRow(c *CgroupStats) []string {
	t, _ := time.Now().UTC().MarshalText()
	return []string{
		string(t),
		c.Name,
		fmt.Sprintf("%f", c.CPU.Usage),
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

func (c CgroupStatsCollection) ToDisplayOutput() *common.DisplayOutput {
	displayOutput := common.DisplayOutput{
		Headers: c.getDisplayHeaders(),
	}

	for _, s := range c {
		displayOutput.Rows = append(displayOutput.Rows, c.toDisplayRow(s))
	}

	return &displayOutput
}

func (_ CgroupStatsCollection) getDisplayHeaders() []interface{} {
	return []interface{}{
		"Name", "CPU Usage", "Throttled Periods", "PIDs", "Anon Memory Usage", "Kernel Memory", "Page Cache",
		"OOM Kills", "TCP Sockets", "UDP Sockets", "Open Files",
	}
}

func (_ CgroupStatsCollection) toDisplayRow(c *CgroupStats) []interface{} {
	cgroupName := common.Shorten(c.Name, 32)
	cpuUsage := fmt.Sprintf("%.2f%%", c.CPU.Usage)
	throttledPeriods := common.DisplayRatio(c.CPU.NumThrottledPeriods, c.CPU.NumRunnablePeriods)
	pids := common.DisplayRatio(c.PID.Current, c.PID.Limit)
	memoryUsage := common.DisplayRatio(c.Memory.Anon.Total, c.Memory.UsageLimit, common.WithBytes())
	kernelMemory := common.DisplayRatio(c.Memory.Kernel.Slab+c.Memory.Kernel.Stack, c.Memory.UsageLimit, common.WithBytes())
	pageCache := common.DisplayRatio(c.Memory.Filesystem.Active, c.Memory.UsageLimit, common.WithBytes())
	numOomKillEvents := fmt.Sprintf("%d", c.MemoryEvent.NumOomKillEvents)
	tcpSockets := fmt.Sprintf("%d (%s)", c.Network.TCPStats.Sockets, common.FormatBytes(c.Network.TCPStats.SocketMemory))
	udpSockets := fmt.Sprintf("%d (%s)", c.Network.UDPStats.Sockets, common.FormatBytes(c.Network.UDPStats.SocketMemory))
	numFDs := fmt.Sprintf("%d", c.ProcStats.NumFD)

	return []interface{}{
		cgroupName,
		cpuUsage,
		throttledPeriods,
		pids,
		memoryUsage,
		kernelMemory,
		pageCache,
		numOomKillEvents,
		tcpSockets,
		udpSockets,
		numFDs,
	}
}

func (c CgroupStatsCollection) ToVerboseOutput(w io.Writer) {
	tbl := table.New()
	tbl.WithWriter(w)
	for _, cgroupStats := range c {
		tbl.AddRow("Name:", cgroupStats.Name)
		tbl.AddRow("PIDs:", common.DisplayRatio(cgroupStats.PID.Current, cgroupStats.PID.Limit, common.WithTotal()))
		tbl.AddRow("CPU Usage:", cgroupStats.CPU.Usage)
		tbl.AddRow("Throttled Periods:", common.DisplayRatio(cgroupStats.CPU.NumThrottledPeriods, cgroupStats.CPU.NumRunnablePeriods, common.WithTotal()))
		tbl.AddRow("Throttled Time:", cgroupStats.CPU.ThrottledTimeInUsec)
		tbl.AddRow("System Usage", common.DisplayRatio(cgroupStats.CPU.SystemTimeInUsec, cgroupStats.CPU.UsageInUsec, common.WithTotal()))
		tbl.AddRow("User Usage", common.DisplayRatio(cgroupStats.CPU.UserTimeInUsec, cgroupStats.CPU.UsageInUsec, common.WithTotal()))
	}

	tbl.Print()
}
