package shell

import (
	"errors"
	"strings"
	"testing"
)

func TestErrorMessage_ColorsAndReset(t *testing.T) {
	err := errors.New("fail")
	msg := ErrorMessage(err)
	color := getColor(colorConfig.ErrorColor)
	if !strings.Contains(msg, color) {
		t.Errorf("ErrorMessage missing ErrorColor: %q", msg)
	}
	if !strings.Contains(msg, ResetColor) {
		t.Errorf("ErrorMessage missing ResetColor: %q", msg)
	}
	if !strings.HasPrefix(msg, color+"Error: ") {
		t.Errorf("ErrorMessage does not start with color and prefix: %q", msg)
	}
	if !strings.HasSuffix(msg, ResetColor+"\n") {
		t.Errorf("ErrorMessage does not end with reset and newline: %q", msg)
	}
}
