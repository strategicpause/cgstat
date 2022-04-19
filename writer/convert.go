package writer

import (
	"cgstat/stats"
	"fmt"
	"time"
)

func FormatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(bytes)/float64(div), "KMGTPE"[exp])
}


func ToCsvRow(cgStats *stats.CgroupStats) []string {
	t, _  := time.Now().UTC().MarshalText()
	return []string {
		string(t),
		cgStats.Name,
		fmt.Sprintf("%f", cgStats.UserCPU),
		fmt.Sprintf("%f", cgStats.KernelCPU),
		fmt.Sprintf("%d", cgStats.CurrentUsage),
		fmt.Sprintf("%d", cgStats.MaxUsage),
		fmt.Sprintf("%d", cgStats.UsageLimit),
		fmt.Sprintf("%d", cgStats.Rss),
		fmt.Sprintf("%d", cgStats.CacheSize),
		fmt.Sprintf("%d", cgStats.DirtySize),
		fmt.Sprintf("%d", cgStats.WriteBack),
		fmt.Sprintf("%d", cgStats.UnderOom),
		fmt.Sprintf("%d", cgStats.OomKill),
	}
}

func ToDisplayRow(cgStats *stats.CgroupStats) []interface{} {
	userCPU := fmt.Sprintf("%.2f%%", cgStats.UserCPU)
	kernelCPU := fmt.Sprintf("%.2f%%", cgStats.KernelCPU)
	currentUsage := fmt.Sprintf("%s (%.2f%%)", FormatBytes(cgStats.CurrentUsage), cgStats.CurrentUtilization)
	maxUsage := fmt.Sprintf("%s (%.2f%%)", FormatBytes(cgStats.MaxUsage), cgStats.CurrentUtilization)
	usageLimit := FormatBytes(cgStats.UsageLimit)
	rss := FormatBytes(cgStats.Rss)
	cacheSize := FormatBytes(cgStats.CacheSize)
	dirtySize := FormatBytes(cgStats.DirtySize)
	writeback := FormatBytes(cgStats.WriteBack)

	return []interface{} { cgStats.Name, userCPU, kernelCPU, currentUsage, maxUsage, usageLimit, rss,
		cacheSize, dirtySize, writeback, cgStats.UnderOom, cgStats.OomKill }
}