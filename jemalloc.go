package loro

import (
	"encoding/json"
	"net/http"
)

// JemallocStats contains jemalloc memory usage statistics.
type JemallocStats struct {
	Allocated uint64 `json:"allocated"` // bytes actively in use by the application
	Active    uint64 `json:"active"`    // bytes in active pages (jemalloc's working set)
	Resident  uint64 `json:"resident"`  // bytes in physically resident pages
	Mapped    uint64 `json:"mapped"`    // bytes in mmap'd regions (total address space)
	Retained  uint64 `json:"retained"`  // bytes in retained (cached) virtual memory
}

// JemallocStatsHandler returns an HTTP handler that serves jemalloc memory
// statistics as JSON.
//
// Usage:
//
//	mux.Handle("/debug/jemalloc/stats", loro.JemallocStatsHandler())
func JemallocStatsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		stats, err := GetJemallocStats()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stats)
	})
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
