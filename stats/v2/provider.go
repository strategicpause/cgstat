package v2

import (
	"github.com/containerd/cgroups/v3/cgroup2"
	"github.com/containerd/cgroups/v3/cgroup2/stats"
	"github.com/strategicpause/cgstat/stats/common"
)

const (
	CgroupPrefix = "/sys/fs/cgroup"
)

type CgroupStatsProvider struct {
	commonProvider *common.CommonCgroupStatsProvider
}

func NewCgroupStatsProvider() common.CgroupStatsProvider {
	return &CgroupStatsProvider{
		commonProvider: common.NewCommonCgroupStatsProvider(CgroupPrefix),
	}
}

func (c *CgroupStatsProvider) ListCgroupsByPrefix(cgroupPrefix string) []string {
	return c.commonProvider.ListCgroupsByPrefix(cgroupPrefix)
}

func (c *CgroupStatsProvider) GetCgroupStatsByPrefix(prefix string) (common.CgroupStatsCollection, error) {
	paths := c.ListCgroupsByPrefix(prefix)
	return c.getCgroupStatsByPath(paths)
}

func (c *CgroupStatsProvider) GetCgroupStatsByName(name string) (common.CgroupStatsCollection, error) {
	paths := []string{name}

	return c.getCgroupStatsByPath(paths)
}

func (c *CgroupStatsProvider) getCgroupStatsByPath(cgroupPaths []string) (common.CgroupStatsCollection, error) {
	var statsCollection CgroupStatsCollection

	for _, cgroupPath := range cgroupPaths {
		mgr, err := cgroup2.Load(cgroupPath)
		if err != nil {
			return nil, err
		}
		metrics, err := mgr.Stat()
		if err != nil {
			return nil, err
		}
		cgroupStats := NewCgroupStat(cgroupPath,
			c.withCPU(metrics.GetCPU()),
			c.withPids(metrics.GetPids()),
			c.withMemory(metrics.GetMemory()),
			c.withMemoryEvents(metrics.GetMemoryEvents()),
			c.WithIO(metrics.GetIo()),
		)
		statsCollection = append(statsCollection, cgroupStats)
	}

	return statsCollection, nil
}

func (c *CgroupStatsProvider) withCPU(cpu *stats.CPUStat) CgroupStatsOpt {
	return func(cgroupStats *CgroupStats) {
		if cpu == nil {
			return
		}
		// TODO
	}
}

func (c *CgroupStatsProvider) withPids(pids *stats.PidsStat) CgroupStatsOpt {
	return func(cgroupStats *CgroupStats) {
		if pids == nil {
			return
		}
		cgroupStats.NumPids = pids.GetCurrent()
		cgroupStats.MaxPids = pids.GetLimit()
	}
}

func (c *CgroupStatsProvider) withMemory(memory *stats.MemoryStat) CgroupStatsOpt {
	return func(cgroupStats *CgroupStats) {
		if memory == nil {
			return
		}
		cgroupStats.CurrentUsage = memory.GetUsage()
		cgroupStats.UsageLimit = memory.GetUsageLimit()
		// TODO
	}
}

func (c *CgroupStatsProvider) withMemoryEvents(memoryEvents *stats.MemoryEvents) CgroupStatsOpt {
	return func(cgroupStats *CgroupStats) {
		if memoryEvents == nil {
			return
		}
		cgroupStats.NumOomEvents = memoryEvents.GetOom()
		cgroupStats.NumOomKillEvents = memoryEvents.GetOomKill()
		cgroupStats.MemoryHigh = memoryEvents.GetHigh()
		cgroupStats.MemoryMax = memoryEvents.GetMax()
		cgroupStats.MemoryLow = memoryEvents.GetLow()
	}
}

func (c *CgroupStatsProvider) WithIO(io *stats.IOStat) CgroupStatsOpt {
	return func(cgroupStats *CgroupStats) {
		if io == nil {
			return
		}
		// TODO
	}
}
