// internal/engine/codec.go

package engine

import (
    "encoding/binary"
    "math"
    "strings"

    "github.com/tamzrod/prism/internal/errors"
)

// Decode transforms a byte slice from the given payload format into
// a normalized internal value. The returned interface{} may be one of
// []byte, []uint64, []int64, []float64, or string depending on the format.
func Decode(format string, data []byte) (interface{}, error) {
    // normalize Unix aliases to their numeric counterparts preserving endian
    if strings.HasPrefix(format, "unix32") {
        format = "u32" + format[len("unix32"):]
    }
    if strings.HasPrefix(format, "unix64") {
        format = "u64" + format[len("unix64"):]
    }

    switch format {
    case "bytes":
        return data, nil
    case "u8":
        out := make([]uint64, len(data))
        for i, b := range data {
            out[i] = uint64(b)
        }
        return out, nil
    case "i16be", "i16le":
        if len(data)%2 != 0 {
            return nil, ErrorCode(errors.CodeLengthMismatch)
        }
        cnt := len(data) / 2
        out := make([]int64, cnt)
        for i := 0; i < cnt; i++ {
            chunk := data[i*2 : i*2+2]
            var v int16
            if strings.HasSuffix(format, "be") {
                v = int16(binary.BigEndian.Uint16(chunk))
            } else {
                v = int16(binary.LittleEndian.Uint16(chunk))
            }
            out[i] = int64(v)
        }
        return out, nil
    case "u16be", "u16le":
        if len(data)%2 != 0 {
            return nil, ErrorCode(errors.CodeLengthMismatch)
        }
        cnt := len(data) / 2
        out := make([]uint64, cnt)
        for i := 0; i < cnt; i++ {
            chunk := data[i*2 : i*2+2]
            if strings.HasSuffix(format, "be") {
                out[i] = uint64(binary.BigEndian.Uint16(chunk))
            } else {
                out[i] = uint64(binary.LittleEndian.Uint16(chunk))
            }
        }
        return out, nil
    case "i32be", "i32le":
        if len(data)%4 != 0 {
            return nil, ErrorCode(errors.CodeLengthMismatch)
        }
        cnt := len(data) / 4
        out := make([]int64, cnt)
        for i := 0; i < cnt; i++ {
            chunk := data[i*4 : i*4+4]
            var v int32
            if strings.HasSuffix(format, "be") {
                v = int32(binary.BigEndian.Uint32(chunk))
            } else {
                v = int32(binary.LittleEndian.Uint32(chunk))
            }
            out[i] = int64(v)
        }
        return out, nil
    case "u32be", "u32le":
        if len(data)%4 != 0 {
            return nil, ErrorCode(errors.CodeLengthMismatch)
        }
        cnt := len(data) / 4
        out := make([]uint64, cnt)
        for i := 0; i < cnt; i++ {
            chunk := data[i*4 : i*4+4]
            if strings.HasSuffix(format, "be") {
                out[i] = uint64(binary.BigEndian.Uint32(chunk))
            } else {
                out[i] = uint64(binary.LittleEndian.Uint32(chunk))
            }
        }
        return out, nil
    case "i64be", "i64le":
        if len(data)%8 != 0 {
            return nil, ErrorCode(errors.CodeLengthMismatch)
        }
        cnt := len(data) / 8
        out := make([]int64, cnt)
        for i := 0; i < cnt; i++ {
            chunk := data[i*8 : i*8+8]
            var v int64
            if strings.HasSuffix(format, "be") {
                v = int64(binary.BigEndian.Uint64(chunk))
            } else {
                v = int64(binary.LittleEndian.Uint64(chunk))
            }
            out[i] = v
        }
        return out, nil
    case "u64be", "u64le":
        if len(data)%8 != 0 {
            return nil, ErrorCode(errors.CodeLengthMismatch)
        }
        cnt := len(data) / 8
        out := make([]uint64, cnt)
        for i := 0; i < cnt; i++ {
            chunk := data[i*8 : i*8+8]
            if strings.HasSuffix(format, "be") {
                out[i] = binary.BigEndian.Uint64(chunk)
            } else {
                out[i] = binary.LittleEndian.Uint64(chunk)
            }
        }
        return out, nil
    case "f32be", "f32le":
        if len(data)%4 != 0 {
            return nil, ErrorCode(errors.CodeLengthMismatch)
        }
        cnt := len(data) / 4
        out := make([]float64, cnt)
        for i := 0; i < cnt; i++ {
            chunk := data[i*4 : i*4+4]
            var u uint32
            if strings.HasSuffix(format, "be") {
                u = binary.BigEndian.Uint32(chunk)
            } else {
                u = binary.LittleEndian.Uint32(chunk)
            }
            out[i] = float64(math.Float32frombits(u))
        }
        return out, nil
    case "f64be", "f64le":
        if len(data)%8 != 0 {
            return nil, ErrorCode(errors.CodeLengthMismatch)
        }
        cnt := len(data) / 8
        out := make([]float64, cnt)
        for i := 0; i < cnt; i++ {
            chunk := data[i*8 : i*8+8]
            var u uint64
            if strings.HasSuffix(format, "be") {
                u = binary.BigEndian.Uint64(chunk)
            } else {
                u = binary.LittleEndian.Uint64(chunk)
            }
            out[i] = math.Float64frombits(u)
        }
        return out, nil
    case "str_utf8":
        return string(data), nil
    default:
        return nil, ErrorCode(errors.CodeUnknownFormat)
    }
}