package writer

import (
	"github.com/strategicpause/cgstat/stats"

	"fmt"
	"github.com/gosuri/uilive"
	"sort"
)

type CgStatsVerboseWriter struct {
	writer *uilive.Writer
}

func NewCgroupStatsVerboseWriter() StatsWriter {
	writer := uilive.New()
	writer.Start()

	return &CgStatsVerboseWriter{
		writer: writer,
	}
}

func (c *CgStatsVerboseWriter) Write(cgStats []*stats.CgroupStats) error {
	for _, s := range cgStats {
		c.printMemStats(s)
		c.printCPUStats(s)
		c.printBlkIOStats(s)
	}
	c.writer.Flush()
	return nil
}

func (c *CgStatsVerboseWriter) printMemStats(cgstats *stats.CgroupStats) {
	fmt.Fprintln(c.writer, "Memory Stats")

	c.printMemUtilization("Usage", cgstats.CurrentUsage, cgstats.UsageLimit)
	c.printMemUtilization("MaxUsage", cgstats.MaxUsage, cgstats.UsageLimit)
	c.printMemStat("RSS", cgstats.Rss)
	c.printMemStat("RSSHuge", cgstats.RssHuge)
	c.printMemStat("Writeback", cgstats.WriteBack)
	c.printMemStat("Cache", cgstats.CacheSize)
	c.printMemStat("Dirty", cgstats.DirtySize)
	c.printMemStat("PgPgIn", cgstats.PgPgIn)
	c.printMemStat("PgPgOut", cgstats.PgPgOut)
	c.printCounter("PgFault", cgstats.PgFault)
	c.printCounter("PgMajFault", cgstats.PgMajFault)
	c.printMemStat("ActiveAnon", cgstats.ActiveAnon)
	c.printMemStat("InactiveAnon", cgstats.InactiveAnon)
	c.printMemStat("ActiveFile", cgstats.ActiveFile)
	c.printMemStat("InactiveFile", cgstats.InactiveFile)
	c.printMemStat("Unevictable", cgstats.Unevictable)
	c.printMemUtilization("KernelUsage", cgstats.KernelUsage, cgstats.KernelUsageLimit)
	c.printMemUtilization("KernelMax", cgstats.KernelMaxUsage,cgstats.KernelUsageLimit)
	c.printMemUtilization("KernelTCPUsage", cgstats.KernelTCPUsage, cgstats.KernelTCPLimit)
	c.printMemUtilization("KernelTCPMax", cgstats.KernelTCPMax, cgstats.KernelTCPLimit)
}

func (c *CgStatsVerboseWriter) printMemUtilization(name string, value uint64, maxValue uint64) {
	percentage := 0.0
	if maxValue != 0 {
		percentage = float64(value) / float64(maxValue) * 100.0
	}

	fmt.Fprintf(c.writer,"\t%s:%s%v / %v (%.2f%%)\n", name, c.getTabs(name), FormatBytes(value), FormatBytes(maxValue), percentage)
}

func (c *CgStatsVerboseWriter) getTabs(name string) string {
	if len(name) < 7  {
		return "\t\t"
	}
	return "\t"
}

func (c *CgStatsVerboseWriter) printMemStat(name string, value uint64) {
	fmt.Fprintf(c.writer,"\t%s:%s%v\n", name, c.getTabs(name), FormatBytes(value))
}

func (c *CgStatsVerboseWriter) printCounter(name string, value uint64) {
	fmt.Fprintf(c.writer,"\t%s:%s%d\n", name, c.getTabs(name), value)
}


func (c *CgStatsVerboseWriter) printCPUStats(cgstats *stats.CgroupStats) {
	fmt.Fprintln(c.writer, "CPU Stats")

	c.printCpuStat("CPU", cgstats.CPU)
	c.printCounter("NumProcesses", cgstats.NumProcesses)
	c.printCounter("ThrottlePeriods", cgstats.ThrottlePeriods)
	c.printCounter("TotalPeriods", cgstats.TotalPeriods)
}

func (c *CgStatsVerboseWriter) printCpuStat(name string, value float64) {
	fmt.Fprintf(c.writer,"\t%s:%s%.2f%%\n", name, c.getTabs(name), value)
}

func (c *CgStatsVerboseWriter) printBlkIOStats(cgstats *stats.CgroupStats) {
	c.printBlkIOStat("IoWaitTime", cgstats.IoWaitTimeRecursive)
	c.printBlkIOStat("IoTimeRecursive", cgstats.IoTimeRecursive)
	c.printBlkIOStat("IoQueuedRecursive", cgstats.IoQueuedRecursive)
	c.printBlkIOStat( "IoMergedRecursive", cgstats.IoMergedRecursive)
	c.printBlkIOStat( "IoServiceBytesRecursive", cgstats.IoServiceBytesRecursive)
	c.printBlkIOStat( "IoServiceTimeRecursive", cgstats.IoServiceTimeRecursive)
	c.printBlkIOStat( "SectorsRecursive", cgstats.SectorsRecursive)
	c.printBlkIOStat( "IoServicedRecursive", cgstats.IoServicedRecursive)
}

func (c *CgStatsVerboseWriter) printBlkIOStat(name string, devices map[string]*stats.BlockDevice) {
	if len(devices) == 0 {
		return
	}

	fmt.Fprintf(c.writer,"%s:\n", name)

	// Print values
	keys := make([]string, 0, len(devices))
	for deviceName := range devices {
		keys = append(keys, deviceName)
	}
	sort.Strings(keys)
	for _, deviceName := range keys {
		device := devices[deviceName]
		fmt.Fprintf(c.writer,"\t%s:\t%v (Read) %v (Write) %v (Sync) %v (Async) %v (Total)\n",
			deviceName, device.Read, device.Write, device.Sync, device.Async, device.Total)
	}
}
