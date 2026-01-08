package bench

import (
	"github.com/vincent119/commons/timex"
	"testing"
	"time"
)

func BenchmarkStartOfDay(b *testing.B) {
	loc, _ := time.LoadLocation("Asia/Taipei")
	t := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = timex.StartOfDay(t, loc)
	}
}
