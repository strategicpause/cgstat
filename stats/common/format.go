package common

import "fmt"

const (
	DividerText = "[...]"
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

// defaultFormat is used when a formatted is not provided. This will simply turn the given uint64 into a string.
func defaultFormat(n uint64) string {
	return fmt.Sprintf("%d", n)
}

type DisplayOptions func(d *displayOptions)

type displayOptions struct {
	withTotal bool
	formatter func(uint64) string
}

func WithTotal() DisplayOptions {
	return func(d *displayOptions) {
		d.withTotal = true
	}
}

func WithBytes() DisplayOptions {
	return func(d *displayOptions) {
		d.formatter = FormatBytes
	}
}

func DisplayRatio(current uint64, total uint64, opts ...DisplayOptions) string {
	dispOpts := displayOptions{}
	for _, opt := range opts {
		opt(&dispOpts)
	}
	if dispOpts.formatter == nil {
		dispOpts.formatter = defaultFormat
	}

	currentStr := dispOpts.formatter(current)

	usageRatio := 0.0
	if total != 0 {
		usageRatio = float64(current) / float64(total) * 100
	}

	if dispOpts.withTotal {
		totalStr := dispOpts.formatter(total)
		return fmt.Sprintf("%s / %s (%.2f%%)", currentStr, totalStr, usageRatio)
	}
	return fmt.Sprintf("%s (%.2f%%)", currentStr, usageRatio)
}

func Shorten(s string, truncateLen int) string {
	if len(DividerText) >= truncateLen {
		return DividerText
	}
	strLen := len(s)
	if strLen <= truncateLen {
		return s
	}
	prefixLen := (truncateLen - len(DividerText)) / 2

	return s[:prefixLen] + DividerText + s[strLen-prefixLen:]
}
