package stats

import (
	"github.com/containerd/cgroups"
	v1 "github.com/containerd/cgroups/stats/v1"
	"io/ioutil"
	"path/filepath"
)

type CgroupStatsProvider struct {

}

const (
	CgroupPrefix = "/sys/fs/cgroup/pids"
)

func NewCgroupStatsProvider() *CgroupStatsProvider {
	return &CgroupStatsProvider{}
}

func (c *CgroupStatsProvider) GetCgroupStatsByPrefix(prefix string) ([]*CgroupStats, error) {
	paths, err := c.getPathsByPrefix(prefix)
	if err != nil {
		return nil, err
	}
	return c.getCgroupStatsByPath(paths)
}

func (c *CgroupStatsProvider) getPathsByPrefix(prefix string) ([]string, error) {
	var cgroupPaths []string

	prefixPath := filepath.Join(CgroupPrefix, prefix)

	files, err := ioutil.ReadDir(prefixPath)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if file.IsDir() {
			cgroupPath := filepath.Join(prefix, file.Name())
			cgroupPaths = append(cgroupPaths, cgroupPath)
		}
	}

	return cgroupPaths, nil
}

func (c *CgroupStatsProvider) GetCgroupStatsByName(name string) ([]*CgroupStats, error) {
	paths := []string {name}

	return c.getCgroupStatsByPath(paths)
}

func (c *CgroupStatsProvider) getCgroupStatsByPath(cgroupPaths []string) ([]*CgroupStats, error) {
	var stats []*CgroupStats
	for _, cgroup := range cgroupPaths {
		control, err := cgroups.Load(cgroups.V1, cgroups.StaticPath(cgroup))
		if err != nil {
			return nil, err
		}
		stat, err := control.Stat(cgroups.IgnoreNotExist)
		if err != nil {
			return nil, err
		}
		stats = append(stats, getCgroupStats(cgroup, stat))
	}
	return stats, nil
}

func getCgroupStats(name string, metrics *v1.Metrics) *CgroupStats {
	cgStats := &CgroupStats {
		Name:               name,
		UnderOom:           metrics.MemoryOomControl.UnderOom,
		OomKill:            metrics.MemoryOomControl.OomKill,
	}
	withCpuStats(cgStats, metrics.CPU)
	withMemoryStats(cgStats, metrics.Memory)
	withIOStats(cgStats, metrics.Blkio)

	return cgStats
}

func withCpuStats(cgStats *CgroupStats,  cpuMetrics *v1.CPUStat) {
	userCPU := 0.0
	kernelCPU := 0.0

	if cpuMetrics.Usage.Total != 0.0 {
		userCPU = float64(cpuMetrics.Usage.User) / float64(cpuMetrics.Usage.Total) * 100
		kernelCPU = float64(cpuMetrics.Usage.Kernel) / float64(cpuMetrics.Usage.Total) * 100
	}

	cgStats.UserCPU = userCPU
	cgStats.KernelCPU = kernelCPU
	cgStats.ThrottlePeriods = cpuMetrics.Throttling.ThrottledPeriods
	cgStats.TotalPeriods = cpuMetrics.Throttling.Periods
}

func withMemoryStats(cgStats *CgroupStats,  memMetrics *v1.MemoryStat) {
	cgStats.CurrentUsage = memMetrics.Usage.Usage
	cgStats.UsageLimit = memMetrics.Usage.Limit
	cgStats.CurrentUtilization = float64(memMetrics.Usage.Usage) / float64(memMetrics.Usage.Limit) * 100.0
	cgStats.MaxUsage = memMetrics.Usage.Max
	cgStats.MaxUtilization = float64(memMetrics.Usage.Max) / float64(memMetrics.Usage.Limit) * 100.0
	cgStats.Rss = memMetrics.RSS
	cgStats.CacheSize = memMetrics.Cache
	cgStats.DirtySize = memMetrics.Dirty
	cgStats.WriteBack = memMetrics.Writeback
}

func withIOStats(cgStats *CgroupStats, ioMetrics *v1.BlkIOStat) {
	cgStats.IoServicedRecursive = getBlockDeviceStats(ioMetrics.IoServiceTimeRecursive)
	cgStats.IoServiceBytesRecursive = getBlockDeviceStats(ioMetrics.IoServiceBytesRecursive)
	cgStats.IoQueuedRecursive = getBlockDeviceStats(ioMetrics.IoQueuedRecursive)
	cgStats.IoTimeRecursive = getBlockDeviceStats(ioMetrics.IoTimeRecursive)
	cgStats.IoMergedRecursive = getBlockDeviceStats(ioMetrics.IoMergedRecursive)
	cgStats.IoWaitTimeRecursive = getBlockDeviceStats(ioMetrics.IoWaitTimeRecursive)
	cgStats.SectorsRecursive = getBlockDeviceStats(ioMetrics.SectorsRecursive)
	cgStats.IoServiceTimeRecursive = getBlockDeviceStats(ioMetrics.IoServiceTimeRecursive)
}

func getBlockDeviceStats(entries []*v1.BlkIOEntry) map[string]*BlockDevice {
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