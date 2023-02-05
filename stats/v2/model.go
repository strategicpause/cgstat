package v2

import (
	"github.com/strategicpause/cgstat/stats/common"
	"io"
)

type CgroupStats struct {
	Name string
	/**     **/
	/** CPU **/
	/**     **/
	//

	/**      **/
	/** PIDs **/
	/**      **/
	// The number of processes currently in the cgroup and its descendants.
	NumPids uint64
	// Hard limit of number of processes.
	MaxPids uint64

	/**        **/
	/** Memory **/
	/**        **/
	// The total amount of memory currently being used by the cgroup and its descendants.
	CurrentUsage uint64
	//
	UsageLimit uint64

	/**               **/
	/** Memory Events **/
	/**               **/
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
	MemoryHigh uint64
	// The number of times the cgroup's memory usage was about to go over the max boundary.  If direct reclaim
	// fails to bring it down, the cgroup goes to OOM state.
	MemoryMax uint64
	// The number of times the cgroup is reclaimed due to high memory pressure even though its usage is under
	// the low boundary.  This usually indicates that the low boundary is over-committed.
	MemoryLow uint64

	/**    **/
	/** IO **/
	/**    **/
	//
}

type CgroupStatsOpt func(*CgroupStats)

func NewCgroupStat(name string, opts ...CgroupStatsOpt) common.CgroupStats {
	cgroupStats := &CgroupStats{
		Name: name,
	}
	for _, opt := range opts {
		opt(cgroupStats)
	}
	return cgroupStats
}

func (c *CgroupStats) ToCsvRow() []string {
	return nil
}

func (c *CgroupStats) ToDisplayRow() []interface{} {
	return nil
}

func (c *CgroupStats) ToVerboseOutput(io.Writer) {

}
