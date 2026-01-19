package mappedslice

import "fmt"

type BoundsError struct {
	index  int
	length int64
}

func (b *BoundsError) Error() string {
	return fmt.Sprintf("index of <%d> is out of bounds with a length of <%d>", b.index, b.length)
}
