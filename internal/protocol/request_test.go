// internal/protocol/request_test.go

package protocol

import "testing"

func TestRequestStruct(t *testing.T) {
    req := Request{PayloadFormat: "bytes", TargetFormat: "u8"}
    if req.PayloadFormat == "" || req.TargetFormat == "" {
        t.Fatal("fields should be set")
    }
}