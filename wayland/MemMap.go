package wayland

/*
#include "mmap.h"
*/
import "C"

import (
	"fmt"
	"unsafe"
)

type MemMapInfo struct {
	Bytes          []byte
	Addr           unsafe.Pointer
	Size           C.size_t
	FileDescriptor C.int
	UnMapped       bool
}

func NewMemMapInfo(fd int, size uint64) (MemMapInfo, error) {
	fdNum := C.int(fd)
	c_size := C.size_t(size)
	addr := C.mmap_fd(fdNum, c_size)
	if addr == C.map_failed() {
		return MemMapInfo{
			Addr:           addr,
			Size:           c_size,
			FileDescriptor: fdNum,
			UnMapped:       true,
		}, fmt.Errorf("failed to mmap fd %d", fdNum)
	}

	info := MemMapInfo{
		Addr:           addr,
		Size:           c_size,
		FileDescriptor: fdNum,
		UnMapped:       false,
	}
	info.Bytes = unsafe.Slice((*byte)(info.Addr), info.Size)
	return info, nil
}

func (m *MemMapInfo) Unmap() {
	if m.UnMapped {
		return
	}

	C.unmap(m.Addr, m.Size)
	m.UnMapped = true
	m.Bytes = nil
}
