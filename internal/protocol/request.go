// internal/protocol/request.go

package protocol

// Request represents the PRISM protocol request shape defined in Docs/PRISM_PROTOCOL.md
// Fields use pointer types where needed to allow exclusivity checks.
type Request struct {
    PayloadFormat string   `json:"payload_format"`
    TargetFormat  string   `json:"target_format"`

    // payload fields. Only one must be non-nil according to the protocol.
    PayloadBytes *[]byte   `json:"payload_bytes,omitempty"`            // raw bytes
    PayloadHex   *string   `json:"payload_hex,omitempty"`              // hex string for bytes

    // numeric arrays (all use 64-bit containers for simplicity)
    PayloadU8    *[]uint64 `json:"payload_u8,omitempty"`
    PayloadI16   *[]int64  `json:"payload_i16,omitempty"`
    PayloadU16   *[]uint64 `json:"payload_u16,omitempty"`
    PayloadI32   *[]int64  `json:"payload_i32,omitempty"`
    PayloadU32   *[]uint64 `json:"payload_u32,omitempty"`
    PayloadI64   *[]int64  `json:"payload_i64,omitempty"`
    PayloadU64   *[]uint64 `json:"payload_u64,omitempty"`
    PayloadF32   *[]float64`json:"payload_f32,omitempty"`
    PayloadF64   *[]float64`json:"payload_f64,omitempty"`
    PayloadStr   *string   `json:"payload_str,omitempty"`             // str_utf8
}
