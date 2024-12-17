package main

import (
	"reflect"
	"testing"
)

func Test_cp(t *testing.T) {
	tests := []struct {
		name string
		path []point
		p    point
		want []point
	}{
		{"a", []point{{0, 0}}, point{1, 1}, []point{{0, 0}, {1, 1}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cp(tt.path, tt.p); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cp() = %v, want %v", got, tt.want)
			}
		})
	}
}
