package mappedslice

import (
	"fmt"
	"testing"
)

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
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRef, err := New[int](tt.args.filepath)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			s := gotRef.Slice()
			s = append(s, 1337)
			s = append(s, 8)
			s = append(s, 8)
			s = append(s, 8)
			s = append(s, 8)
			s = append(s, 8)
			s = append(s, 8)
			s = append(s, 8)
			s = append(s, 8)
			s = append(s, 8)
			s = append(s, 8)
			s = append(s, 8)
			s = append(s, 8)
			s = append(s, 8)
			s = append(s, 8)
			s = append(s, 8)
			s = append(s, 8)
			s = append(s, 8)
			s = append(s, 8)
			s = append(s, 8)
			s = append(s, 8)
			s = append(s, 8)
			fmt.Println("S", s, len(s), cap(s))

			gotRef.Close()

			gotRef, err = New[int](tt.args.filepath)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			fmt.Println("S", s, len(s), cap(s))
		})
	}
}
