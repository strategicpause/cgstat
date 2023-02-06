package v2

import (
	"fmt"
	"github.com/containerd/cgroups/v3/cgroup2"
	"github.com/containerd/cgroups/v3/cgroup2/stats"
	"github.com/strategicpause/cgstat/stats/common"
	"github.com/struCoder/pidusage"
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
			c.withCPU(mgr, metrics.GetCPU()),
			c.withPids(metrics.GetPids()),
			c.withMemory(metrics.GetMemory()),
			c.withMemoryEvents(metrics.GetMemoryEvents()),
		)
		statsCollection = append(statsCollection, cgroupStats)
	}

	return statsCollection, nil
}

func (c *CgroupStatsProvider) withCPU(mgr *cgroup2.Manager, cpu *stats.CPUStat) CgroupStatsOpt {
	return func(cgroupStats *CgroupStats) {
		cgroupStats.CPU = &CPUStats{
			NumThrottledPeriods: cpu.GetNrThrottled(),
			NumRunnablePeriods:  cpu.GetNrPeriods(),
			UsageInUsec:         cpu.GetUsageUsec(),
			SystemTimeInUsec:    cpu.GetSystemUsec(),
			UserTimeInUsec:      cpu.GetUserUsec(),
			ThrottledTimeInUsec: cpu.GetThrottledUsec(),
		}

		pids, err := mgr.Procs(true)
		if err != nil {
			fmt.Println("Error fetching cgroup processes: ", err)
			return
		}

		for _, pid := range pids {
			if stat, pidErr := pidusage.GetStat(int(pid)); pidErr == nil {
				cgroupStats.CPU.Usage += stat.CPU
			} else {
				fmt.Println("Error fetching CPU usage for PID: ", pid)
			}
		}
	}
}

func (c *CgroupStatsProvider) withPids(pids *stats.PidsStat) CgroupStatsOpt {
	return func(cgroupStats *CgroupStats) {
		cgroupStats.PID = &PidStats{
			Current: pids.GetCurrent(),
			Limit:   pids.GetLimit(),
		}
	}
}

func (c *CgroupStatsProvider) withMemory(memory *stats.MemoryStat) CgroupStatsOpt {
	return func(cgroupStats *CgroupStats) {
		cgroupStats.Memory = &MemoryStats{
			Usage:       memory.GetUsage(),
			UsageLimit:  memory.GetUsageLimit(),
			Unevictable: memory.GetUnevictable(),
			Anon: &AnonymousMemoryStats{
				Total:                memory.GetAnon(),
				Active:               memory.GetActiveAnon(),
				Inactive:             memory.GetInactiveAnon(),
				TransparentHugepages: memory.GetAnonThp(),
			},
			PageCache: &PageCacheStats{
				Activate:   memory.GetPgactivate(),
				Deactivate: memory.GetPgdeactivate(),
				Fault:      memory.GetPgfault(),
				LazyFree:   memory.GetPglazyfree(),
				LazyFreed:  memory.GetPglazyfreed(),
				MajorFault: memory.GetPgmajfault(),
				Refill:     memory.GetPgrefill(),
				Scan:       memory.GetPgscan(),
				Steal:      memory.GetPgsteal(),
			},
			Kernel: &KernelMemoryStats{
				Slab:              memory.GetSlab(),
				SlabReclaimable:   memory.GetSlabReclaimable(),
				SlabUnreclaimable: memory.GetSlabUnreclaimable(),
				Stack:             memory.GetKernelStack(),
			},
			Network: &NetworkMemoryStats{
				Socket: memory.GetSock(),
			},
			Swap: &SwapMemoryStats{
				Limit: memory.GetSwapLimit(),
				Usage: memory.GetSwapUsage(),
			},
			Filesystem: &FilesystemMemoryStats{
				Current:   memory.GetFile(),
				Active:    memory.GetActiveFile(),
				Inactive:  memory.GetInactiveFile(),
				Dirty:     memory.GetFileDirty(),
				Mapped:    memory.GetFileMapped(),
				Writeback: memory.GetFileWriteback(),
				Shmem:     memory.GetShmem(),
			},
			Workingset: &WorkingsetMemoryStats{
				Refault:     memory.GetWorkingsetRefault(),
				Activate:    memory.GetWorkingsetActivate(),
				Nodereclaim: memory.GetWorkingsetNodereclaim(),
			},
			TransparentHugepage: &TransparentHugepageMemoryStats{
				TransparentHugepageFaultAlloc:    memory.GetThpFaultAlloc(),
				TransparentHugepageCollapseAlloc: memory.GetThpCollapseAlloc(),
			},
		}
	}
}

func (c *CgroupStatsProvider) withMemoryEvents(memoryEvents *stats.MemoryEvents) CgroupStatsOpt {
	return func(cgroupStats *CgroupStats) {
		cgroupStats.MemoryEvent = &MemoryEventStats{
			NumOomEvents:     memoryEvents.GetOom(),
			NumOomKillEvents: memoryEvents.GetOomKill(),
			High:             memoryEvents.GetHigh(),
			Max:              memoryEvents.GetMax(),
			Low:              memoryEvents.GetLow(),
		}
	}
}
