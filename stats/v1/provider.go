package v1

import (
	"github.com/containerd/cgroups"
	v1 "github.com/containerd/cgroups/stats/v1"
	"github.com/strategicpause/cgstat/stats/common"
	"github.com/struCoder/pidusage"
	"log"
	"os"
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

func (c *CgroupStatsProvider) GetCgroupStatsByPrefix(prefix string) ([]*common.CgroupStats, error) {
	paths := c.ListCgroupsByPrefix(prefix)
	return c.getCgroupStatsByPath(paths)
}

func (c *CgroupStatsProvider) ListCgroupsByPrefix(cgroupPrefix string) []string {
	var cgroupPaths []string
	queue := []string{cgroupPrefix}

	for len(queue) > 0 {
		prefix := queue[0]
		queue = queue[1:]

		prefixPath := filepath.Join(CgroupPrefix, prefix)
		files, err := os.ReadDir(prefixPath)
		if err != nil {
			log.Println(err)
		}

		for _, file := range files {
			if file.IsDir() {
				cgroupPath := filepath.Join(prefix, file.Name())
				cgroupPaths = append(cgroupPaths, cgroupPath)
				if prefix != cgroupPath {
					queue = append(queue, cgroupPath)
				}
			}
		}
	}
	return cgroupPaths
}

func (c *CgroupStatsProvider) GetCgroupStatsByName(name string) ([]*common.CgroupStats, error) {
	paths := []string{name}

	return c.getCgroupStatsByPath(paths)
}

func (c *CgroupStatsProvider) getCgroupStatsByPath(cgroupPaths []string) ([]*common.CgroupStats, error) {
	var stats []*common.CgroupStats
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

func (c *CgroupStatsProvider) getCgroupStats(name string, control cgroups.Cgroup) (*common.CgroupStats, error) {
	metrics, err := control.Stat(cgroups.IgnoreNotExist)

	if err != nil {
		return nil, err
	}
	cgStats := &common.CgroupStats{
		Name: name,
	}

	err = c.withCpuStats(cgStats, control, metrics.CPU)
	if err != nil {
		return nil, err
	}
	c.withMemoryOomControl(cgStats, metrics.MemoryOomControl)
	c.withMemoryStats(cgStats, metrics.Memory)
	c.withIOStats(cgStats, metrics.Blkio)

	return cgStats, nil
}

func (c *CgroupStatsProvider) withCpuStats(cgStats *common.CgroupStats, control cgroups.Cgroup, cpuMetrics *v1.CPUStat) error {
	processes, err := control.Processes("cpu", true)
	if err != nil {
		return err
	}
	for _, proc := range processes {
		stat, pidErr := pidusage.GetStat(proc.Pid)
		if pidErr != nil {
			return pidErr
		}
		cgStats.CPU += stat.CPU
	}

	cgStats.NumProcesses = uint64(len(processes))
	cgStats.ThrottlePeriods = cpuMetrics.Throttling.ThrottledPeriods
	cgStats.TotalPeriods = cpuMetrics.Throttling.Periods

	return nil
}

func (c *CgroupStatsProvider) withMemoryOomControl(cgStats *common.CgroupStats, oomMetrics *v1.MemoryOomControl) {
	if oomMetrics == nil {
		return
	}

	cgStats.UnderOom = oomMetrics.UnderOom
	cgStats.OomKill = oomMetrics.OomKill
}

func (c *CgroupStatsProvider) withMemoryStats(cgStats *common.CgroupStats, memMetrics *v1.MemoryStat) {
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

func (c *CgroupStatsProvider) withIOStats(cgStats *common.CgroupStats, ioMetrics *v1.BlkIOStat) {
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

func (c *CgroupStatsProvider) getBlockDeviceStats(entries []*v1.BlkIOEntry) map[string]*common.BlockDevice {
	devices := make(map[string]*common.BlockDevice)
	for _, entry := range entries {
		deviceName := entry.Device
		device, ok := devices[deviceName]
		if !ok {
			device = &common.BlockDevice{}
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
