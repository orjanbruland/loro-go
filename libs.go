package loro

// #cgo darwin,amd64 LDFLAGS: -L${SRCDIR}/libs/x86_64-apple-darwin -lloro
// #cgo darwin,arm64 LDFLAGS: -L${SRCDIR}/libs/aarch64-apple-darwin -lloro
// #cgo linux,amd64 LDFLAGS: -L${SRCDIR}/libs/x86_64-unknown-linux-musl -lloro -lm
// #cgo linux,arm64 LDFLAGS: -L${SRCDIR}/libs/aarch64-unknown-linux-musl -lloro -lm
// #cgo windows,386 LDFLAGS: -L${SRCDIR}/libs/i686-pc-windows-gnu -lloro -lkernel32 -lbcrypt -lsynchronization
// #cgo windows,amd64 LDFLAGS: -L${SRCDIR}/libs/x86_64-pc-windows-gnu -lloro -lkernel32 -lbcrypt -lsynchronization
import "C"
