package cgo

/*
#cgo CXXFLAGS: -std=c++11
//#cgo CPPFLAGS: -I${SRCDIR}/../submodules/facebook/rocksdb/include
//#cgo LDFLAGS: -L${SRCDIR}/../submodules/facebook/rocksdb -lrocksdb
//#cgo CPPFLAGS: -I${SRCDIR}/rocksdb/include
//#cgo LDFLAGS:  -L${SRCDIR}/rocksdb/lib -lrocksdb -lstdc++ -lm -lz
// 与import "C"之间不能有空行，否则就会有这样的error："could not determine kind of name for C.xxxx"
*/
import "C"
