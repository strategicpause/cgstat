package common

import "fmt"

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

func DisplayRatio(current uint64, total uint64, withTotal bool) string {
	usageRatio := float64(current) / float64(total)
	if withTotal {
		return fmt.Sprintf("%d / %d (%.2f%%)", current, total, usageRatio)
	}
	return fmt.Sprintf("%d (%.2f%%)", current, usageRatio)
}
