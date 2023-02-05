package common

import (
	"fmt"
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
	TotalPeriods    uint64
	NumProcesses    uint64
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
	currentUsage := fmt.Sprintf("%s (%.2f%%)", FormatBytes(c.CurrentUsage), c.CurrentUtilization)
	maxUsage := fmt.Sprintf("%s (%.2f%%)", FormatBytes(c.MaxUsage), c.MaxUtilization)
	usageLimit := FormatBytes(c.UsageLimit)
	rss := FormatBytes(c.Rss)
	cacheSize := FormatBytes(c.CacheSize)
	dirtySize := FormatBytes(c.DirtySize)
	writeback := FormatBytes(c.WriteBack)

	return []interface{}{c.Name, CPU, c.NumProcesses, currentUsage, maxUsage, usageLimit, rss,
		cacheSize, dirtySize, writeback, c.UnderOom, c.OomKill}
}
