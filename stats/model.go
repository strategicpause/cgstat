package stats

import (
	"time"
)

type CgstatArgs struct {
	CgroupName string
	CgroupPrefix string
	VerboseOutput bool
	OutputFile string
	FollowMode bool
	RefreshInterval float64
}

func (c *CgstatArgs) HasPrefix() bool {
	return c.CgroupPrefix != ""
}

func (c *CgstatArgs) HasOutputFile() bool {
	return c.OutputFile != ""
}

func (c *CgstatArgs) GetRefreshInterval() time.Duration {
	return time.Duration(c.RefreshInterval * float64(time.Second))
}

type BlockDevice struct {
	Read uint64
	Write uint64
	Async uint64
	Sync uint64
	Total uint64
}

type CgroupStats struct {
	Name string
	// CPU
	CPU float64
	ThrottlePeriods uint64
	TotalPeriods uint64
	// Memory
	CurrentUsage uint64
	UsageLimit uint64
	CurrentUtilization float64
	MaxUsage uint64
	MaxUtilization float64
	Rss uint64
	TotalRss uint64
	RssHuge uint64
	RssHugeTotal uint64
	KernelUsage uint64
	KernelMaxUsage uint64
	KernelUsageLimit uint64
	KernelTCPUsage uint64
	KernelTCPMax uint64
	KernelTCPLimit uint64
	CacheSize uint64
	TotalCacheSize uint64
	DirtySize uint64
	TotalDirtySize uint64
	WriteBack uint64
	TotalWriteBack uint64
	// The number of processes belonging to this cgroup killed by any kind of OOM killer.
	OomKill uint64
	// The cgroup is under OOM, tasks may be stopped.
	UnderOom uint64
	// IO Stats
	IoWaitTimeRecursive map[string]*BlockDevice
	IoTimeRecursive map[string]*BlockDevice
	IoQueuedRecursive map[string]*BlockDevice
	IoMergedRecursive map[string]*BlockDevice
	IoServiceBytesRecursive map[string]*BlockDevice
	IoServiceTimeRecursive map[string]*BlockDevice
	SectorsRecursive map[string]*BlockDevice
	IoServicedRecursive map[string]*BlockDevice
}
