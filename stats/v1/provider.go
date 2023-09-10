package v1

import (
	cgroups "github.com/containerd/cgroups/v3/cgroup1"
	v1 "github.com/containerd/cgroups/v3/cgroup1/stats"
	"github.com/strategicpause/cgstat/stats/common"
	"time"
)

type CgroupStatsProvider struct {
	commonProvider               *common.CommonCgroupStatsProvider
	previousCPUStatsByCgroupPath map[string]*CgroupStats
}

const (
	CgroupPrefix = "/sys/fs/cgroup/pids"
)

func NewCgroupStatsProvider() *CgroupStatsProvider {
	return &CgroupStatsProvider{
		commonProvider:               common.NewCommonCgroupStatsProvider(CgroupPrefix),
		previousCPUStatsByCgroupPath: map[string]*CgroupStats{},
	}
}

func (c *CgroupStatsProvider) GetCgroupStatsByPrefix(prefix string) (common.CgroupStatsCollection, error) {
	paths := c.ListCgroupsByPrefix(prefix)
	return c.getCgroupStatsByPath(paths)
}

func (c *CgroupStatsProvider) ListCgroupsByPrefix(cgroupPrefix string) []string {
	return c.commonProvider.ListCgroupsByPrefix(cgroupPrefix)
}

func (c *CgroupStatsProvider) GetCgroupStatsByName(name string) (common.CgroupStatsCollection, error) {
	paths := []string{name}

	return c.getCgroupStatsByPath(paths)
}

func (c *CgroupStatsProvider) getCgroupStatsByPath(cgroupPaths []string) (common.CgroupStatsCollection, error) {
	var stats []*CgroupStats
	for _, cgroupPath := range cgroupPaths {
		control, err := cgroups.Load(cgroups.StaticPath(cgroupPath))
		if err != nil {
			return nil, err
		}
		cgroupStats, err := c.getCgroupStats(cgroupPath, control)
		if err != nil {
			return nil, err
		}
		stats = append(stats, cgroupStats)
	}
	return NewCollection(stats), nil
}

func (c *CgroupStatsProvider) getCgroupStats(name string, control cgroups.Cgroup) (*CgroupStats, error) {
	metrics, err := control.Stat(cgroups.IgnoreNotExist)

	if err != nil {
		return nil, err
	}
	cgStats := &CgroupStats{
		Name: name,
	}
	processes, err := control.Processes("cpu", true)
	if err != nil {
		return nil, err
	}

	prevStats := c.previousCPUStatsByCgroupPath[name]

	c.withProcessStats(cgStats, processes)
	c.withCpuStats(cgStats, metrics.CPU, prevStats)
	c.withMemoryOomControl(cgStats, metrics.MemoryOomControl)
	c.withMemoryStats(cgStats, metrics.Memory)
	c.withIOStats(cgStats, metrics.Blkio)

	c.previousCPUStatsByCgroupPath[name] = cgStats

	return cgStats, nil
}

func (c *CgroupStatsProvider) withProcessStats(cgStats *CgroupStats, processes []cgroups.Process) {
	cgStats.NumProcesses = uint64(len(processes))
}

func (c *CgroupStatsProvider) withCpuStats(cgStats *CgroupStats, cpuMetrics *v1.CPUStat, prevStats *CgroupStats) {
	cgStats.SystemTime = time.Now().UnixMicro()
	cgStats.CPUUsage = cpuMetrics.GetUsage().Total
	cgStats.ThrottlePeriods = cpuMetrics.Throttling.ThrottledPeriods
	cgStats.TotalPeriods = cpuMetrics.Throttling.Periods

	if prevStats == nil {
		cgStats.CPUUtilization = 0.0
	} else {
		cpuUsageDelta := float64(cgStats.CPUUsage - prevStats.CPUUsage)
		systemTimeDelta := float64(cgStats.SystemTime - prevStats.SystemTime)
		cgStats.CPUUtilization = (cpuUsageDelta / systemTimeDelta) * 100.0
	}
}

func (c *CgroupStatsProvider) withMemoryOomControl(cgStats *CgroupStats, oomMetrics *v1.MemoryOomControl) {
	if oomMetrics == nil {
		return
	}

	cgStats.UnderOom = oomMetrics.UnderOom
	cgStats.OomKill = oomMetrics.OomKill
}

func (c *CgroupStatsProvider) withMemoryStats(cgStats *CgroupStats, memMetrics *v1.MemoryStat) {
	if memMetrics == nil {
		return
	}
	cgStats.CurrentUsage = memMetrics.Usage.Usage
	cgStats.UsageLimit = memMetrics.Usage.Limit
	cgStats.CurrentUtilization = float64(memMetrics.Usage.Usage) / float64(memMetrics.Usage.Limit) * 100.0
	cgStats.MaxUsage = memMetrics.Usage.Max
	cgStats.MaxUtilization = float64(memMetrics.Usage.Max) / float64(memMetrics.Usage.Limit) * 100.0
	cgStats.Rss = memMetrics.RSS
	cgStats.PgPgIn = memMetrics.PgPgIn
	cgStats.PgPgOut = memMetrics.TotalPgPgOut
	cgStats.PgMajFault = memMetrics.PgMajFault
	cgStats.ActiveAnon = memMetrics.ActiveAnon
	cgStats.InactiveAnon = memMetrics.InactiveAnon
	cgStats.ActiveFile = memMetrics.ActiveFile
	cgStats.InactiveFile = memMetrics.TotalInactiveFile
	cgStats.Unevictable = memMetrics.Unevictable
	cgStats.CacheSize = memMetrics.Cache
	cgStats.DirtySize = memMetrics.Dirty
	cgStats.WriteBack = memMetrics.Writeback
}

func (c *CgroupStatsProvider) withIOStats(cgStats *CgroupStats, ioMetrics *v1.BlkIOStat) {
	if ioMetrics == nil {
		return
	}
	cgStats.IoServicedRecursive = c.getBlockDeviceStats(ioMetrics.IoServiceTimeRecursive)
	cgStats.IoServiceBytesRecursive = c.getBlockDeviceStats(ioMetrics.IoServiceBytesRecursive)
	cgStats.IoQueuedRecursive = c.getBlockDeviceStats(ioMetrics.IoQueuedRecursive)
	cgStats.IoTimeRecursive = c.getBlockDeviceStats(ioMetrics.IoTimeRecursive)
	cgStats.IoMergedRecursive = c.getBlockDeviceStats(ioMetrics.IoMergedRecursive)
	cgStats.IoWaitTimeRecursive = c.getBlockDeviceStats(ioMetrics.IoWaitTimeRecursive)
	cgStats.SectorsRecursive = c.getBlockDeviceStats(ioMetrics.SectorsRecursive)
	cgStats.IoServiceTimeRecursive = c.getBlockDeviceStats(ioMetrics.IoServiceTimeRecursive)
}

func (c *CgroupStatsProvider) getBlockDeviceStats(entries []*v1.BlkIOEntry) map[string]*BlockDevice {
	devices := make(map[string]*BlockDevice)
	for _, entry := range entries {
		deviceName := entry.Device
		device, ok := devices[deviceName]
		if !ok {
			device = &BlockDevice{}
			devices[deviceName] = device
		}
		switch entry.Op {
		case "Read":
			device.Read = entry.Value
		case "Write":
			device.Write = entry.Value
		case "Sync":
			device.Sync = entry.Value
		case "Total":
			device.Total = entry.Value
		case "Async":
			device.Async = entry.Value
		}
	}
	return devices
}
