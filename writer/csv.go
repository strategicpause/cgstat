package writer

import (
	"encoding/csv"
	"os"

	"github.com/strategicpause/cgstat/stats/common"
)

// CgroupStatsCsvWriter is an implementation of StatsWriter which will write cgroup stats to a CSV file.
type CgroupStatsCsvWriter struct {
	writer *csv.Writer
}

func NewCgroupStatsCsvWriter(filename string) (*CgroupStatsCsvWriter, error) {
	fileWriter, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	statsWriter := CgroupStatsCsvWriter{
		writer: csv.NewWriter(fileWriter),
	}

	return &statsWriter, nil
}

func (c *CgroupStatsCsvWriter) Write(collection common.CgroupStatsCollection) error {
	csvOutput := collection.ToCsvOutput()
	if err := c.writer.Write(csvOutput.Headers); err != nil {
		return err
	}

	for _, row := range csvOutput.Rows {
		if err := c.writer.Write(row); err != nil {
			return err
		}
	}
	c.writer.Flush()

	return nil
}
