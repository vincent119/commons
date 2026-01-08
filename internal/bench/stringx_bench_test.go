package bench

import (
	"github.com/vincent119/commons/stringx"
	"testing"
)

func BenchmarkToSnake(b *testing.B) {
	inputs := []string{"UserID", "userName", "user-name", "LongHTTPHeaderName", "  spaced  Name  "}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, in := range inputs {
			_ = stringx.ToSnake(in)
		}
	}
}
