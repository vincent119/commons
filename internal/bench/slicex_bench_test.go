package bench

import (
	"github.com/vincent119/commons/slicex"
	"testing"
)

func BenchmarkContains_Sizes(b *testing.B) {
	types := []struct {
		name string
		n    int
		want int // -1=absent, >=0 present index
	}{
		{"n=10_present", 10, 5},
		{"n=10_absent", 10, -1},
		{"n=1k_present", 1000, 777},
		{"n=1k_absent", 1000, -1},
		{"n=100k_present", 100000, 54321},
		{"n=100k_absent", 100000, -1},
	}
	for _, tc := range types {
		b.Run(tc.name, func(b *testing.B) {
			// 準備資料：遞增整數 slice
			s := make([]int, tc.n)
			for i := 0; i < tc.n; i++ {
				s[i] = i
			}
			var target int
			if tc.want >= 0 {
				target = s[tc.want]
			} else {
				target = -1
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = slicex.Contains(s, target)
			}
		})
	}
}

func BenchmarkIndexOf_Sizes(b *testing.B) {
	types := []struct {
		name string
		n    int
		want int
	}{
		{"n=10_present", 10, 5},
		{"n=10_absent", 10, -1},
		{"n=1k_present", 1000, 777},
		{"n=1k_absent", 1000, -1},
		{"n=100k_present", 100000, 54321},
		{"n=100k_absent", 100000, -1},
	}
	for _, tc := range types {
		b.Run(tc.name, func(b *testing.B) {
			s := make([]int, tc.n)
			for i := 0; i < tc.n; i++ {
				s[i] = i
			}
			var target int
			if tc.want >= 0 {
				target = s[tc.want]
			} else {
				target = -1
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = slicex.IndexOf(s, target)
			}
		})
	}
}

func BenchmarkMap_Filter(b *testing.B) {
	s := make([]int, 10000)
	for i := 0; i < len(s); i++ {
		s[i] = i
	}
	b.Run("Map", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = slicex.Map(s, func(v int) int { return v * 2 })
		}
	})
	b.Run("Filter_even", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = slicex.Filter(s, func(v int) bool { return v%2 == 0 })
		}
	})
}
