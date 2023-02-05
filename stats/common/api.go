package common

import "io"

type CgroupStatsProvider interface {
	ListCgroupsByPrefix(cgroupPrefix string) []string
	GetCgroupStatsByPrefix(prefix string) (CgroupStatsCollection, error)
	GetCgroupStatsByName(name string) (CgroupStatsCollection, error)
}

type CgroupStats interface {
	ToCsvRow() []string
	ToDisplayRow() []interface{}
	ToVerboseOutput(io.Writer)
}

type CgroupStatsCollection interface {
	GetCSVHeaders() []string
	GetDisplayHeaders() []string
	GetCgroupStats() []CgroupStats
}
