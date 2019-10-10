package dump

import (
	"runtime"
	"testing"
)

func Test_delta64(t *testing.T) {
	type args struct {
		prev    uint64
		current uint64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "big positive",
			args: args{
				prev:    3800000000, // 3.82G ,
				current: 3900000000, // 3.92G ,
			},
			want: 100000000,
		},

		{
			name: "big positive 2",
			args: args{
				prev:    uint64(3*giga) + uint64(800*mega), // 3.82G ,
				current: uint64(3*giga) + uint64(900*mega), // 3.92G ,
			},
			want: int64(100 * mega),
		},

		{
			name: "big negative",
			args: args{
				prev:    3800000000, // 3.82G ,
				current: 3900000000, // 3.92G ,
			},
			want: 100000000,
		},

		{
			name: "big negative 2",
			args: args{
				prev:    uint64(3*giga) + uint64(900*mega), // 3.92G ,
				current: uint64(3*giga) + uint64(800*mega), // 3.82G ,
			},
			want: -int64(100 * mega),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := delta64(tt.args.current, tt.args.prev); got != tt.want {
				t.Errorf("delta64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_printMem(t *testing.T) {
	type args struct {
		message  string
		prevM    *runtime.MemStats
		currentM *runtime.MemStats
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "big negative",
			args: args{
				message: "big negative",
				prevM: &runtime.MemStats{
					HeapAlloc: uint64(3*giga) + uint64(900*mega),
				},
				currentM: &runtime.MemStats{
					HeapAlloc: uint64(3*giga) + uint64(800*mega),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			printMem(tt.name, tt.args.prevM, tt.args.currentM)
		})
	}
}
