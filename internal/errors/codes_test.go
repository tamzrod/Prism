// internal/errors/codes_test.go

package errors

import "testing"

func TestErrorCodes(t *testing.T) {
    if CodeInvalidRequest != 1 || CodeUnsupported != 6 {
        t.Errorf("error codes constants not set correctly")
    }
}