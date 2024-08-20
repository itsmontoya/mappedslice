package mappedslice

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

var exampleSlice *Slice[int]

func TestNew(t *testing.T) {
	type args struct {
		filepath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "basic",
			args: args{
				filepath: "test.bat",
			},
			wantErr: false,
		},
		{
			name: "no filepath",
			args: args{
				filepath: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ref, err := New[int](tt.args.filepath, 32)
			if err == nil {
				defer os.Remove(tt.args.filepath)
				defer ref.Close()
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestSlice_Get(t *testing.T) {
	type args struct {
		index int
	}

	tests := []struct {
		name            string
		numberOfEntries int
		args            args

		want    int
		wantErr bool
	}{
		{
			name:            "basic",
			numberOfEntries: 3,
			args: args{
				index: 2,
			},
			want:    2,
			wantErr: false,
		},
		{
			name:            "large set",
			numberOfEntries: 128,
			args: args{
				index: 127,
			},
			want:    127,
			wantErr: false,
		},
		{
			name:            "negative index",
			numberOfEntries: 3,
			args: args{
				index: -1,
			},
			want:    0,
			wantErr: true,
		},
		{
			name:            "out of bounds index",
			numberOfEntries: 3,
			args: args{
				index: 5,
			},
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := getTestSlice(tt.numberOfEntries)
			if err != nil {
				t.Errorf("Slice.Get(): error preparing: %v", err)
				return
			}
			defer os.Remove(m.f.Name())

			got, err := m.Get(tt.args.index)
			if (err != nil) != tt.wantErr {
				t.Errorf("Slice.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("Slice.Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlice_Set(t *testing.T) {
	type args struct {
		index int
		value int
	}

	tests := []struct {
		name            string
		numberOfEntries int
		args            args

		want    []int
		wantErr bool
	}{
		{
			name:            "last index",
			numberOfEntries: 3,
			args: args{
				index: 2,
				value: 7,
			},
			want:    []int{0, 1, 7},
			wantErr: false,
		},
		{
			name:            "middle index",
			numberOfEntries: 3,
			args: args{
				index: 1,
				value: 7,
			},
			want:    []int{0, 7, 2},
			wantErr: false,
		},
		{
			name:            "first index",
			numberOfEntries: 3,
			args: args{
				index: 0,
				value: 7,
			},
			want:    []int{7, 1, 2},
			wantErr: false,
		},
		{
			name:            "negative index",
			numberOfEntries: 3,
			args: args{
				index: -1,
				value: 7,
			},
			want:    []int{0, 1, 2},
			wantErr: true,
		},
		{
			name:            "out of bounds index",
			numberOfEntries: 3,
			args: args{
				index: 5,
				value: 7,
			},
			want:    []int{0, 1, 2},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := getTestSlice(tt.numberOfEntries)
			if err != nil {
				t.Errorf("Slice.Set(): error preparing: %v", err)
				return
			}
			defer os.Remove(m.f.Name())

			err = m.Set(tt.args.index, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Slice.Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			got := m.Slice()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Slice.Set() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlice_InsertAt(t *testing.T) {
	type args struct {
		index int
		value int
	}

	tests := []struct {
		name            string
		numberOfEntries int
		args            args

		want    []int
		wantErr bool
	}{
		{
			name:            "last index",
			numberOfEntries: 3,
			args: args{
				index: 2,
				value: 7,
			},
			want:    []int{0, 1, 7, 2},
			wantErr: false,
		},
		{
			name:            "middle index",
			numberOfEntries: 3,
			args: args{
				index: 1,
				value: 7,
			},
			want:    []int{0, 7, 1, 2},
			wantErr: false,
		},
		{
			name:            "first index",
			numberOfEntries: 3,
			args: args{
				index: 0,
				value: 7,
			},
			want:    []int{7, 0, 1, 2},
			wantErr: false,
		},
		{
			name:            "negative index",
			numberOfEntries: 3,
			args: args{
				index: -1,
			},
			want:    []int{0, 1, 2},
			wantErr: true,
		},
		{
			name:            "out of bounds index",
			numberOfEntries: 3,
			args: args{
				index: 5,
			},
			want:    []int{0, 1, 2},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := getTestSlice(tt.numberOfEntries)
			if err != nil {
				t.Errorf("Slice.InsertAt(): error preparing: %v", err)
				return
			}
			defer os.Remove(m.f.Name())

			err = m.InsertAt(tt.args.index, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Slice.InsertAt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			got := m.Slice()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Slice.InsertAt() got = %v, want %v", got, tt.want)
				return
			}
		})
	}
}

func TestSlice_RemoveAt(t *testing.T) {
	type args struct {
		index int
	}

	tests := []struct {
		name            string
		numberOfEntries int
		args            args

		want    []int
		wantErr bool
	}{
		{
			name:            "last index",
			numberOfEntries: 3,
			args: args{
				index: 2,
			},
			want:    []int{0, 1},
			wantErr: false,
		},
		{
			name:            "middle index",
			numberOfEntries: 3,
			args: args{
				index: 1,
			},
			want:    []int{0, 2},
			wantErr: false,
		},
		{
			name:            "first index",
			numberOfEntries: 3,
			args: args{
				index: 0,
			},
			want:    []int{1, 2},
			wantErr: false,
		},
		{
			name:            "negative index",
			numberOfEntries: 3,
			args: args{
				index: -1,
			},
			want:    []int{0, 1, 2},
			wantErr: true,
		},
		{
			name:            "out of bounds index",
			numberOfEntries: 3,
			args: args{
				index: 5,
			},
			want:    []int{0, 1, 2},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := getTestSlice(tt.numberOfEntries)
			if err != nil {
				t.Errorf("Slice.RemoveAt(): error preparing: %v", err)
				return
			}
			defer os.Remove(m.f.Name())

			err = m.RemoveAt(tt.args.index)
			if (err != nil) != tt.wantErr {
				t.Errorf("Slice.RemoveAt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			got := m.Slice()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Slice.RemoveAt() got = %v, want %v", got, tt.want)
				return
			}
		})
	}
}

func TestSlice_ForEach(t *testing.T) {
	type args struct {
		end bool
	}

	tests := []struct {
		name            string
		numberOfEntries int
		args            args

		want      []int
		wantEnded bool
	}{
		{
			name:            "basic",
			numberOfEntries: 3,
			args: args{
				end: false,
			},
			want:      []int{0, 1, 2},
			wantEnded: false,
		},
		{
			name:            "with end",
			numberOfEntries: 3,
			args: args{
				end: true,
			},
			want:      []int{0},
			wantEnded: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := getTestSlice(tt.numberOfEntries)
			if err != nil {
				t.Errorf("Slice.ForEach(): error preparing: %v", err)
				return
			}
			defer os.Remove(m.f.Name())

			var got []int
			gotEnded := m.ForEach(func(v int) (end bool) {
				got = append(got, v)
				return tt.args.end
			})

			if gotEnded != tt.wantEnded {
				t.Errorf("Slice.ForEach() gotEnded = %v, wantEnded %v", gotEnded, tt.wantEnded)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Slice.ForEach() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlice_Cursor(t *testing.T) {
	type args struct {
		seek int
		err  error
	}

	tests := []struct {
		name            string
		numberOfEntries int
		args            args

		want         []int
		wantNonExist bool
	}{
		{
			name:            "basic",
			numberOfEntries: 3,
			args: args{
				seek: 0,
				err:  nil,
			},
			want: []int{0, 1, 2},
		},
		{
			name:            "with seek",
			numberOfEntries: 3,
			args: args{
				seek: 1,
				err:  nil,
			},
			want: []int{1, 2},
		},
		{
			name:            "with end seek",
			numberOfEntries: 3,
			args: args{
				seek: 2,
				err:  nil,
			},
			want: []int{2},
		},
		{
			name:            "with out of bounds seek",
			numberOfEntries: 3,
			args: args{
				seek: 3,
				err:  nil,
			},
			want:         nil,
			wantNonExist: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := getTestSlice(tt.numberOfEntries)
			if err != nil {
				t.Errorf("Slice.Cursor(): error preparing: %v", err)
				return
			}
			defer os.Remove(m.f.Name())

			var got []int
			cur := m.Cursor()
			v, ok := cur.Seek(tt.args.seek)
			if !ok && !tt.wantNonExist {
				t.Errorf("Slice.Cursor(): error seeking: %v", tt.args.seek)
				return
			}

			if ok {
				got = append(got, v)
				for {
					v, ok := cur.Next()
					if !ok {
						break
					}

					got = append(got, v)
				}
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Slice.Cursor() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlice_Cursor_Prev(t *testing.T) {
	type args struct {
		seek int
		err  error
	}

	tests := []struct {
		name            string
		numberOfEntries int
		args            args

		want         []int
		wantNonExist bool
	}{
		{
			name:            "basic",
			numberOfEntries: 3,
			args: args{
				seek: 0,
				err:  nil,
			},
			want:         []int{0},
			wantNonExist: false,
		},
		{
			name:            "with seek",
			numberOfEntries: 3,
			args: args{
				seek: 1,
				err:  nil,
			},
			want:         []int{1, 0},
			wantNonExist: false,
		},
		{
			name:            "with end seek",
			numberOfEntries: 3,
			args: args{
				seek: 2,
				err:  nil,
			},
			want:         []int{2, 1, 0},
			wantNonExist: false,
		},
		{
			name:            "with out of bounds seek",
			numberOfEntries: 3,
			args: args{
				seek: 3,
				err:  nil,
			},
			want:         nil,
			wantNonExist: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := getTestSlice(tt.numberOfEntries)
			if err != nil {
				t.Errorf("Slice.Cursor(): error preparing: %v", err)
				return
			}
			defer os.Remove(m.f.Name())

			var got []int
			cur := m.Cursor()
			v, ok := cur.Seek(tt.args.seek)
			if !ok && !tt.wantNonExist {
				t.Errorf("Slice.Cursor(): error seeking: %v", tt.args.seek)
				return
			}

			if ok {
				got = append(got, v)
				for {
					v, ok := cur.Prev()
					if !ok {
						break
					}

					got = append(got, v)
				}
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Slice.Cursor() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlice_Len(t *testing.T) {
	tests := []struct {
		name            string
		numberOfEntries int

		want int
	}{
		{
			name:            "basic",
			numberOfEntries: 3,
			want:            3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := getTestSlice(tt.numberOfEntries)
			if err != nil {
				t.Errorf("Slice.Len(): error preparing: %v", err)
				return
			}
			defer os.Remove(m.f.Name())

			got := m.Len()
			if got != tt.want {
				t.Errorf("Slice.Len() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func getTestSlice(count int) (m *Slice[int], err error) {
	if m, err = New[int]("test.bat", 32); err != nil {
		return
	}

	for i := 0; i < count; i++ {
		if err = m.Append(i); err != nil {
			return
		}
	}

	return
}

func ExampleNew() {
	var err error
	if exampleSlice, err = New[int]("myfile.bat", 32); err != nil {
		// Handle error here
		return
	}
}

func ExampleSlice_Get() {
	var (
		v   int
		err error
	)

	if v, err = exampleSlice.Get(0); err != nil {
		// Handle error here
		return
	}

	fmt.Println("Value", v)
}

func ExampleSlice_Set() {
	var err error
	if err = exampleSlice.Set(0, 1337); err != nil {
		// Handle error here
		return
	}
}

func ExampleSlice_Append() {
	var err error
	if err = exampleSlice.Append(1337); err != nil {
		// Handle error here
		return
	}
}

func ExampleSlice_InsertAt() {
	var err error
	if err = exampleSlice.InsertAt(0, 1337); err != nil {
		// Handle error here
		return
	}
}

func ExampleSlice_RemoveAt() {
	var err error
	if err = exampleSlice.RemoveAt(0); err != nil {
		// Handle error here
		return
	}
}

func ExampleSlice_ForEach() {
	exampleSlice.ForEach(func(v int) (end bool) {
		fmt.Println("Value", v)
		return
	})
}

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

func ExampleSlice_Len() {
	fmt.Println("Length", exampleSlice.Len())
}

func ExampleSlice_Slice() {
	fmt.Println("Slice copy", exampleSlice.Slice())
}
