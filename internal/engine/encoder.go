// internal/engine/encoder.go

package engine

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/tamzrod/prism/internal/errors"
	"github.com/tamzrod/prism/internal/protocol"
)

// EncodeResult takes a normalized internal value (as produced by Decode)
// and converts it into a form suitable for inclusion in the response JSON.
// For numeric formats the returned value will be a slice of numbers. For
// string outputs the result is a Go string.
func EncodeResult(value interface{}, format string) (interface{}, error) {
	switch format {
	case "hex":
		b, ok := value.([]byte)
		if !ok {
			return nil, ErrorCode(errors.CodeInvalidPayload)
		}
		return hex.EncodeToString(b), nil
	case "base64":
		b, ok := value.([]byte)
		if !ok {
			return nil, ErrorCode(errors.CodeInvalidPayload)
		}
		return base64.StdEncoding.EncodeToString(b), nil
	case "rfc3339":
		t, ok := value.(time.Time)
		if !ok {
			return nil, ErrorCode(errors.CodeInvalidPayload)
		}
		return t.Format(time.RFC3339), nil
	default:
		// numeric formats simply expect the correct slice type and return it
		switch fmt.Sprintf("%T", value) {
		case "[]uint64", "[]int64", "[]float64":
			return value, nil
		default:
			return nil, ErrorCode(errors.CodeUnknownFormat)
		}
	}
}

// Process performs a full request -> response operation. It validates the
// request, decodes the payload, encodes the result to the target format, and
// returns a JSON‑ready byte slice representing either a success or an error
// according to the protocol.
func Process(req *protocol.Request) ([]byte, error) {
	if err := ValidateRequest(req); err != nil {
		code := int(err.(ErrorCode))
		return json.Marshal(map[string]int{"error_code": code})
	}

	internal, code := extractInternal(req)
	if code != 0 {
		return json.Marshal(map[string]int{"error_code": int(code)})
	}

	// if target is a string format operating on bytes we may already have
	// the raw data packaged above.  otherwise we attempt to convert the
	// internal value via Decode (rules repeated from earlier code).
	var out interface{}
	var err error
	if req.TargetFormat == "hex" || req.TargetFormat == "base64" {
		// For byte-oriented string targets, ensure we have raw bytes.
		switch v := internal.(type) {
		case []byte:
			out, err = EncodeResult(v, req.TargetFormat)
			if err != nil {
				if ec, ok := err.(ErrorCode); ok {
					return json.Marshal(map[string]int{"error_code": int(ec)})
				}
				return json.Marshal(map[string]int{"error_code": 6})
			}
		default:
			// attempt to convert numeric internal arrays to raw bytes
			b, ec := internalToBytes(req.PayloadFormat, v)
			if ec != 0 {
				return json.Marshal(map[string]int{"error_code": int(ec)})
			}
			out, err = EncodeResult(b, req.TargetFormat)
			if err != nil {
				if ec2, ok := err.(ErrorCode); ok {
					return json.Marshal(map[string]int{"error_code": int(ec2)})
				}
				return json.Marshal(map[string]int{"error_code": 6})
			}
		}
	} else {
		// if we have a []byte from decoding or payload we still need to run
		// it through Decode to transform to the desired numeric layout.
		switch v := internal.(type) {
		case []byte:
			out, err = Decode(req.TargetFormat, v)
		default:
			// non-byte values (numeric arrays etc) are passed straight to
			// EncodeResult which will handle them appropriately.
			out, err = EncodeResult(v, req.TargetFormat)
			// EncodeResult returns internal or error; we can skip final
			// marshalling below because we already got the formatted value.
			if err == nil {
				// we already have the proper out in this case, return below
				goto build
			}
		}
		if err != nil {
			if ec, ok := err.(ErrorCode); ok {
				return json.Marshal(map[string]int{"error_code": int(ec)})
			}
			return json.Marshal(map[string]int{"error_code": 4})
		}
	}

build:
	// build response map: always use value_<target_format> as the key
	key := "value_" + req.TargetFormat
	return json.Marshal(map[string]interface{}{key: out})
}

// internalToBytes converts a normalized internal numeric value into a raw
// byte slice according to the provided payload format (which encodes
// element width and endianness). Returns an ErrorCode on failure.
func internalToBytes(format string, v interface{}) ([]byte, ErrorCode) {
	switch format {
	case "u8":
		arr, ok := v.([]uint64)
		if !ok {
			return nil, ErrorCode(errors.CodeInvalidPayload)
		}
		b := make([]byte, len(arr))
		for i, x := range arr {
			b[i] = byte(x & 0xFF)
		}
		return b, 0
	case "i16be", "i16le":
		arr, ok := v.([]int64)
		if !ok {
			return nil, ErrorCode(errors.CodeInvalidPayload)
		}
		b := make([]byte, 2*len(arr))
		for i, x := range arr {
			ux := uint16(int16(x))
			if strings.HasSuffix(format, "be") {
				binary.BigEndian.PutUint16(b[i*2:], ux)
			} else {
				binary.LittleEndian.PutUint16(b[i*2:], ux)
			}
		}
		return b, 0
	case "u16be", "u16le":
		arr, ok := v.([]uint64)
		if !ok {
			return nil, ErrorCode(errors.CodeInvalidPayload)
		}
		b := make([]byte, 2*len(arr))
		for i, x := range arr {
			ux := uint16(x)
			if strings.HasSuffix(format, "be") {
				binary.BigEndian.PutUint16(b[i*2:], ux)
			} else {
				binary.LittleEndian.PutUint16(b[i*2:], ux)
			}
		}
		return b, 0
	case "i32be", "i32le":
		arr, ok := v.([]int64)
		if !ok {
			return nil, ErrorCode(errors.CodeInvalidPayload)
		}
		b := make([]byte, 4*len(arr))
		for i, x := range arr {
			ux := uint32(int32(x))
			if strings.HasSuffix(format, "be") {
				binary.BigEndian.PutUint32(b[i*4:], ux)
			} else {
				binary.LittleEndian.PutUint32(b[i*4:], ux)
			}
		}
		return b, 0
	case "u32be", "u32le", "unix32be", "unix32le":
		arr, ok := v.([]uint64)
		if !ok {
			return nil, ErrorCode(errors.CodeInvalidPayload)
		}
		b := make([]byte, 4*len(arr))
		for i, x := range arr {
			ux := uint32(x)
			if strings.HasSuffix(format, "be") {
				binary.BigEndian.PutUint32(b[i*4:], ux)
			} else {
				binary.LittleEndian.PutUint32(b[i*4:], ux)
			}
		}
		return b, 0
	case "i64be", "i64le":
		arr, ok := v.([]int64)
		if !ok {
			return nil, ErrorCode(errors.CodeInvalidPayload)
		}
		b := make([]byte, 8*len(arr))
		for i, x := range arr {
			ux := uint64(x)
			if strings.HasSuffix(format, "be") {
				binary.BigEndian.PutUint64(b[i*8:], ux)
			} else {
				binary.LittleEndian.PutUint64(b[i*8:], ux)
			}
		}
		return b, 0
	case "u64be", "u64le", "unix64be", "unix64le":
		arr, ok := v.([]uint64)
		if !ok {
			return nil, ErrorCode(errors.CodeInvalidPayload)
		}
		b := make([]byte, 8*len(arr))
		for i, x := range arr {
			if strings.HasSuffix(format, "be") {
				binary.BigEndian.PutUint64(b[i*8:], x)
			} else {
				binary.LittleEndian.PutUint64(b[i*8:], x)
			}
		}
		return b, 0
	case "f32be", "f32le":
		arr, ok := v.([]float64)
		if !ok {
			return nil, ErrorCode(errors.CodeInvalidPayload)
		}
		b := make([]byte, 4*len(arr))
		for i, x := range arr {
			ux := math.Float32bits(float32(x))
			if strings.HasSuffix(format, "be") {
				binary.BigEndian.PutUint32(b[i*4:], ux)
			} else {
				binary.LittleEndian.PutUint32(b[i*4:], ux)
			}
		}
		return b, 0
	case "f64be", "f64le":
		arr, ok := v.([]float64)
		if !ok {
			return nil, ErrorCode(errors.CodeInvalidPayload)
		}
		b := make([]byte, 8*len(arr))
		for i, x := range arr {
			ux := math.Float64bits(x)
			if strings.HasSuffix(format, "be") {
				binary.BigEndian.PutUint64(b[i*8:], ux)
			} else {
				binary.LittleEndian.PutUint64(b[i*8:], ux)
			}
		}
		return b, 0
	case "str_utf8":
		s, ok := v.(string)
		if !ok {
			return nil, ErrorCode(errors.CodeInvalidPayload)
		}
		return []byte(s), 0
	default:
		return nil, ErrorCode(errors.CodeUnknownFormat)
	}
}

// extractInternal examines the Request and returns a normalized internal
// representation suitable for feeding into EncodeResult/Decode. It handles
// hex string payloads and numeric fields. On error it returns a non‑zero
// ErrorCode.
func extractInternal(req *protocol.Request) (interface{}, ErrorCode) {
	switch req.PayloadFormat {
	case "bytes", "hex":
		if req.PayloadBytes != nil {
			return *req.PayloadBytes, 0
		}
		if req.PayloadHex != nil {
			b, err := hex.DecodeString(*req.PayloadHex)
			if err != nil {
				return nil, ErrorCode(errors.CodeInvalidPayload)
			}
			return b, 0
		}
		return nil, ErrorCode(errors.CodeInvalidPayload)
	case "u8":
		if req.PayloadU8 != nil {
			return *req.PayloadU8, 0
		}
		if req.PayloadBytes != nil {
			v, err := Decode("u8", *req.PayloadBytes)
			if err != nil {
				if ec, ok := err.(ErrorCode); ok {
					return nil, ec
				}
				return nil, ErrorCode(errors.CodeInvalidPayload)
			}
			return v, 0
		}
		return nil, ErrorCode(errors.CodeInvalidPayload)
	case "i16be", "i16le":
		if req.PayloadI16 != nil {
			return *req.PayloadI16, 0
		}
		if req.PayloadBytes != nil {
			v, err := Decode(req.PayloadFormat, *req.PayloadBytes)
			if err != nil {
				if ec, ok := err.(ErrorCode); ok {
					return nil, ec
				}
				return nil, ErrorCode(errors.CodeInvalidPayload)
			}
			return v, 0
		}
		return nil, ErrorCode(errors.CodeInvalidPayload)
	case "u16be", "u16le":
		if req.PayloadU16 != nil {
			return *req.PayloadU16, 0
		}
		if req.PayloadBytes != nil {
			v, err := Decode(req.PayloadFormat, *req.PayloadBytes)
			if err != nil {
				if ec, ok := err.(ErrorCode); ok {
					return nil, ec
				}
				return nil, ErrorCode(errors.CodeInvalidPayload)
			}
			return v, 0
		}
		return nil, ErrorCode(errors.CodeInvalidPayload)
	case "i32be", "i32le":
		if req.PayloadI32 != nil {
			return *req.PayloadI32, 0
		}
		if req.PayloadBytes != nil {
			v, err := Decode(req.PayloadFormat, *req.PayloadBytes)
			if err != nil {
				if ec, ok := err.(ErrorCode); ok {
					return nil, ec
				}
				return nil, ErrorCode(errors.CodeInvalidPayload)
			}
			return v, 0
		}
		return nil, ErrorCode(errors.CodeInvalidPayload)
	case "u32be", "u32le", "unix32be", "unix32le":
		if req.PayloadU32 != nil {
			return *req.PayloadU32, 0
		}
		if req.PayloadBytes != nil {
			v, err := Decode(req.PayloadFormat, *req.PayloadBytes)
			if err != nil {
				if ec, ok := err.(ErrorCode); ok {
					return nil, ec
				}
				return nil, ErrorCode(errors.CodeInvalidPayload)
			}
			return v, 0
		}
		return nil, ErrorCode(errors.CodeInvalidPayload)
	case "i64be", "i64le":
		if req.PayloadI64 != nil {
			return *req.PayloadI64, 0
		}
		if req.PayloadBytes != nil {
			v, err := Decode(req.PayloadFormat, *req.PayloadBytes)
			if err != nil {
				if ec, ok := err.(ErrorCode); ok {
					return nil, ec
				}
				return nil, ErrorCode(errors.CodeInvalidPayload)
			}
			return v, 0
		}
		return nil, ErrorCode(errors.CodeInvalidPayload)
	case "u64be", "u64le", "unix64be", "unix64le":
		if req.PayloadU64 != nil {
			return *req.PayloadU64, 0
		}
		if req.PayloadBytes != nil {
			v, err := Decode(req.PayloadFormat, *req.PayloadBytes)
			if err != nil {
				if ec, ok := err.(ErrorCode); ok {
					return nil, ec
				}
				return nil, ErrorCode(errors.CodeInvalidPayload)
			}
			return v, 0
		}
		return nil, ErrorCode(errors.CodeInvalidPayload)
	case "f32be", "f32le":
		if req.PayloadF32 != nil {
			return *req.PayloadF32, 0
		}
		if req.PayloadBytes != nil {
			v, err := Decode(req.PayloadFormat, *req.PayloadBytes)
			if err != nil {
				if ec, ok := err.(ErrorCode); ok {
					return nil, ec
				}
				return nil, ErrorCode(errors.CodeInvalidPayload)
			}
			return v, 0
		}
		return nil, ErrorCode(errors.CodeInvalidPayload)
	case "f64be", "f64le":
		if req.PayloadF64 != nil {
			return *req.PayloadF64, 0
		}
		if req.PayloadBytes != nil {
			v, err := Decode(req.PayloadFormat, *req.PayloadBytes)
			if err != nil {
				if ec, ok := err.(ErrorCode); ok {
					return nil, ec
				}
				return nil, ErrorCode(errors.CodeInvalidPayload)
			}
			return v, 0
		}
		return nil, ErrorCode(errors.CodeInvalidPayload)
	case "str_utf8":
		if req.PayloadStr != nil {
			return *req.PayloadStr, 0
		}
		if req.PayloadBytes != nil {
			return string(*req.PayloadBytes), 0
		}
		return nil, ErrorCode(errors.CodeInvalidPayload)
	default:
		return nil, ErrorCode(errors.CodeUnknownFormat)
	}
}
