// +build static

package cgo

// #cgo LDFLAGS: -l:librocksdb.a -l:libstdc++.a -lm -ldl
import "C"
