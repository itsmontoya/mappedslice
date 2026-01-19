package mappedslice

var _ Cursor[int] = &cursor[int]{}

type Cursor[T any] interface {
	Seek(index int) (T, error)
	Next() (T, error)
	Prev() (T, error)
	Close() error
}

type cursor[T any] struct {
	index int
	s     *Slice[T]
}

func (c *cursor[T]) Seek(index int) (t T, err error) {
	c.index = index
	if err = c.s.boundsCheck(c.index); err != nil {
		return
	}

	t = c.s.s[c.index]
	return
}

func (c *cursor[T]) Next() (next T, err error) {
	c.index++
	if err = c.s.boundsCheck(c.index); err != nil {
		return
	}

	next = c.s.s[c.index]
	return
}

func (c *cursor[T]) Prev() (prev T, err error) {
	c.index--
	if err = c.s.boundsCheck(c.index); err != nil {
		return
	}

	prev = c.s.s[c.index]
	return
}

func (c *cursor[T]) Close() error {
	c.s = nil
	return nil
}
