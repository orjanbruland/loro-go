package loro

/*
#include "jemalloc_pprof.h"
*/
import "C"

import (
	"errors"
	"net/http"
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
			return nil, errors.New("jemalloc profiling not enabled (set MALLOC_CONF=prof:true,prof_active:true)")
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

// JemallocProfileHandler returns an HTTP handler that serves jemalloc heap
// profiles in pprof format. Compatible with go tool pprof.
//
// Usage:
//
//	mux.Handle("/debug/pprof/jemalloc", loro.JemallocProfileHandler())
func JemallocProfileHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := DumpJemallocProfile()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment; filename=jemalloc.pb.gz")
		w.Write(data)
	})
}
