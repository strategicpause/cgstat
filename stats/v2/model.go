package v2

type CPUStats struct {
	// SystemTime in Microseconds.
	SystemTime int64
	// CPU Utilization (percentage)
	Utilization float64
	// Number of runnable periods in which the application used its entire quota and was throttled
	NumThrottledPeriods uint64
	// Number of periods that any process in the cgroup was runnable
	NumRunnablePeriods uint64
	// The total amount of time, in microseconds, processes within the cgroup were throttled
	ThrottledTimeInUsec uint64
	// Total CPU usage, in microseconds.
	UsageInUsec uint64
	// System CPU usage, in microseconds.
	SystemTimeInUsec uint64
	// Userspace CPU usage, in microseconds.
	UserTimeInUsec uint64
}

type ProcStats struct {
	// The total number of open file descriptors for processes in the container
	NumFD uint64
}

type PidStats struct {
	// The number of processes currently in the cgroup and its descendants.
	Current uint64
	// Hard limit of number of processes.
	Limit uint64
}

// AnonymousMemoryStats describes memory that is not backed by a filesystem. This is memory which is created for a
// program's stack & heap or by explicit calls to the mmap system call.
type AnonymousMemoryStats struct {
	// Total amount of anonymous memory being used.
	Total uint64
	// Anonymous memory that has been used more recently and usually not swapped out.
	Active uint64
	// Anonymous memory that has not been used recently and can be swapped out.
	Inactive uint64
	// Amount of memory used in anonymous mappings backed by transparent hugepages.
	TransparentHugepages uint64
}

// PageCacheStats describes the memory which is read from and written to disk. When data is written to disk it is
// written to the cache, and then written to the backing storage device. The page cache is an LRU cache, so that
// least recently used entries are evicted and frequently used entries are kept to avoid expensive reads from the
// storage device.
type PageCacheStats struct {
	// Amount of pages moved to the active LRU list.
	Activate uint64
	// Amount of pages moved to the inactive LRU list.
	Deactivate uint64
	// Total number of page faults incurred.
	Fault uint64
	// Amount of pages postponed to be freed under memory pressure.
	LazyFree uint64
	// Amount of reclaimed lazyfree pages.
	LazyFreed uint64
	// Number of major page faults incurred.
	MajorFault uint64
	// Amount of scanned pages (in an active LRU list).
	Refill uint64
	// Amount of scanned pages (in an inactive LRU list).
	Scan uint64
	// Amount of reclaimed pages.
	Steal uint64
}

type KernelMemoryStats struct {
	// Amount of memory used for storing in-kernel data structures
	Slab uint64
	// Part of “slab” that might be reclaimed, such as dentries and inodes.
	SlabReclaimable uint64
	// Part of “slab” that cannot be reclaimed on memory pressure.
	SlabUnreclaimable uint64
	// The memory the kernel stack uses. This is not reclaimable.
	Stack uint64
}

type NetworkMemoryStats struct {
	// Amount of memory used in network transmission buffers.
	Socket uint64
}

type SwapMemoryStats struct {
	// Swap usage hard limit.  If a cgroup's swap usage reaches this limit, anonymous memory of the cgroup will not
	// be swapped out.
	Limit uint64
	// The total amount of swap currently being used by the cgroup and its descendants.
	Usage uint64
}

type FilesystemMemoryStats struct {
	// Amount of memory used to cache filesystem data, including tmpfs and shared memory.
	Current uint64
	// Pagecache memory that has been used more recently and usually not reclaimed until needed.
	Active uint64
	// Pagecache memory that can be reclaimed without huge performance impact.
	Inactive uint64
	// Amount of cached filesystem data that was modified but not yet written back to disk.
	Dirty uint64
	// Amount of cached filesystem data mapped with mmap()
	Mapped uint64
	// Amount of cached filesystem data that was modified and is currently being written back to disk.
	Writeback uint64
	// Amount of cached filesystem data that is swap-backed, such as tmpfs, shm segments, shared anonymous mmap()s.
	Shmem uint64
}

type WorkingsetMemoryStats struct {
	// Number of refaults of previously evicted pages.
	Refault uint64
	// Number of refaulted pages that were immediately activated.
	Activate uint64
	// Number of times a shadow node has been reclaimed.
	Nodereclaim uint64
}

type TransparentHugepageMemoryStats struct {
	// Number of transparent hugepages which were allocated to satisfy a page fault.
	TransparentHugepageFaultAlloc uint64
	// Number of transparent hugepages which were allocated to allow collapsing an existing range of pages.
	TransparentHugepageCollapseAlloc uint64
}

type MemoryEventStats struct {
	// The number of times the cgroup's memory usage reached the limit and allocation was about to fail.
	// Depending on context, the result could be invoking the OOM killer and retrying allocation, or
	// failing allocation.
	NumOomEvents uint64
	// The number of processes in this cgroup or its subtree killed by any kind of OOM killer. This could be
	// because of a breach of the cgroup’s memory limit, one of its ancestors’ memory limits, or an overall
	// system memory shortage.
	NumOomKillEvents uint64
	// The number of times processes of the cgroup are throttled and routed to perform direct memory reclaim
	// because the high memory boundary was exceeded.  For a cgroup whose memory usage is capped by the high limit
	// rather than global memory pressure, this event's occurrences are expected.
	High uint64
	// The number of times the cgroup's memory usage was about to go over the max boundary.  If direct reclaim
	// fails to bring it down, the cgroup goes to OOM state.
	Max uint64
	// The number of times the cgroup is reclaimed due to high memory pressure even though its usage is under
	// the low boundary.  This usually indicates that the low boundary is over-committed.
	Low uint64
}

type MemoryStats struct {
	// The total amount of memory currently being used by the cgroup and its descendants.
	Usage uint64
	// The maximum amount of memory that can be used by the cgroup and its descendants.
	UsageLimit uint64
	// Memory that cannot be reclaimed.
	Unevictable uint64
	//
	Anon *AnonymousMemoryStats
	//
	PageCache *PageCacheStats
	//
	Kernel *KernelMemoryStats
	//
	Network *NetworkMemoryStats
	//
	Swap *SwapMemoryStats
	//
	Filesystem *FilesystemMemoryStats
	//
	Workingset *WorkingsetMemoryStats
	//
	TransparentHugepage *TransparentHugepageMemoryStats
}

type TCPNetworkStats struct {
	// Number of TCP sockets which are not in the CLOSED state.
	Sockets uint64
	// Amount of bytes allocated to TCP.
	SocketMemory uint64
	// TCP receive queue length.
	RxQueueLength uint64
	// TCP transfer queue length.
	TxQueueLength uint64
}

type UDPNetworkStats struct {
	// Number of UDP sockets which are not in the CLOSED state.
	Sockets uint64
	// Amount of bytes allocated to UDP.
	SocketMemory uint64
	// UDP receive queue length.
	RxQueueLength uint64
	// UDP transfer queue length.
	TxQueueLength uint64
}

type NetworkStats struct {
	TCPStats *TCPNetworkStats
	UDPStats *UDPNetworkStats
}

type CgroupStats struct {
	//
	Name string
	//
	CPU *CPUStats
	//
	PID *PidStats
	//
	ProcStats *ProcStats
	//
	Memory *MemoryStats
	//
	MemoryEvent *MemoryEventStats
	//
	Network *NetworkStats
}

type CgroupStatsOpt func(*CgroupStats)

func NewCgroupStat(name string, opts ...CgroupStatsOpt) *CgroupStats {
	cgroupStats := &CgroupStats{
		Name: name,
	}
	for _, opt := range opts {
		opt(cgroupStats)
	}
	return cgroupStats
}
