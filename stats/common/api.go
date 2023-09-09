package common

import "io"

type CgroupStatsProvider interface {
	// ListCgroupsByPrefix will return a list of cgroup names that start with the given prefix.
	ListCgroupsByPrefix(cgroupPrefix string) []string
	// GetCgroupStatsByPrefix will return stats for cgroups that start with the given prefix.
	GetCgroupStatsByPrefix(prefix string) (CgroupStatsCollection, error)
	// GetCgroupStatsByName will return stats for the cgroup that matches the given name.
	GetCgroupStatsByName(name string) (CgroupStatsCollection, error)
}

type CgroupStatsCollection interface {
	// ToCsvOutput will transform the underlying collection into a format which can be written to a CSV file.
	ToCsvOutput() *CsvOutput
	// ToDisplayOutput will transform the underlying collection into a format which can be displayed to the screen.
	ToDisplayOutput() *DisplayOutput
	// ToVerboseOutput will transform the write the given collection to the provided writer. There is no guarantee about
	// the format of the data that is written to the given writer.
	ToVerboseOutput(writer io.Writer)
}

type CsvOutput struct {
	Headers []string
	Rows    [][]string
}

type DisplayOutput struct {
	Headers []interface{}
	Rows    [][]interface{}
}
