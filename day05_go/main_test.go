package main

import "testing"

func Test_isSubset(t *testing.T) {
	type args struct {
		whole []int
		part  []int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"empty", args{[]int{}, []int{}}, true},
		{"empty part", args{[]int{1, 2, 3}, []int{}}, true},
		{"empty whole", args{[]int{}, []int{1, 2, 3}}, false},
		{"equal", args{[]int{1, 2, 3}, []int{1, 2, 3}}, false},
		{"subset1", args{[]int{1, 2, 3, 4, 5}, []int{2, 3, 5}}, true},
		{"subset2", args{[]int{6, 2, 3, 4, 5}, []int{6, 2, 3, 5}}, true},
		{"fail", args{[]int{6, 2, 3, 4, 5}, []int{2, 3, 5, 6}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isSubset(tt.args.whole, tt.args.part); got != tt.want {
				t.Errorf("isSubset() = %v, want %v", got, tt.want)
			}
		})
	}
}
