package writer

import (
	"encoding/csv"
	"github.com/strategicpause/cgstat/stats/common"
	"os"
)

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

	err = statsWriter.addHeader()
	if err != nil {
		return nil, err
	}

	return &statsWriter, nil
}

func (c *CgroupStatsCsvWriter) addHeader() error {
	header := []string{"Time", "Name", "UserCPU", "CurrentUsage", "MaxUsage", "UsageLimit", "RSS",
		"Cache", "Dirty", "WriteBack", "UnderOom", "OomKill"}
	return c.writer.Write(header)
}

func (c *CgroupStatsCsvWriter) Write(cgroupStats []*common.CgroupStats) error {
	for _, s := range cgroupStats {
		err := c.writer.Write(s.ToCsvRow())
		if err != nil {
			return err
		}
	}

	c.writer.Flush()
	return nil
}
