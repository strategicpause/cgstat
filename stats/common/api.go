package common

type CgroupStatsProvider interface {
	ListCgroupsByPrefix(cgroupPrefix string) []string
	GetCgroupStatsByPrefix(prefix string) ([]*CgroupStats, error)
	GetCgroupStatsByName(name string) ([]*CgroupStats, error)
}
