package mappedslice

import (
	"fmt"
	"os"
	"reflect"
	"unsafe"

	"github.com/edsrzf/mmap-go"
)

const mappedCeiling = 1<<37 - 1

var headerSize = reflect.SliceHeader{}

func New[T any](filepath string) (ref *Mapped[T], err error) {
	var m Mapped[T]
	if m.f, err = os.OpenFile(filepath, os.O_CREATE|os.O_RDWR, 0644); err != nil {
		return
	}

	var info os.FileInfo
	if info, err = m.f.Stat(); err != nil {
		return
	}

	var t T
	m.sizeOf = int64(unsafe.Sizeof(t))
	m.cap = (info.Size() - 8) / m.sizeOf
	if info.Size() == 0 {
		if err = m.Grow(32); err != nil {
			return
		}
	} else {
		m.associate()
	}

	return &m, nil
}

type Mapped[T any] struct {
	f  *os.File
	mm mmap.MMap
	s  []T

	sizeOf int64

	len *int64
	cap int64
}

func (m *Mapped[T]) Len() int {
	return len(m.s)
}

func (m *Mapped[T]) Slice() []T {
	return m.s
}

func (m *Mapped[T]) Close() (err error) {
	if err = m.unmap(); err != nil {
		return
	}

	return m.f.Close()
}

func (m *Mapped[T]) unmap() (err error) {
	if m.mm == nil {
		return
	}

	if err = m.mm.Unmap(); err != nil {
		return
	}

	m.mm = nil
	return
}

func (m *Mapped[T]) associate() (err error) {
	if m.mm, err = mmap.Map(m.f, os.O_RDWR, 0); err != nil {
		return
	}

	m.len = (*int64)(unsafe.Pointer(&m.mm[0]))
	m.s = (*(*[mappedCeiling - 8]T)(unsafe.Pointer(&m.mm[8])))[:*m.len]

	// Use unsafe to manipulate the slice
	header := (*reflect.SliceHeader)(unsafe.Pointer(&m.s))
	header.Cap = int(m.cap)

	fmt.Println("Associated", *m.len, m.s, m.cap)
	return
}

func (m *Mapped[T]) Grow(newCapacity int) (err error) {
	if cap(m.s) >= newCapacity {
		return
	}

	if err = m.unmap(); err != nil {
		return
	}

	newSize := 32 + int64(newCapacity)*m.sizeOf

	if err = m.f.Truncate(newSize); err != nil {
		return
	}

	m.cap = int64(newCapacity)

	if err = m.associate(); err != nil {
		return
	}

	return
}
