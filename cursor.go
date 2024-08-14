package mappedslice

var _ Cursor[int] = &cursor[int]{}

type Cursor[T any] interface {
	Seek(index int) (T, bool)
	Next() (T, bool)
	Prev() (T, bool)
	Close() error
}

type cursor[T any] struct {
	index int
	s     *Slice[T]
}

func (c *cursor[T]) Seek(index int) (t T, ok bool) {
	c.index = index
	if !c.s.isInBounds(c.index) {
		return
	}

	return c.s.s[c.index], true
}

func (c *cursor[T]) Next() (next T, ok bool) {
	c.index++
	if !c.s.isInBounds(c.index) {
		return
	}

	next = c.s.s[c.index]
	ok = true
	return
}

func (c *cursor[T]) Prev() (prev T, ok bool) {
	c.index--
	if !c.s.isInBounds(c.index) {
		return
	}

	prev = c.s.s[c.index]
	ok = true
	return
}

func (c *cursor[T]) Close() error {
	c.s = nil
	return nil
}
