// pkg/prism/prism_test.go

package prism

import "testing"

func TestVersion(t *testing.T) {
    if Version() == "" {
        t.Error("Version should not be empty")
    }
}