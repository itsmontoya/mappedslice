package mappedslice

type Cursor[T any] struct {
	index int
	s     *Slice[T]
}

func (c *Cursor[T]) Seek(index int) (t T, ok bool) {
	c.index = index
	if !c.s.isInBounds(c.index) {
		return
	}

	return c.s.s[c.index], true
}

func (c *Cursor[T]) Next() (next T, ok bool) {
	c.index++
	if !c.s.isInBounds(c.index) {
		return
	}

	next = c.s.s[c.index]
	ok = true
	return
}

func (c *Cursor[T]) Prev() (prev T, ok bool) {
	c.index--
	if !c.s.isInBounds(c.index) {
		return
	}

	prev = c.s.s[c.index]
	ok = true
	return
}
