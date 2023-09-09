package common

import "io"

type CgroupStatsProvider interface {
	ListCgroupsByPrefix(cgroupPrefix string) []string
	GetCgroupStatsByPrefix(prefix string) (CgroupStatsCollection, error)
	GetCgroupStatsByName(name string) (CgroupStatsCollection, error)
}

type CgroupStatsCollection interface {
	ToCsvOutput() *CsvOutput
	ToDisplayOutput() *DisplayOutput
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
