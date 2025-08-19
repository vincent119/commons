package errorx

import (
	"errors"
	"io"
	"testing"
)

func TestWrapAndCause(t *testing.T) {
	err := errors.New("root")
	wrapped := Wrap(err, "context")
	if Cause(wrapped).Error() != "root" {
		t.Fatalf("expected root cause, got %v", Cause(wrapped))
	}
}

func TestIsAndAs(t *testing.T) {
	var targetErr = io.EOF
	err := Wrap(targetErr, "reading failed")
	if !Is(err, io.EOF) {
		t.Fatal("expected error to be EOF")
	}

	var eof error
	if !As(err, &eof) {
		t.Fatal("expected As to succeed")
	}
}