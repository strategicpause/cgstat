package v2

import (
	"fmt"
	"github.com/containerd/cgroups/v3/cgroup2"
	"github.com/containerd/cgroups/v3/cgroup2/stats"
	"github.com/prometheus/procfs"
	"github.com/strategicpause/cgstat/stats/common"
	"github.com/struCoder/pidusage"
)

const (
	CgroupPrefix = "/sys/fs/cgroup"
	// SocketPageSizeInBytes tells us the size of pages which are allocated to either TCP or UDP.
	SocketPageSizeInBytes = 4096
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
		cgroupStats, err := c.getStatsByCgroupPath(cgroupPath)
		if err != nil {
			return nil, err
		}
		statsCollection = append(statsCollection, cgroupStats)
	}

	return statsCollection, nil
}

func (c *CgroupStatsProvider) getStatsByCgroupPath(cgroupPath string) (*CgroupStats, error) {
	mgr, err := cgroup2.Load(cgroupPath)
	if err != nil {
		return nil, err
	}
	metrics, err := mgr.Stat()
	if err != nil {
		return nil, err
	}
	return NewCgroupStat(cgroupPath,
		c.withCPU(mgr, metrics.GetCPU()),
		c.withPids(metrics.GetPids()),
		c.withProcStats(mgr),
		c.withMemory(metrics.GetMemory()),
		c.withMemoryEvents(metrics.GetMemoryEvents()),
		c.withNetwork(mgr),
	), nil
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

func (c *CgroupStatsProvider) withProcStats(mgr *cgroup2.Manager) CgroupStatsOpt {
	return func(cgroupStats *CgroupStats) {
		pids, err := mgr.Procs(true)
		if err != nil {
			return
		}

		procStats := ProcStats{}
		for _, pid := range pids {
			if proc, err := procfs.NewProc(int(pid)); err == nil {
				fds, _ := proc.FileDescriptorsLen()
				procStats.NumFD += uint64(fds)

			}
		}
		cgroupStats.ProcStats = &procStats
	}
}

func (c *CgroupStatsProvider) withNetwork(mgr *cgroup2.Manager) CgroupStatsOpt {
	return func(cgroupStats *CgroupStats) {
		pids, err := mgr.Procs(true)
		if err != nil {
			return
		}

		tcpStats := &TCPNetworkStats{}
		udpStats := &UDPNetworkStats{}

		for _, pid := range pids {
			procPath := fmt.Sprintf("/proc/%d", pid)
			if fs, err := procfs.NewFS(procPath); err == nil {
				if tcpSummary, err := fs.NetTCPSummary(); err == nil {
					tcpStats.TxQueueLength += tcpSummary.TxQueueLength
					tcpStats.RxQueueLength += tcpSummary.RxQueueLength
				}
				if udpSummary, err := fs.NetUDPSummary(); err == nil {
					udpStats.TxQueueLength += udpSummary.TxQueueLength
					udpStats.RxQueueLength += udpSummary.RxQueueLength
				}
				if netSockStat, err := fs.NetSockstat(); err == nil {
					for _, protocol := range netSockStat.Protocols {
						if protocol.Protocol == "TCP" {
							tcpStats.Sockets += uint64(protocol.InUse)
							// The Mem value is reported in pages, so we need to convert pages to bytes in order to
							// determine how much memory is being used.
							tcpStats.SocketMemory += uint64(*protocol.Mem * SocketPageSizeInBytes)
						} else if protocol.Protocol == "UDP" {
							udpStats.Sockets += uint64(protocol.InUse)
							udpStats.SocketMemory += uint64(*protocol.Mem * SocketPageSizeInBytes)
						}
					}
				}
			}
		}

		cgroupStats.Network = &NetworkStats{
			TCPStats: tcpStats,
			UDPStats: udpStats,
		}
	}
}
