package cgo

import "bytes"

type BytesComparator struct{}

func (cmp *BytesComparator) Name() string { return "gorocksdb.bytes-reverse" }
func (cmp *BytesComparator) Compare(a, b []byte) int {
	rtn := bytes.Compare(a, b)
	return rtn
}
