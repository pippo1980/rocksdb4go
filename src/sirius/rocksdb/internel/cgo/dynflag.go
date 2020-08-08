// +build !linux !static

package cgo

// #cgo LDFLAGS: -lrocksdb -lstdc++ -lm -lz -ldl
import "C"
