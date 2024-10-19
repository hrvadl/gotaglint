//go:build abc

package examples

import "testing"

func Test_sum(t *testing.T) {
	t.Parallel()
	type args struct {
		a int
		b int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "ShouldSumCorrectly",
			args: args{
				a: 1,
				b: 1,
			},
			want: 2,
		},
		{
			name: "ShouldSumCorrectly",
			args: args{
				a: 2,
				b: 2,
			},
			want: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := sum(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("sum() = %v, want %v", got, tt.want)
			}
		})
	}
}
