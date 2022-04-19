package writer

import (
	"cgstat/stats"
	"encoding/csv"
	"os"
)

type CgroupStatsCsvWriter struct {
	writer *csv.Writer
}

func NewCgroupStatsCsvWriter(args *stats.CgstatArgs) (*CgroupStatsCsvWriter, error) {
	fileWriter, err := os.Create(args.OutputFile)
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
	header := []string{"Time", "Name", "UserCPU", "KernelCPU", "CurrentUsage", "MaxUsage", "UsageLimit", "RSS",
		"Cache", "Dirty", "WriteBack", "UnderOom", "OomKill"}
	return c.writer.Write(header)
}

func (c *CgroupStatsCsvWriter) Write(cgroupStats []*stats.CgroupStats) error {
	for _, s := range cgroupStats {
		err := c.writer.Write(ToCsvRow(s))
		if err != nil {
			return err
		}
	}

	c.writer.Flush()
	return nil
}
