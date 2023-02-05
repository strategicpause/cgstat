package v2

import (
	"github.com/strategicpause/cgstat/stats/common"
)

const (
	CgroupPrefix = "/sys/fs/cgroup"
)

type CgroupStatsProvider struct {
	commonProvider *common.CommonCgroupStatsProvider
}

func NewCgroupStatsProvider() common.CgroupStatsProvider {
	return &CgroupStatsProvider{
		commonProvider: common.NewCommonCgroupStatsProvider(CgroupPrefix),
	}
}

func (c *CgroupStatsProvider) ListCgroupsByPrefix(cgroupPrefix string) []string {
	return c.commonProvider.ListCgroupsByPrefix(cgroupPrefix)
}

func (c *CgroupStatsProvider) GetCgroupStatsByPrefix(prefix string) ([]*common.CgroupStats, error) {
	paths := c.ListCgroupsByPrefix(prefix)
	return c.getCgroupStatsByPath(paths)
}

func (c *CgroupStatsProvider) GetCgroupStatsByName(name string) ([]*common.CgroupStats, error) {
	paths := []string{name}

	return c.getCgroupStatsByPath(paths)
}

func (c *CgroupStatsProvider) getCgroupStatsByPath(cgroupPaths []string) ([]*common.CgroupStats, error) {
	return nil, nil
}
