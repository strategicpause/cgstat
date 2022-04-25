package writer

import (
	"github.com/strategicpause/cgstat/stats"

	"github.com/gosuri/uilive"
	"github.com/rodaine/table"
)

// CgStatsDisplayWriter will display stats for a set of cgroups to the screen
type CgStatsDisplayWriter struct {
	writer *uilive.Writer
}

func NewCgStatsDisplayWriter() StatsWriter {
	writer := uilive.New()
	writer.Start()

	return &CgStatsDisplayWriter{
		writer: writer,
	}
}

func (c *CgStatsDisplayWriter) Write(cgroupStats []*stats.CgroupStats) error {
	tbl := table.New("Name", "CPU", "NumProcesses", "CurrentUsage", "MaxUsage", "UsageLimit", "RSS",
		"Cache", "Dirty", "WriteBack", "UnderOom", "OomKill")
	tbl.WithWriter(c.writer)

	for _, cgStats := range cgroupStats {
		row := ToDisplayRow(cgStats)
		// Write to display
		tbl.AddRow(row...)
	}
	tbl.Print()
	c.writer.Flush()

	return nil
}
