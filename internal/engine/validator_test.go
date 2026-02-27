// internal/engine/validator_test.go

package engine

import (
    "testing"

    "github.com/tamzrod/prism/internal/protocol"
)

func TestValidateRequest_ReturnsNil(t *testing.T) {
    req := &protocol.Request{
        PayloadFormat: "bytes",
        TargetFormat:  "hex",
        PayloadBytes:  &[]byte{0x01, 0x02},
    }
    if err := ValidateRequest(req); err != nil {
        t.Errorf("expected nil error from valid request, got %v", err)
    }
}

func TestValidateRequest_PayloadExclusivity(t *testing.T) {
    req := &protocol.Request{PayloadFormat: "bytes", TargetFormat: "hex"}
    if err := ValidateRequest(req); err == nil {
        t.Error("expected error for missing payload")
    }

    // cannot supply more than one payload field
    req = &protocol.Request{
        PayloadFormat: "bytes",
        TargetFormat:  "hex",
        PayloadBytes:  &[]byte{0x01},
        PayloadHex:    new(string),
    }
    if err := ValidateRequest(req); err == nil {
        t.Error("expected error for multiple payloads")
    }
}

func TestValidateRequest_HexPayloadAllowed(t *testing.T) {
    s := "deadbeef"
    req := &protocol.Request{
        PayloadFormat: "bytes",
        TargetFormat:  "hex",
        PayloadHex:    &s,
    }
    if err := ValidateRequest(req); err != nil {
        t.Errorf("unexpected validation error: %v", err)
    }
}

func TestErrorCode_ImplementsError(t *testing.T) {
    var err error = ErrorCode(3)
    if err.Error() == "" {
        t.Error("ErrorCode.Error() returned empty string")
    }
}