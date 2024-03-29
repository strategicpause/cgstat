package writer

import (
	"fmt"
)

func WithCSVWriter(filename string) ViewWriterOptions {
	return func() (StatsWriter, error) {
		return NewCgroupStatsCsvWriter(filename)
	}
}

func WithDisplayWriter(verbosity DisplayVerbosity) ViewWriterOptions {
	return func() (StatsWriter, error) {
		if verbosity == Verbose {
			return NewCgroupStatsVerboseWriter(), nil
		}
		return NewCgStatsDisplayWriter(), nil
	}
}

func NewViewWriters(options []ViewWriterOptions) []StatsWriter {
	var writers []StatsWriter
	for _, opt := range options {
		writer, err := opt()
		if err != nil {
			fmt.Println(err)
			continue
		}
		writers = append(writers, writer)
	}
	return writers
}
