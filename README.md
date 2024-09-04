# mappedslice [![GoDoc](https://godoc.org/github.com/itsmontoya/mappedslice?status.svg)](https://godoc.org/github.com/itsmontoya/mappedslice) ![Status](https://img.shields.io/badge/status-beta-yellow.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/itsmontoya/mappedslice)](https://goreportcard.com/report/github.com/itsmontoya/mappedslice) ![Go Test Coverage](https://img.shields.io/badge/coverage-86%25-brightgreen)
`mappedslice` is a Go library for efficient, memory map-backed slices. It provides a way to work with large datasets by mapping files directly into memory, allowing for fast and efficient access and manipulation of data. Ideal for applications needing high-performance data handling without loading entire files into RAM.

## Features:
- Memory-mapped slices for efficient file access
- Supports large datasets with minimal memory overhead
- Easy-to-use API for integrating with existing Go applications

## Usage

### New
```go
func ExampleNew() {
	var err error
	if exampleSlice, err = New[int]("myfile.bat", 32); err != nil {
		// Handle error here
		return
	}
}
```

### Slice.Get
```go
func ExampleSlice_Get() {
	var (
		v   int
		ok bool
	)

	if v, ok = exampleSlice.Get(0); !ok {
		// Missing entry here
		return
	}

	fmt.Println("Value", v)
}
```

### Slice.Set
```go
func ExampleSlice_Set() {
	var err error
	if err = exampleSlice.Set(0, 1337); err != nil {
		// Handle error here
		return
	}
}
```

### Slice.Append
```go
func ExampleSlice_Append() {
	var err error
	if err = exampleSlice.Append(1337); err != nil {
		// Handle error here
		return
	}
}
```

### Slice.InsertAt
```go
func ExampleSlice_InsertAt() {
	var err error
	if err = exampleSlice.InsertAt(0, 1337); err != nil {
		// Handle error here
		return
	}
}
```

### Slice.RemoveAt
```go
func ExampleSlice_RemoveAt() {
	var err error
	if err = exampleSlice.RemoveAt(0); err != nil {
		// Handle error here
		return
	}
}
```

### Slice.ForEach
```go
func ExampleSlice_ForEach() {
	exampleSlice.ForEach(func(v int) (end bool) {
		fmt.Println("Value", v)
		return
	})
}
```

### Slice.Cursor
```go

func ExampleSlice_Cursor() {
	cur := exampleSlice.Cursor()
	v, ok := cur.Seek(1337)
	if !ok {
		fmt.Println("index is missing")
		return
	}

	fmt.Println("My seek value!", v)

	for ok {
		v, ok = cur.Next()
		fmt.Println("My next value!", v)
	}
}
```

### Slice.Cursor (using previous)
```go
func ExampleSlice_Cursor_prev() {
	cur := exampleSlice.Cursor()
	v, ok := cur.Seek(1337)
	if !ok {
		fmt.Println("index is missing")
		return
	}

	fmt.Println("My seek value!", v)

	for ok {
		v, ok = cur.Prev()
		fmt.Println("My previous value!", v)
	}
}
```

### Slice.Len
```go
func ExampleSlice_Len() {
	fmt.Println("Length", exampleSlice.Len())
}
```

### Slice.Slice
```go
func ExampleSlice_Slice() {
	fmt.Println("Slice copy", exampleSlice.Slice())
}
```