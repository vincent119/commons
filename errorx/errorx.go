package errorx

import (
	"errors"
	"fmt"
)

// Wrap 包裝錯誤並加上訊息。
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", msg, err)
}

// Is 判斷錯誤鏈是否包含 target。
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As 嘗試將錯誤轉型成目標類型。
func As[T any](err error, target *T) bool {
	return errors.As(err, target)
}

// Cause 取出最底層錯誤。
func Cause(err error) error {
	for {
		unwrapped := errors.Unwrap(err)
		if unwrapped == nil {
			return err
		}
		err = unwrapped
	}
}
