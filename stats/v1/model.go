package v1

import (
	"fmt"
	"github.com/strategicpause/cgstat/stats/common"
	"io"
	"sort"
	"time"
)

type Cgroup struct {
	Name string
}

type BlockDevice struct {
	Read  uint64
	Write uint64
	Async uint64
	Sync  uint64
	Total uint64
}

type CgroupStats struct {
	Name string
	/** CPU **/
	CPU float64
	// The total CPU throttled time
	ThrottlePeriods uint64
	//
	TotalPeriods uint64
	// The number of processes currently in the cgroup and its descendants.
	NumProcesses uint64
	// Hard limit of number of processes.
	MaxProcesses uint64
	/** Memory **/
	CurrentUsage       uint64
	UsageLimit         uint64
	CurrentUtilization float64
	MaxUsage           uint64
	MaxUtilization     float64
	// Number of bytes of anonymous and swap cache memory (includes transparent hugepages).
	Rss uint64
	// Number of bytes of anonymous transparent hugepages.
	RssHuge          uint64
	KernelUsage      uint64
	KernelMaxUsage   uint64
	KernelUsageLimit uint64
	KernelTCPUsage   uint64
	KernelTCPMax     uint64
	KernelTCPLimit   uint64
	CacheSize        uint64
	// The total amount of memory waiting to be written back to the disk.
	DirtySize uint64
	// The total amount of memory actively being written back to the disk.
	WriteBack uint64
	// Number of bytes the system has paged in from disk per second.
	PgPgIn uint64
	// Number of bytes the system has paged out to disk per second.
	PgPgOut uint64
	// Number of page faults the system has made per second.
	PgFault uint64
	// Number of major faults per second the system required loading a memory page from disk.
	PgMajFault uint64
	// The amount of anonymous and tmpfs/shmem memory, that is in active use, or was in active use since
	// the last time the system moved something to swap.
	ActiveAnon uint64
	// The amount of anonymous and tmpfs/shmem memory, that is a candidate for eviction
	InactiveAnon uint64
	// The amount of file cache memory that is in active use, or was in active use since the last time the
	// system reclaimed memory.
	ActiveFile uint64
	// The amount of file cache memory that is newly loaded from the disk, or is a candidate for reclaiming.
	InactiveFile uint64
	// The amount of memory discovered by the pageout code, that is not evictable because it is locked into
	// memory by user programs.
	Unevictable uint64
	// The number of processes belonging to this cgroup killed by any kind of OOM killer.
	OomKill uint64
	// The cgroup is under OOM, tasks may be stopped.
	UnderOom uint64
	/** IO Stats **/
	// The total amount of time the IOs for this cgroup spent waiting in the scheduler queues for service.
	IoWaitTimeRecursive map[string]*BlockDevice
	// The disk time allocated to cgroup per device in milliseconds.
	IoTimeRecursive map[string]*BlockDevice
	// The total number of requests queued up at any given instant for this cgroup.
	IoQueuedRecursive map[string]*BlockDevice
	// The total number of bios/requests merged into requests belonging to this cgroup.
	IoMergedRecursive map[string]*BlockDevice
	// Number of bytes transferred to and from the block device.
	IoServiceBytesRecursive map[string]*BlockDevice
	// The total amount of time between request dispatch and request completion for the IOs done by this cgroup.
	IoServiceTimeRecursive map[string]*BlockDevice
	// The number of sectors transferred to/from disk by the group.
	SectorsRecursive map[string]*BlockDevice
	// The number of IOs (bio) issued to the disk by the group.
	IoServicedRecursive map[string]*BlockDevice
}

func (c *CgroupStats) ToCsvRow() []string {
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

func (c *CgroupStats) ToDisplayRow() []interface{} {
	CPU := fmt.Sprintf("%.2f%%", c.CPU)
	currentUsage := fmt.Sprintf("%s (%.2f%%)", common.FormatBytes(c.CurrentUsage), c.CurrentUtilization)
	maxUsage := fmt.Sprintf("%s (%.2f%%)", common.FormatBytes(c.MaxUsage), c.MaxUtilization)
	usageLimit := common.FormatBytes(c.UsageLimit)
	rss := common.FormatBytes(c.Rss)
	cacheSize := common.FormatBytes(c.CacheSize)
	dirtySize := common.FormatBytes(c.DirtySize)
	writeback := common.FormatBytes(c.WriteBack)

	return []interface{}{c.Name, CPU, c.NumProcesses, currentUsage, maxUsage, usageLimit, rss,
		cacheSize, dirtySize, writeback, c.UnderOom, c.OomKill}
}

func (c *CgroupStats) ToVerboseOutput(w io.Writer) {
	c.printMemStats(w)
	c.printCPUStats(w)
	c.printBlkIOStats(w)
}

func (c *CgroupStats) printMemStats(writer io.Writer) {
	fmt.Fprintln(writer, "Memory Stats")

	c.printMemUtilization(writer, "Usage", c.CurrentUsage, c.UsageLimit)
	c.printMemUtilization(writer, "MaxUsage", c.MaxUsage, c.UsageLimit)
	c.printMemStat(writer, "RSS", c.Rss)
	c.printMemStat(writer, "RSSHuge", c.RssHuge)
	c.printMemStat(writer, "Writeback", c.WriteBack)
	c.printMemStat(writer, "Cache", c.CacheSize)
	c.printMemStat(writer, "Dirty", c.DirtySize)
	c.printMemStat(writer, "PgPgIn", c.PgPgIn)
	c.printMemStat(writer, "PgPgOut", c.PgPgOut)
	c.printCounter(writer, "PgFault", c.PgFault)
	c.printCounter(writer, "PgMajFault", c.PgMajFault)
	c.printMemStat(writer, "ActiveAnon", c.ActiveAnon)
	c.printMemStat(writer, "InactiveAnon", c.InactiveAnon)
	c.printMemStat(writer, "ActiveFile", c.ActiveFile)
	c.printMemStat(writer, "InactiveFile", c.InactiveFile)
	c.printMemStat(writer, "Unevictable", c.Unevictable)
	c.printMemUtilization(writer, "KernelUsage", c.KernelUsage, c.KernelUsageLimit)
	c.printMemUtilization(writer, "KernelMax", c.KernelMaxUsage, c.KernelUsageLimit)
	c.printMemUtilization(writer, "KernelTCPUsage", c.KernelTCPUsage, c.KernelTCPLimit)
	c.printMemUtilization(writer, "KernelTCPMax", c.KernelTCPMax, c.KernelTCPLimit)
}

func (c *CgroupStats) printMemUtilization(w io.Writer, name string, value uint64, maxValue uint64) {
	percentage := 0.0
	if maxValue != 0 {
		percentage = float64(value) / float64(maxValue) * 100.0
	}

	fmt.Fprintf(w, "\t%s:%s%v / %v (%.2f%%)\n", name, c.getTabs(name), common.FormatBytes(value), common.FormatBytes(maxValue), percentage)
}

func (c *CgroupStats) getTabs(name string) string {
	if len(name) < 7 {
		return "\t\t"
	}
	return "\t"
}

func (c *CgroupStats) printMemStat(w io.Writer, name string, value uint64) {
	fmt.Fprintf(w, "\t%s:%s%v\n", name, c.getTabs(name), common.FormatBytes(value))
}

func (c *CgroupStats) printCounter(w io.Writer, name string, value uint64) {
	fmt.Fprintf(w, "\t%s:%s%d\n", name, c.getTabs(name), value)
}

func (c *CgroupStats) printCPUStats(w io.Writer) {
	fmt.Fprintln(w, "CPU Stats")

	c.printCpuStat(w, "CPU", c.CPU)
	c.printCounter(w, "NumProcesses", c.NumProcesses)
	c.printCounter(w, "ThrottlePeriods", c.ThrottlePeriods)
	c.printCounter(w, "TotalPeriods", c.TotalPeriods)
}

func (c *CgroupStats) printCpuStat(w io.Writer, name string, value float64) {
	fmt.Fprintf(w, "\t%s:%s%.2f%%\n", name, c.getTabs(name), value)
}

func (c *CgroupStats) printBlkIOStats(w io.Writer) {
	c.printBlkIOStat(w, "IoWaitTime", c.IoWaitTimeRecursive)
	c.printBlkIOStat(w, "IoTimeRecursive", c.IoTimeRecursive)
	c.printBlkIOStat(w, "IoQueuedRecursive", c.IoQueuedRecursive)
	c.printBlkIOStat(w, "IoMergedRecursive", c.IoMergedRecursive)
	c.printBlkIOStat(w, "IoServiceBytesRecursive", c.IoServiceBytesRecursive)
	c.printBlkIOStat(w, "IoServiceTimeRecursive", c.IoServiceTimeRecursive)
	c.printBlkIOStat(w, "SectorsRecursive", c.SectorsRecursive)
	c.printBlkIOStat(w, "IoServicedRecursive", c.IoServicedRecursive)
}

func (c *CgroupStats) printBlkIOStat(w io.Writer, name string, devices map[string]*BlockDevice) {
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
