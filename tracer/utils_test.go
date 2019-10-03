package tracer

import (
	"testing"

	"github.com/ToQoz/gopwt/assert"
)

func TestEquals(t *testing.T) {
	type args struct {
		a float64
		b float64
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "equals",
			args: args{
				a: 2.345685,
				b: 2.345684,
			},
			want: true,
		},
		{
			name: "not equals",
			args: args{
				a: 2.44567,
				b: 2.34568,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Equals(tt.args.a, tt.args.b)

			assert.OK(t, result == tt.want, "float comparison")
		})
	}
}
