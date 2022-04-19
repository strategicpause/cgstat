package writer

import "cgstat/stats"

type NullWriter struct {

}

func NewNullWriter() StatsWriter {
	return &NullWriter{}
}

func (n *NullWriter) Write(_ []*stats.CgroupStats) error {
	return nil
}
