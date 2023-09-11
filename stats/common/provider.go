package common

import (
	"io/fs"
	"path/filepath"
	"strings"
)

type CommonCgroupStatsProvider struct {
	cgroupRootDir    string
	cgroupRootDirLen int
}

func NewCommonCgroupStatsProvider(rootDir string) *CommonCgroupStatsProvider {
	return &CommonCgroupStatsProvider{
		cgroupRootDir:    rootDir,
		cgroupRootDirLen: len(rootDir),
	}
}

func (c *CommonCgroupStatsProvider) ListCgroupsByPrefix(cgroupPrefix string) []string {
	var cgroupPaths []string
	queue := []string{cgroupPrefix}

	for len(queue) > 0 {
		prefix := queue[0]
		queue = queue[1:]

		prefixPath := filepath.Join(c.cgroupRootDir, prefix)

		_ = filepath.WalkDir(filepath.Dir(prefixPath), func(currPath string, d fs.DirEntry, err error) error {
			if d.IsDir() && strings.HasPrefix(currPath, prefixPath) {
				cgroupPaths = append(cgroupPaths, currPath[c.cgroupRootDirLen:])
			}
			return nil
		})
	}
	return cgroupPaths
}
