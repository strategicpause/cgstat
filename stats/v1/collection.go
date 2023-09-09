package v1

import (
	"fmt"
	"github.com/strategicpause/cgstat/stats/common"
	"io"
	"sort"
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
		"Time", "Name", "UserCPU", "CurrentUsage", "MaxUsage", "UsageLimit", "RSS",
		"Cache", "Dirty", "WriteBack", "UnderOom", "OomKill",
	}
}

func (_ CgroupStatsCollection) toCSVRow(c *CgroupStats) []string {
	t, _ := time.Now().UTC().MarshalText()
	return []string{
		string(t),
		c.Name,
		fmt.Sprintf("%f", c.CPU),
		fmt.Sprintf("%d", c.CurrentUsage),
		fmt.Sprintf("%d", c.MaxUsage),
		fmt.Sprintf("%d", c.UsageLimit),
		fmt.Sprintf("%d", c.Rss),
		fmt.Sprintf("%d", c.CacheSize),
		fmt.Sprintf("%d", c.DirtySize),
		fmt.Sprintf("%d", c.WriteBack),
		fmt.Sprintf("%d", c.UnderOom),
		fmt.Sprintf("%d", c.OomKill),
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
		"Name", "CPU", "NumProcesses", "CurrentUsage", "MaxUsage", "UsageLimit",
		"RSS", "Cache", "Dirty", "WriteBack", "UnderOom", "OomKill",
	}
}

func (_ CgroupStatsCollection) toDisplayRow(c *CgroupStats) []interface{} {
	CPU := fmt.Sprintf("%.2f%%", c.CPU)
	numProcess := fmt.Sprintf("%d", c.NumProcesses)
	currentUsage := fmt.Sprintf("%s (%.2f%%)", common.FormatBytes(c.CurrentUsage), c.CurrentUtilization)
	maxUsage := fmt.Sprintf("%s (%.2f%%)", common.FormatBytes(c.MaxUsage), c.MaxUtilization)
	usageLimit := common.FormatBytes(c.UsageLimit)
	rss := common.FormatBytes(c.Rss)
	cacheSize := common.FormatBytes(c.CacheSize)
	dirtySize := common.FormatBytes(c.DirtySize)
	writeback := common.FormatBytes(c.WriteBack)
	underOom := fmt.Sprintf("%d", c.UnderOom)
	oomKill := fmt.Sprintf("%d", c.OomKill)

	return []interface{}{c.Name, CPU, numProcess, currentUsage, maxUsage, usageLimit, rss,
		cacheSize, dirtySize, writeback, underOom, oomKill}
}

func (c CgroupStatsCollection) ToVerboseOutput(w io.Writer) {
	for _, cgropStats := range c {
		c.printMemStats(w, cgropStats)
		c.printCPUStats(w, cgropStats)
		c.printBlkIOStats(w, cgropStats)
	}
}

func (c CgroupStatsCollection) printMemStats(writer io.Writer, s *CgroupStats) {
	fmt.Fprintln(writer, "Memory Stats")

	c.printMemUtilization(writer, "Usage", s.CurrentUsage, s.UsageLimit)
	c.printMemUtilization(writer, "MaxUsage", s.MaxUsage, s.UsageLimit)
	c.printMemStat(writer, "RSS", s.Rss)
	c.printMemStat(writer, "RSSHuge", s.RssHuge)
	c.printMemStat(writer, "Writeback", s.WriteBack)
	c.printMemStat(writer, "Cache", s.CacheSize)
	c.printMemStat(writer, "Dirty", s.DirtySize)
	c.printMemStat(writer, "PgPgIn", s.PgPgIn)
	c.printMemStat(writer, "PgPgOut", s.PgPgOut)
	c.printCounter(writer, "PgFault", s.PgFault)
	c.printCounter(writer, "PgMajFault", s.PgMajFault)
	c.printMemStat(writer, "ActiveAnon", s.ActiveAnon)
	c.printMemStat(writer, "InactiveAnon", s.InactiveAnon)
	c.printMemStat(writer, "ActiveFile", s.ActiveFile)
	c.printMemStat(writer, "InactiveFile", s.InactiveFile)
	c.printMemStat(writer, "Unevictable", s.Unevictable)
	c.printMemUtilization(writer, "KernelUsage", s.KernelUsage, s.KernelUsageLimit)
	c.printMemUtilization(writer, "KernelMax", s.KernelMaxUsage, s.KernelUsageLimit)
	c.printMemUtilization(writer, "KernelTCPUsage", s.KernelTCPUsage, s.KernelTCPLimit)
	c.printMemUtilization(writer, "KernelTCPMax", s.KernelTCPMax, s.KernelTCPLimit)
}

func (c CgroupStatsCollection) printMemUtilization(w io.Writer, name string, value uint64, maxValue uint64) {
	percentage := 0.0
	if maxValue != 0 {
		percentage = float64(value) / float64(maxValue) * 100.0
	}

	fmt.Fprintf(w, "\t%s:%s%v / %v (%.2f%%)\n", name, c.getTabs(name), common.FormatBytes(value), common.FormatBytes(maxValue), percentage)
}

func (_ CgroupStatsCollection) getTabs(name string) string {
	if len(name) < 7 {
		return "\t\t"
	}
	return "\t"
}

func (c CgroupStatsCollection) printMemStat(w io.Writer, name string, value uint64) {
	fmt.Fprintf(w, "\t%s:%s%v\n", name, c.getTabs(name), common.FormatBytes(value))
}

func (c CgroupStatsCollection) printCounter(w io.Writer, name string, value uint64) {
	fmt.Fprintf(w, "\t%s:%s%d\n", name, c.getTabs(name), value)
}

func (c CgroupStatsCollection) printCPUStats(w io.Writer, s *CgroupStats) {
	fmt.Fprintln(w, "CPU Stats")

	c.printCpuStat(w, "CPU", s.CPU)
	c.printCounter(w, "NumProcesses", s.NumProcesses)
	c.printCounter(w, "ThrottlePeriods", s.ThrottlePeriods)
	c.printCounter(w, "TotalPeriods", s.TotalPeriods)
}

func (c CgroupStatsCollection) printCpuStat(w io.Writer, name string, value float64) {
	fmt.Fprintf(w, "\t%s:%s%.2f%%\n", name, c.getTabs(name), value)
}

func (c CgroupStatsCollection) printBlkIOStats(w io.Writer, s *CgroupStats) {
	c.printBlkIOStat(w, "IoWaitTime", s.IoWaitTimeRecursive)
	c.printBlkIOStat(w, "IoTimeRecursive", s.IoTimeRecursive)
	c.printBlkIOStat(w, "IoQueuedRecursive", s.IoQueuedRecursive)
	c.printBlkIOStat(w, "IoMergedRecursive", s.IoMergedRecursive)
	c.printBlkIOStat(w, "IoServiceBytesRecursive", s.IoServiceBytesRecursive)
	c.printBlkIOStat(w, "IoServiceTimeRecursive", s.IoServiceTimeRecursive)
	c.printBlkIOStat(w, "SectorsRecursive", s.SectorsRecursive)
	c.printBlkIOStat(w, "IoServicedRecursive", s.IoServicedRecursive)
}

func (c CgroupStatsCollection) printBlkIOStat(w io.Writer, name string, devices map[string]*BlockDevice) {
	if len(devices) == 0 {
		return
	}

	fmt.Fprintf(w, "%s:\n", name)

	// Print values
	keys := make([]string, 0, len(devices))
	for deviceName := range devices {
		keys = append(keys, deviceName)
	}
	sort.Strings(keys)
	for _, deviceName := range keys {
		device := devices[deviceName]
		fmt.Fprintf(w, "\t%s:\t%v (Read) %v (Write) %v (Sync) %v (Async) %v (Total)\n",
			deviceName, device.Read, device.Write, device.Sync, device.Async, device.Total)
	}
}
