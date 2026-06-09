package loro

/*
#include "jemalloc_pprof.h"
*/
import "C"

import (
	"errors"
	"unsafe"
)

// JemallocProfilingEnabled reports whether jemalloc profiling is active.
// Profiling is enabled via MALLOC_CONF=prof:true,prof_active:true at process startup.
func JemallocProfilingEnabled() bool {
	return C.loro_jemalloc_prof_enabled() == 1
}

// DumpJemallocProfile captures a jemalloc heap profile in pprof protobuf format.
func DumpJemallocProfile() ([]byte, error) {
	result := C.loro_jemalloc_dump_pprof()
	if result.error_code != 0 {
		switch result.error_code {
		case 1:
			return nil, errors.New("jemalloc profiling not enabled (set _RJEM_MALLOC_CONF=prof:true,prof_active:true)")
		case 2:
			return nil, errors.New("jemalloc profiling lock is busy, try again")
		case 3:
			return nil, errors.New("jemalloc profile dump failed")
		default:
			return nil, errors.New("jemalloc profile dump failed with unknown error")
		}
	}
	defer C.loro_jemalloc_free_profile(result.data, result.len)
	return C.GoBytes(unsafe.Pointer(result.data), C.int(result.len)), nil
}

// GetJemallocStats returns current jemalloc memory statistics.
func GetJemallocStats() (*JemallocStats, error) {
	result := C.loro_jemalloc_stats()
	if result.error_code != 0 {
		return nil, errors.New("jemalloc stats: mallctl failed")
	}
	return &JemallocStats{
		Allocated: uint64(result.allocated),
		Active:    uint64(result.active),
		Resident:  uint64(result.resident),
		Mapped:    uint64(result.mapped),
		Retained:  uint64(result.retained),
	}, nil
}
