package writer

import (
	"cgstat/stats"
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

	c.printMemStat("Usage", cgstats.CurrentUsage, cgstats.UsageLimit)
	c.printMemStat("MaxUsage", cgstats.MaxUsage, cgstats.UsageLimit)
	c.printMemStat("RSS", cgstats.Rss, cgstats.TotalRss)
	c.printMemStat("RSSHuge", cgstats.RssHuge, cgstats.RssHugeTotal)
	c.printMemStat("Writeback", cgstats.WriteBack, cgstats.TotalWriteBack)
	c.printMemStat("Cache", cgstats.CacheSize, cgstats.TotalCacheSize)
	c.printMemStat("Dirty", cgstats.DirtySize, cgstats.TotalDirtySize)
	c.printMemStat("KernelUsage", cgstats.KernelUsage, cgstats.KernelUsageLimit)
	c.printMemStat("KernelMax", cgstats.KernelMaxUsage,cgstats.KernelUsageLimit)
	c.printMemStat("KernelTCPUsage", cgstats.KernelTCPUsage, cgstats.KernelTCPLimit)
	c.printMemStat("KernelTCPMax", cgstats.KernelTCPMax, cgstats.KernelTCPLimit)
}

func (c *CgStatsVerboseWriter) printMemStat(name string, value uint64, maxValue uint64) {
	percentage := 0.0
	if maxValue != 0 {
		percentage = float64(value) / float64(maxValue) * 100
	}
	tab := "\t"
	if len(name) < 7  {
		tab = "\t\t"
	}
	fmt.Fprintf(c.writer,"\t%s:%s%v (%.2f%%)\n", name, tab, FormatBytes(value), percentage)
}

func (c *CgStatsVerboseWriter) printCPUStats(cgstats *stats.CgroupStats) {
	fmt.Fprintln(c.writer, "CPU Stats")

	throttling := 0.0
	if cgstats.TotalPeriods != 0 {
		throttling = float64(cgstats.ThrottlePeriods) / float64(cgstats.TotalPeriods) * 100.0
	}

	c.printCpuStat("UserCPU", cgstats.UserCPU)
	c.printCpuStat("KernelCPU", cgstats.KernelCPU)
	c.printCpuStat("ThrottlePeriods", throttling)
}

func (c *CgStatsVerboseWriter) printCpuStat(name string, value float64) {
	tab := "\t"
	if len(name) < 10  {
		tab = "\t\t"
	}
	fmt.Fprintf(c.writer,"\t%s:%s%.2f%%\n", name, tab, value)
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
