package cgo

// #include <stdlib.h>
import "C"
import "unsafe"

// Slice is used as a wrapper for non-copy values
type Slice struct {
	data  *C.char
	size  C.size_t
	freed bool
}

type Slices []*Slice

func (slices Slices) Destroy() {
	for _, s := range slices {
		s.Free()
	}
}

// NewSlice returns a slice with the given data.
func NewSlice(data *C.char, size C.size_t) *Slice {
	return &Slice{data, size, false}
}

// StringToSlice is similar to NewSlice, but can be called with
// a Go string type. This exists to make testing integration
// with Gorocksdb easier.
func StringToSlice(data string) *Slice {
	return NewSlice(C.CString(data), C.size_t(len(data)))
}

// Data returns the data of the slice.
func (s *Slice) Data() []byte {
	return charToByte(s.data, s.size)
}

// Size returns the size of the data.
func (s *Slice) Size() int {
	return int(s.size)
}

// Free frees the slice data.
func (s *Slice) Free() {
	if !s.freed {
		C.free(unsafe.Pointer(s.data))
		s.freed = true
	}
}

func RetrieveSliceData(slice *Slice) []byte {
	if slice == nil {
		return nil
	}

	if slice.Size() == 0 {
		slice.Free()
		return nil
	}

	// 这里要主动copy一次, 否则[]byte和result.Data()共享指针, 有可能在results.Destroy()后返回空数据
	value := make([]byte, slice.Size())
	copy(value, slice.Data())
	slice.Free()
	return value
}
