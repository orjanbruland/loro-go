#ifndef JEMALLOC_PPROF_H
#define JEMALLOC_PPROF_H

#include <stddef.h>
#include <stdint.h>

typedef struct {
    uint8_t *data;
    size_t len;
    int32_t error_code;
} JemallocProfileResult;

// Check if jemalloc profiling is enabled.
// Returns 1 if enabled, 0 if not.
int32_t loro_jemalloc_prof_enabled(void);

// Dump jemalloc heap profile as gzip-compressed pprof protobuf.
// error_code: 0=ok, 1=not enabled, 2=lock busy, 3=dump failed
JemallocProfileResult loro_jemalloc_dump_pprof(void);

// Free a profile buffer previously returned by loro_jemalloc_dump_pprof.
void loro_jemalloc_free_profile(uint8_t *data, size_t len);

#endif
