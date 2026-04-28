//go:build !linux

package loro

import (
	"errors"
)

func JemallocProfilingEnabled() bool { return false }

func DumpJemallocProfile() ([]byte, error) {
	return nil, errors.New("jemalloc profiling is not supported on this platform")
}

var emptyStats = &JemallocStats{}

// GetJemallocStats returns current jemalloc memory statistics.
func GetJemallocStats() (*JemallocStats, error) {
	return emptyStats, nil
}
