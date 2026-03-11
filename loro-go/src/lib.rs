#[cfg(not(target_env = "msvc"))]
#[global_allocator]
static GLOBAL: tikv_jemallocator::Jemalloc = tikv_jemallocator::Jemalloc;

loro_ffi::uniffi_reexport_scaffolding!();

#[cfg(not(target_env = "msvc"))]
mod jemalloc_profiling {
    #[repr(C)]
    pub struct JemallocProfileResult {
        pub data: *mut u8,
        pub len: usize,
        pub error_code: i32,
    }

    #[no_mangle]
    pub extern "C" fn loro_jemalloc_prof_enabled() -> i32 {
        if jemalloc_pprof::PROF_CTL.as_ref().is_some() {
            1
        } else {
            0
        }
    }

    #[no_mangle]
    pub extern "C" fn loro_jemalloc_dump_pprof() -> JemallocProfileResult {
        let Some(prof_ctl) = jemalloc_pprof::PROF_CTL.as_ref() else {
            return JemallocProfileResult {
                data: std::ptr::null_mut(),
                len: 0,
                error_code: 1,
            };
        };

        let mut guard = match prof_ctl.try_lock() {
            Ok(g) => g,
            Err(_) => {
                return JemallocProfileResult {
                    data: std::ptr::null_mut(),
                    len: 0,
                    error_code: 2,
                };
            }
        };

        match guard.dump_pprof() {
            Ok(pprof_data) => {
                let len = pprof_data.len();
                let mut boxed = pprof_data.into_boxed_slice();
                let data = boxed.as_mut_ptr();
                std::mem::forget(boxed);
                JemallocProfileResult {
                    data,
                    len,
                    error_code: 0,
                }
            }
            Err(_) => JemallocProfileResult {
                data: std::ptr::null_mut(),
                len: 0,
                error_code: 3,
            },
        }
    }

    #[no_mangle]
    pub unsafe extern "C" fn loro_jemalloc_free_profile(data: *mut u8, len: usize) {
        if !data.is_null() && len > 0 {
            let _ = Vec::from_raw_parts(data, len, len);
        }
    }

    #[repr(C)]
    pub struct JemallocStats {
        pub allocated: usize,
        pub active: usize,
        pub resident: usize,
        pub mapped: usize,
        pub retained: usize,
        pub error_code: i32,
    }

    #[no_mangle]
    pub extern "C" fn loro_jemalloc_stats() -> JemallocStats {
        use tikv_jemalloc_ctl::{epoch, stats};

        // Advance the epoch so stats are fresh.
        if epoch::advance().is_err() {
            return JemallocStats {
                allocated: 0,
                active: 0,
                resident: 0,
                mapped: 0,
                retained: 0,
                error_code: 1,
            };
        }

        let allocated = stats::allocated::read().unwrap_or(0);
        let active = stats::active::read().unwrap_or(0);
        let resident = stats::resident::read().unwrap_or(0);
        let mapped = stats::mapped::read().unwrap_or(0);
        let retained = stats::retained::read().unwrap_or(0);

        JemallocStats {
            allocated,
            active,
            resident,
            mapped,
            retained,
            error_code: 0,
        }
    }
}
