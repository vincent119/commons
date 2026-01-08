package bench

import (
	"errors"
	"github.com/vincent119/commons/errorx"
	"io"
	"testing"
)

func BenchmarkWrap(b *testing.B) {
	root := io.EOF
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = errorx.Wrap(root, "read failed")
	}
}

func BenchmarkIs(b *testing.B) {
	root := io.EOF
	err := errorx.Wrap(root, "context")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = errorx.Is(err, io.EOF)
	}
}

func BenchmarkCause_Depth(b *testing.B) {
	mk := func(depth int) error {
		err := errors.New("root")
		for i := 0; i < depth; i++ {
			err = errorx.Wrap(err, "ctx")
		}
		return err
	}
	for _, d := range []int{1, 3, 5, 8} {
		b.Run("depth="+itoa(d), func(b *testing.B) {
			e := mk(d)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = errorx.Cause(e)
			}
		})
	}
}

// 避免引入 strconv 導致干擾，簡單寫個小工具。
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}
