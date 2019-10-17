package tracer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewIntersection(t *testing.T) {
	type args struct {
		o Object
		t float64
	}
	tests := []struct {
		name string
		args args
		want Intersection
	}{
		{
			name: "test1",
			args: args{
				o: NewUnitSphere(),
				t: 3.0,
			},
			want: Intersection{
				o: NewUnitSphere(),
				t: 3.0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewIntersection(tt.args.o, tt.args.t), "should be equal")
		})
	}
}

func TestIntersections_Hit(t *testing.T) {
	tests := []struct {
		name    string
		i       Intersections
		want    Intersection
		wantErr bool
	}{
		{
			name: "t all positive",
			i: NewIntersections(
				NewIntersection(NewUnitSphere(), 1),
				NewIntersection(NewUnitSphere(), 2),
			),
			want:    NewIntersection(NewUnitSphere(), 1),
			wantErr: false,
		},
		{
			name: "t some negative",
			i: NewIntersections(
				NewIntersection(NewUnitSphere(), -1),
				NewIntersection(NewUnitSphere(), 1),
			),
			want:    NewIntersection(NewUnitSphere(), 1),
			wantErr: false,
		},
		{
			name: "t negative",
			i: NewIntersections(
				NewIntersection(NewUnitSphere(), -2),
				NewIntersection(NewUnitSphere(), -1),
			),
			wantErr: true,
		},
		{
			name: "t some negative 2",
			i: NewIntersections(
				NewIntersection(NewUnitSphere(), 5),
				NewIntersection(NewUnitSphere(), 7),
				NewIntersection(NewUnitSphere(), -3),
				NewIntersection(NewUnitSphere(), 2),
			),
			want:    NewIntersection(NewUnitSphere(), 2),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.i.Hit()

			if tt.wantErr {
				assert.Error(t, err, "should error")
			} else {
				assert.Equal(t, tt.want, got, "should equal")
			}
		})
	}
}
