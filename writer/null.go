package writer

import "github.com/strategicpause/cgstat/stats"

// NullWriter implements the StasWriter interface, but does nothing.
type NullWriter struct {
}

func NewNullWriter() StatsWriter {
	return &NullWriter{}
}

func (n *NullWriter) Write(_ []*stats.CgroupStats) error {
	return nil
}
