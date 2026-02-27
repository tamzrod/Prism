// internal/engine/validator.go

package engine

import (
    "fmt"

    "github.com/tamzrod/prism/internal/errors"
    "github.com/tamzrod/prism/internal/protocol"
)

// ErrorCode is a simple numeric error type returned by various engine
// functions. It mirrors the constants defined in the internal/errors
// package but is defined here so callers need not import the protocol
// package just to receive an error.

type ErrorCode int

func (e ErrorCode) Error() string {
    return fmt.Sprintf("error code %d", int(e))
}

// supported payload and target formats (V1)
var payloadFormats = map[string]struct{}{
    "bytes": {}, "hex": {}, "u8": {},
    "i16be": {}, "i16le": {}, "u16be": {}, "u16le": {},
    "i32be": {}, "i32le": {}, "u32be": {}, "u32le": {},
    "i64be": {}, "i64le": {}, "u64be": {}, "u64le": {},
    "f32be": {}, "f32le": {}, "f64be": {}, "f64le": {},
    "unix32be": {}, "unix32le": {}, "unix64be": {}, "unix64le": {},
    "str_utf8": {},
}

var targetFormats = map[string]struct{}{
    "u8": {},
    "i16be": {}, "i16le": {}, "u16be": {}, "u16le": {},
    "i32be": {}, "i32le": {}, "u32be": {}, "u32le": {},
    "i64be": {}, "i64le": {}, "u64be": {}, "u64le": {},
    "f32be": {}, "f32le": {}, "f64be": {}, "f64le": {},
    "unix32be": {}, "unix32le": {}, "unix64be": {}, "unix64le": {},
    "hex": {}, "base64": {}, "rfc3339": {},
}

// ValidateRequest ensures that the protocol constraints are met, such as
// exactly one payload field and known formats. It returns an ErrorCode on
// failure (non‑zero) or nil on success.
func ValidateRequest(req *protocol.Request) error {
    if req == nil {
        return ErrorCode(errors.CodeInvalidRequest)
    }
    if req.PayloadFormat == "" || req.TargetFormat == "" {
        return ErrorCode(errors.CodeInvalidRequest)
    }

    // payload exclusivity; exactly one payload_* field must be non‑nil
    count := 0
    if req.PayloadBytes != nil {
        count++
    }
    if req.PayloadHex != nil {
        count++
    }
    if req.PayloadU8 != nil {
        count++
    }
    if req.PayloadI16 != nil {
        count++
    }
    if req.PayloadU16 != nil {
        count++
    }
    if req.PayloadI32 != nil {
        count++
    }
    if req.PayloadU32 != nil {
        count++
    }
    if req.PayloadI64 != nil {
        count++
    }
    if req.PayloadU64 != nil {
        count++
    }
    if req.PayloadF32 != nil {
        count++
    }
    if req.PayloadF64 != nil {
        count++
    }
    if req.PayloadStr != nil {
        count++
    }
    if count != 1 {
        return ErrorCode(errors.CodePayloadExclusivity)
    }

    if _, ok := payloadFormats[req.PayloadFormat]; !ok {
        return ErrorCode(errors.CodeUnknownFormat)
    }
    if _, ok := targetFormats[req.TargetFormat]; !ok {
        return ErrorCode(errors.CodeUnknownFormat)
    }

    return nil
}
