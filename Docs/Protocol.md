# PRISM_PROTOCOL.md

Deterministic Binary Transformation Protocol -- V1

------------------------------------------------------------------------

## 1. Scope

PRISM is a deterministic binary transformation engine.

PRISM performs representation transformation only:

(payload + payload_format + target_format) → exactly ONE output value

PRISM does NOT:

-   Apply semantics (no voltage/current/status meaning)
-   Apply scaling or units
-   Perform alarm logic
-   Decode protocol framing (Modbus/DNP/etc.)
-   Perform multi-field decoding per request
-   Perform bit-level packing
-   Store or persist data

V1 is strictly byte-aligned only.

------------------------------------------------------------------------

## 2. Request Model

A valid request MUST contain:

-   payload_format
-   target_format
-   Exactly ONE payload\_\* field

If none or more than one payload field is present → error.

------------------------------------------------------------------------

## 3. Payload Encodings (Mutually Exclusive)

Exactly one of the following fields must be present.

### 3.1 payload_u8

payload_u8: \[0, 255, ...\]

-   Array of unsigned 8-bit integers
-   Each value MUST be 0--255
-   Normalized internally to \[\]byte

### 3.2 payload_u16

payload_u16: \[0, 65535, ...\]

-   Array of unsigned 16-bit integers
-   Each value MUST be 0--65535
-   Interpretation depends on payload_format
-   Converted to bytes using explicit endian rules

### 3.3 payload_hex

payload_hex: "42F6E666"

-   Hex string
-   Must contain even number of characters
-   Upper/lowercase allowed
-   No implicit whitespace trimming (whitespace → error)
-   Decoded directly to bytes

### 3.4 payload_base64

payload_base64: "QvbmZg=="

-   Standard Base64 encoding
-   Must decode successfully
-   No URL-safe variant in V1

### 3.5 payload_str

payload_str: "hello"

-   UTF-8 string
-   Converted to bytes via UTF-8 encoding
-   Only valid when payload_format supports string interpretation

------------------------------------------------------------------------

## 4. Payload Exclusivity Rule

If:

-   zero payload fields provided → ERR_PAYLOAD_EXCLUSIVE
-   more than one payload field provided → ERR_PAYLOAD_EXCLUSIVE

There is NO precedence resolution.

Fail fast.

------------------------------------------------------------------------

## 5. payload_format (Source Interpretation)

Defines how PRISM interprets the provided payload bytes.

### 5.1 Raw

  Code       Description
  ---------- -----------------------------------------
  bytes      Raw byte stream
  u8         u8 array (normalized to bytes)
  u16be      u16 array → bytes (big endian words)
  u16le      u16 array → bytes (little endian words)
  str_utf8   UTF-8 string → bytes

### 5.2 Integer Decoding From Bytes

  Code    Bytes   Description
  ------- ------- -------------
  i16be   2       int16 BE
  i16le   2       int16 LE
  u16be   2       uint16 BE
  u16le   2       uint16 LE
  i32be   4       int32 BE
  i32le   4       int32 LE
  u32be   4       uint32 BE
  u32le   4       uint32 LE
  i64be   8       int64 BE
  i64le   8       int64 LE
  u64be   8       uint64 BE
  u64le   8       uint64 LE

### 5.3 Floating Point

  Code    Bytes   Description
  ------- ------- -------------
  f32be   4       float32 BE
  f32le   4       float32 LE
  f64be   8       float64 BE
  f64le   8       float64 LE

Word swapping is NOT supported in V1.

### 5.4 Timestamps

  Code       Bytes   Description
  ---------- ------- ----------------------------
  unix32be   4       uint32 seconds since epoch
  unix32le   4       uint32 seconds since epoch
  unix64be   8       uint64 seconds since epoch
  unix64le   8       uint64 seconds since epoch

No implicit formatting occurs during decode.

------------------------------------------------------------------------

## 6. target_format (Output Representation)

### 6.1 Scalar Numeric

  Code     Output Type
  -------- -------------
  i16      int16
  u16      uint16
  i32      int32
  u32      uint32
  i64      int64
  u64      uint64
  f32      float32
  f64      float64
  unix32   uint32
  unix64   uint64

### 6.2 Arrays

  Code           Description
  -------------- --------------------------
  u8_array       \[\]uint8
  u16_array_be   \[\]uint16 (BE grouping)
  u16_array_le   \[\]uint16 (LE grouping)

### 6.3 String Outputs

  -----------------------------------------------------------------------
  Code                    Description
  ----------------------- -----------------------------------------------
  hex                     hex string

  base64                  base64 string

  rfc3339                 RFC3339 timestamp string (only valid when
                          source is unix32/unix64)
  -----------------------------------------------------------------------

------------------------------------------------------------------------

## 7. Length Validation Rules

If byte length does not match expected size for transformation:

→ ERR_LENGTH_MISMATCH

Examples:

-   f32 requires exactly 4 bytes
-   f64 requires exactly 8 bytes
-   u16_array_be requires byte length divisible by 2

No partial decoding allowed.

------------------------------------------------------------------------

## 8. Response Model

### Success Response

Returns exactly ONE value.

Examples:

123.45

\[17142,58982\]

"2026-01-10T12:30:00Z"

### Error Response

Structured format:

{ "error": { "code": "ERR_INVALID_PAYLOAD", "message": "payload_hex
contains non-hex characters" } }

Optional:

"details": { ... }

------------------------------------------------------------------------

## 9. Error Codes (V1)

  Code                    Meaning
  ----------------------- -----------------------------------------
  ERR_INVALID_REQUEST     Missing required fields
  ERR_PAYLOAD_EXCLUSIVE   None or multiple payload\_\* provided
  ERR_UNKNOWN_FORMAT      Unknown payload_format or target_format
  ERR_INVALID_PAYLOAD     Malformed hex/base64/out-of-range array
  ERR_LENGTH_MISMATCH     Byte length invalid for transform
  ERR_UNSUPPORTED         Recognized but not supported in V1

Errors are deterministic and stable.

------------------------------------------------------------------------

## 10. Determinism Guarantees

PRISM guarantees:

-   No implicit endian assumptions
-   No implicit formatting
-   No fallback behavior
-   No silent truncation
-   No partial decoding
-   No precedence rules

Every invalid condition results in an explicit error.

------------------------------------------------------------------------

## 11. V1 Final Constraints

-   Byte-aligned only
-   One transformation per request
-   Exactly one output value
-   No semantic interpretation
-   No multi-field decoding
-   No bit-level packing

------------------------------------------------------------------------

End of PRISM_PROTOCOL.md V1
