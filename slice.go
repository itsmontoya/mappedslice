package mappedslice

import (
	"fmt"
	"os"
	"unsafe"

	"github.com/edsrzf/mmap-go"
)

const mappedCeiling = 1<<37 - 1

func New[T any](filepath string) (ref *Slice[T], err error) {
	var d Slice[T]
	if d.f, err = os.OpenFile(filepath, os.O_CREATE|os.O_RDWR, 0644); err != nil {
		return
	}

	var t T
	d.sizeOf = int64(unsafe.Sizeof(t))

	var info os.FileInfo
	if info, err = d.f.Stat(); err != nil {
		return
	}

	if info.Size() == 0 {
		if err = d.growTo(32); err != nil {
			return
		}
	} else {
		d.associate()
	}

	return &d, nil
}

type Slice[T any] struct {
	f  *os.File
	mm mmap.MMap
	s  []T

	sizeOf int64

	len *int64
	cap *int64
}

func (s *Slice[T]) Get(index int) (kv T, err error) {
	if err = s.boundsCheck(index); err != nil {
		return
	}

	kv = s.s[index]
	return
}

func (s *Slice[T]) Set(index int, t T) (err error) {
	if err = s.boundsCheck(index); err != nil {
		return
	}

	s.s[index] = t
	return
}

func (s *Slice[T]) Append(t T) (err error) {
	if err = s.grow(); err != nil {
		return
	}

	s.s = append(s.s, t)
	*s.len++
	return
}

func (s *Slice[T]) InsertAt(index int, value T) (err error) {
	if err = s.boundsCheck(index); err != nil {
		return
	}

	if err = s.Append(value); err != nil {
		return
	}

	copy(s.s[index+1:], s.s[index:])
	s.s[index] = value
	return
}

func (s *Slice[T]) RemoveAt(index int) (err error) {
	if err = s.boundsCheck(index); err != nil {
		return
	}

	first := s.s[:index]
	second := s.s[index+1:]
	s.s = append(first, second...)
	*s.len--
	return
}

func (s *Slice[T]) ForEach(fn func(T) (end bool)) (ended bool) {
	for _, t := range s.s {
		if ended = fn(t); ended {
			return
		}
	}

	return
}

func (s *Slice[T]) Len() int {
	return len(s.s)
}

func (s *Slice[T]) Slice() (out []T) {
	out = make([]T, len(s.s))
	copy(out, s.s)
	return
}

func (s *Slice[T]) Close() (err error) {
	if err = s.unmap(); err != nil {
		return
	}

	return s.f.Close()
}

func (s *Slice[T]) unmap() (err error) {
	if s.mm == nil {
		return
	}

	if err = s.mm.Unmap(); err != nil {
		return
	}

	s.mm = nil
	return
}

func (s *Slice[T]) associate() (err error) {
	if s.mm, err = mmap.Map(s.f, os.O_RDWR, 0); err != nil {
		return
	}

	s.len = (*int64)(unsafe.Pointer(&s.mm[0]))
	s.cap = (*int64)(unsafe.Pointer(&s.mm[8]))
	s.s = (*(*[mappedCeiling]T)(unsafe.Pointer(&s.mm[16])))[:*s.len]
	return
}

func (s *Slice[T]) grow() (err error) {
	if *s.len < *s.cap {
		return
	}

	newCapacity := *s.cap * 2
	if newCapacity == 0 {
		newCapacity = 32
	}

	return s.growTo(newCapacity)
}

func (s *Slice[T]) growTo(newCapacity int64) (err error) {
	if err = s.unmap(); err != nil {
		return
	}

	newSize := 16 + newCapacity*s.sizeOf
	if err = s.f.Truncate(newSize); err != nil {
		return
	}

	if err = s.associate(); err != nil {
		return
	}

	*s.cap = newCapacity
	return
}

func (s *Slice[T]) boundsCheck(index int) (err error) {
	if s.isInBounds(index) {
		return
	}

	return fmt.Errorf("index of <%d> is out of bounds with a length of <%d>", index, *s.len)
}

func (s *Slice[T]) isInBounds(index int) (ok bool) {
	switch {
	case index > int(*s.len):
		return false
	case index < 0:
		return false
	default:
		return true
	}
}
