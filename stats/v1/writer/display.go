package writer

import (
	"github.com/gosuri/uilive"
	"github.com/rodaine/table"
	"github.com/strategicpause/cgstat/stats/common"
)

type DisplayVerbosity int

const (
	Normal  DisplayVerbosity = 0
	Verbose DisplayVerbosity = 1
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

func (c *CgStatsDisplayWriter) Write(cgroupStats []*common.CgroupStats) error {
	tbl := table.New("Name", "CPU", "NumProcesses", "CurrentUsage", "MaxUsage", "UsageLimit", "RSS",
		"Cache", "Dirty", "WriteBack", "UnderOom", "OomKill")
	tbl.WithWriter(c.writer)

	for _, cgStats := range cgroupStats {
		row := cgStats.ToDisplayRow()
		// Write to display
		tbl.AddRow(row...)
	}
	tbl.Print()
	c.writer.Flush()

	return nil
}
