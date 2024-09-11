package mappedslice

import (
	"fmt"
	"os"
	"unsafe"

	"github.com/edsrzf/mmap-go"
)

var (
	intSize    = int64(unsafe.Sizeof(int64(0)))
	lenOffset  = 0
	capOffset  = intSize
	dataOffset = intSize * 2
)

func New[T any](filepath string, initialCapacity int64) (ref *Slice[T], err error) {
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
		if err = d.growTo(initialCapacity); err != nil {
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

	sizeOf int64

	len *int64
	cap *int64
}

func (s *Slice[T]) Get(index int) (v T, ok bool) {
	if !s.isInBounds(index) {
		return
	}

	v = *s.get(index)
	ok = true
	return
}

func (s *Slice[T]) Set(index int, t T) (err error) {
	if err = s.boundsCheck(index); err != nil {
		return
	}

	v := s.get(index)
	*v = t
	return
}

func (s *Slice[T]) Append(t T) (err error) {
	if err = s.grow(); err != nil {
		return
	}

	index := int(*s.len)
	*s.len++
	v := s.get(index)
	*v = t
	return
}

func (s *Slice[T]) InsertAt(index int, value T) (err error) {
	if err = s.boundsCheck(index); err != nil {
		return
	}

	if err = s.Append(value); err != nil {
		return
	}

	s.mm = append(s.mm[:])

	for i := int(*s.len - 1); i >= index; i-- {
		var source *T
		if i > index {
			source = s.get(i - 1)
		} else {
			source = &value
		}

		target := s.get(i)
		*target = *source
	}

	return
}

func (s *Slice[T]) RemoveAt(index int) (err error) {
	if err = s.boundsCheck(index); err != nil {
		return
	}

	end := int(*s.len) - 1
	for i := index; i < end; i++ {
		source := s.get(i + 1)
		target := s.get(i)
		*target = *source
	}

	*s.len--
	return
}

func (s *Slice[T]) ForEach(fn func(T) (end bool)) (ended bool) {
	for i := 0; i < int(*s.len); i++ {
		t := s.get(i)
		if ended = fn(*t); ended {
			return
		}
	}

	return
}

func (s *Slice[T]) Cursor() (out Cursor[T]) {
	var c cursor[T]
	c.s = s
	return &c
}

func (s *Slice[T]) Len() int {
	return int(*s.len)
}

func (s *Slice[T]) Slice() (out []T) {
	out = make([]T, 0, *s.len)
	s.ForEach(func(t T) (end bool) {
		out = append(out, t)
		return false
	})

	return
}

func (s *Slice[T]) Close() (err error) {
	if err = s.unmap(); err != nil {
		return
	}

	return s.f.Close()
}

func (s *Slice[T]) get(index int) (t *T) {
	bs := s.getBytes(index)
	t = (*T)(unsafe.Pointer(&bs[0]))
	return
}

func (s *Slice[T]) getBytes(index int) (bs []byte) {
	byteIndex := dataOffset + (int64(index) * s.sizeOf)
	bs = s.mm[byteIndex : byteIndex+s.sizeOf]
	return
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

	s.len = (*int64)(unsafe.Pointer(&s.mm[lenOffset]))
	s.cap = (*int64)(unsafe.Pointer(&s.mm[capOffset]))
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
	case index >= int(*s.len):
		return false
	case index < 0:
		return false
	default:
		return true
	}
}
