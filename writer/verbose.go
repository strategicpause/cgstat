package writer

import (
	"github.com/gosuri/uilive"
	"github.com/strategicpause/cgstat/stats/common"
)

type CgStatsVerboseWriter struct {
	writer *uilive.Writer
}

func NewCgroupStatsVerboseWriter() StatsWriter {
	writer := uilive.New()
	writer.Start()

	return &CgStatsVerboseWriter{
		writer: writer,
	}
}

func (c *CgStatsVerboseWriter) Write(cgStats common.CgroupStatsCollection) error {
	for _, s := range cgStats.GetCgroupStats() {
		s.ToVerboseOutput(c.writer)
	}
	return c.writer.Flush()
}
