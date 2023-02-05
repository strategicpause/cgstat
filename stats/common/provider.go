package common

import (
	"log"
	"os"
	"path/filepath"
)

type CommonCgroupStatsProvider struct {
	cgroupRootDir string
}

func NewCommonCgroupStatsProvider(rootDir string) *CommonCgroupStatsProvider {
	return &CommonCgroupStatsProvider{
		cgroupRootDir: rootDir,
	}
}

func (c *CommonCgroupStatsProvider) ListCgroupsByPrefix(cgroupPrefix string) []string {
	var cgroupPaths []string
	queue := []string{cgroupPrefix}

	for len(queue) > 0 {
		prefix := queue[0]
		queue = queue[1:]

		prefixPath := filepath.Join(c.cgroupRootDir, prefix)
		files, err := os.ReadDir(prefixPath)
		if err != nil {
			log.Println(err)
		}

		for _, file := range files {
			if file.IsDir() {
				cgroupPath := filepath.Join(prefix, file.Name())
				cgroupPaths = append(cgroupPaths, cgroupPath)
				if prefix != cgroupPath {
					queue = append(queue, cgroupPath)
				}
			}
		}
	}
	return cgroupPaths
}
