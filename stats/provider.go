package stats

import (
	"github.com/containerd/cgroups"
	v1 "github.com/containerd/cgroups/stats/v1"
	"github.com/struCoder/pidusage"
	"io/ioutil"
	"path/filepath"
	"time"
)

const (
	// CPURefreshInterval determines how often to refresh CPU metrics
	CPURefreshInterval = 100 * time.Millisecond
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
		cgroupStats, err := c.getCgroupStats(cgroup, control)
		if err != nil {
			return nil, err
		}
		stats = append(stats, cgroupStats)
	}
	return stats, nil
}

func (c *CgroupStatsProvider) getCgroupStats(name string, control cgroups.Cgroup) (*CgroupStats, error) {
	metrics, err := control.Stat(cgroups.IgnoreNotExist)

	if err != nil {
		return nil, err
	}
	cgStats := &CgroupStats {
		Name:               name,
		UnderOom:           metrics.MemoryOomControl.UnderOom,
		OomKill:            metrics.MemoryOomControl.OomKill,
	}

	err = c.withCpuStats(cgStats, control, metrics.CPU)
	if err != nil {
		return nil, err
	}
	c.withMemoryStats(cgStats, metrics.Memory)
	c.withIOStats(cgStats, metrics.Blkio)

	return cgStats, nil
}

func (c *CgroupStatsProvider) withCpuStats(cgStats *CgroupStats, control cgroups.Cgroup, cpuMetrics *v1.CPUStat) error {
	processes, err := control.Processes("cpu", true)
	if err != nil {
		return err
	}
	for _, proc := range processes {
		stat, _ := pidusage.GetStat(proc.Pid)
		cgStats.CPU += stat.CPU
	}

	cgStats.ThrottlePeriods = cpuMetrics.Throttling.ThrottledPeriods
	cgStats.TotalPeriods = cpuMetrics.Throttling.Periods

	return nil
}

func (c *CgroupStatsProvider) withMemoryStats(cgStats *CgroupStats,  memMetrics *v1.MemoryStat) {
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

func (c *CgroupStatsProvider) withIOStats(cgStats *CgroupStats, ioMetrics *v1.BlkIOStat) {
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