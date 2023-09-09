package v1

import (
	"fmt"
	"github.com/strategicpause/cgstat/stats/common"
	"io"
	"sort"
	"time"
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
		"Time", "Name", "UserCPU", "CurrentUsage", "MaxUsage", "UsageLimit", "RSS",
		"Cache", "Dirty", "WriteBack", "UnderOom", "OomKill",
	}
}

func toCSVRow(c *CgroupStats) []string {
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

func getDisplayHeaders() []interface{} {
	return []interface{}{
		"Name", "CPU", "NumProcesses", "CurrentUsage", "MaxUsage", "UsageLimit",
		"RSS", "Cache", "Dirty", "WriteBack", "UnderOom", "OomKill",
	}
}

func toDisplayRow(c *CgroupStats) []interface{} {
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

func toVerboseOutput(w io.Writer, c []*CgroupStats) {
	for _, cgropStats := range c {
		printMemStats(w, cgropStats)
		printCPUStats(w, cgropStats)
		printBlkIOStats(w, cgropStats)
	}
}

func printMemStats(writer io.Writer, s *CgroupStats) {
	fmt.Fprintln(writer, "Memory Stats")

	printMemUtilization(writer, "Usage", s.CurrentUsage, s.UsageLimit)
	printMemUtilization(writer, "MaxUsage", s.MaxUsage, s.UsageLimit)
	printMemStat(writer, "RSS", s.Rss)
	printMemStat(writer, "RSSHuge", s.RssHuge)
	printMemStat(writer, "Writeback", s.WriteBack)
	printMemStat(writer, "Cache", s.CacheSize)
	printMemStat(writer, "Dirty", s.DirtySize)
	printMemStat(writer, "PgPgIn", s.PgPgIn)
	printMemStat(writer, "PgPgOut", s.PgPgOut)
	printCounter(writer, "PgFault", s.PgFault)
	printCounter(writer, "PgMajFault", s.PgMajFault)
	printMemStat(writer, "ActiveAnon", s.ActiveAnon)
	printMemStat(writer, "InactiveAnon", s.InactiveAnon)
	printMemStat(writer, "ActiveFile", s.ActiveFile)
	printMemStat(writer, "InactiveFile", s.InactiveFile)
	printMemStat(writer, "Unevictable", s.Unevictable)
	printMemUtilization(writer, "KernelUsage", s.KernelUsage, s.KernelUsageLimit)
	printMemUtilization(writer, "KernelMax", s.KernelMaxUsage, s.KernelUsageLimit)
	printMemUtilization(writer, "KernelTCPUsage", s.KernelTCPUsage, s.KernelTCPLimit)
	printMemUtilization(writer, "KernelTCPMax", s.KernelTCPMax, s.KernelTCPLimit)
}

func printMemUtilization(w io.Writer, name string, value uint64, maxValue uint64) {
	percentage := 0.0
	if maxValue != 0 {
		percentage = float64(value) / float64(maxValue) * 100.0
	}

	fmt.Fprintf(w, "\t%s:%s%v / %v (%.2f%%)\n", name, getTabs(name), common.FormatBytes(value), common.FormatBytes(maxValue), percentage)
}

func getTabs(name string) string {
	if len(name) < 7 {
		return "\t\t"
	}
	return "\t"
}

func printMemStat(w io.Writer, name string, value uint64) {
	fmt.Fprintf(w, "\t%s:%s%v\n", name, getTabs(name), common.FormatBytes(value))
}

func printCounter(w io.Writer, name string, value uint64) {
	fmt.Fprintf(w, "\t%s:%s%d\n", name, getTabs(name), value)
}

func printCPUStats(w io.Writer, s *CgroupStats) {
	fmt.Fprintln(w, "CPU Stats")

	printCpuStat(w, "CPU", s.CPU)
	printCounter(w, "NumProcesses", s.NumProcesses)
	printCounter(w, "ThrottlePeriods", s.ThrottlePeriods)
	printCounter(w, "TotalPeriods", s.TotalPeriods)
}

func printCpuStat(w io.Writer, name string, value float64) {
	fmt.Fprintf(w, "\t%s:%s%.2f%%\n", name, getTabs(name), value)
}

func printBlkIOStats(w io.Writer, s *CgroupStats) {
	printBlkIOStat(w, "IoWaitTime", s.IoWaitTimeRecursive)
	printBlkIOStat(w, "IoTimeRecursive", s.IoTimeRecursive)
	printBlkIOStat(w, "IoQueuedRecursive", s.IoQueuedRecursive)
	printBlkIOStat(w, "IoMergedRecursive", s.IoMergedRecursive)
	printBlkIOStat(w, "IoServiceBytesRecursive", s.IoServiceBytesRecursive)
	printBlkIOStat(w, "IoServiceTimeRecursive", s.IoServiceTimeRecursive)
	printBlkIOStat(w, "SectorsRecursive", s.SectorsRecursive)
	printBlkIOStat(w, "IoServicedRecursive", s.IoServicedRecursive)
}

func printBlkIOStat(w io.Writer, name string, devices map[string]*BlockDevice) {
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
