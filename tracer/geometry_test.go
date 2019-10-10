package tracer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTuple_Equals(t *testing.T) {
	type args struct {
		t Tuple
	}
	tests := []struct {
		name string
		s    Tuple
		args args
		want bool
	}{
		{
			name: "equals",
			s:    NewTuple(1, 1, 1, 0),
			args: args{
				t: NewTuple(1, 1, 1, 0),
			},
			want: true,
		},
		{
			name: "not equals",
			s:    NewTuple(1, 1, 1, 1),
			args: args{
				t: NewTuple(1, 2, 1, 1),
			},
			want: false,
		},
		{
			name: "not equals",
			s:    NewTuple(1, 1, 1, 1),
			args: args{
				t: NewTuple(1, 1, 1, 0),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.s.Equals(tt.args.t), tt.want, "should be equal")
		})
	}
}
