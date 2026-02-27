# PRISM_PROTOCOL.md

Deterministic Binary Transformation Protocol --- V1 (LOCKED)

------------------------------------------------------------------------

## 1. Purpose

PRISM performs deterministic representation transformation only.

(payload + payload_format + target_format) → exactly ONE output value

PRISM does NOT:
- Apply semantics
- Apply scaling or units
- Decode protocol framing
- Perform word swapping
- Perform bit-level packing
- Persist or cache data

V1 is strictly byte-aligned.

------------------------------------------------------------------------

## 2. Request Requirements

A valid request MUST contain:

- payload_format
- target_format
- Exactly ONE payload_* field

Zero or multiple payload_* fields → error_code 2.

------------------------------------------------------------------------

## 3. Supported payload_format Values

Raw:
- bytes
- u8

Multi-byte (explicit endian required):
- i16be / i16le
- u16be / u16le
- i32be / i32le
- u32be / u32le
- i64be / i64le
- u64be / u64le
- f32be / f32le
- f64be / f64le
- unix32be / unix32le
- unix64be / unix64le

String:
- str_utf8

Missing endian for multi-byte types → error_code 3.

No implicit defaults.
No implicit ABCD.
No word swapping.

------------------------------------------------------------------------

## 4. Supported target_format Values

Numeric:
- u8
- i16be / i16le
- u16be / u16le
- i32be / i32le
- u32be / u32le
- i64be / i64le
- u64be / u64le
- f32be / f32le
- f64be / f64le
- unix32be / unix32le
- unix64be / unix64le

Numeric outputs are ALWAYS returned as arrays.
Array length MUST be ≥ 1.

String:
- hex
- base64
- rfc3339

String outputs return:

{ "value_string": "<string>" }

------------------------------------------------------------------------

## 5. Response Model

Success (numeric):

{ "value_<target_format>": [ ... ] }

Success (string):

{ "value_string": "<string>" }

Error:

{ "error_code": <integer> }

Exactly one field MUST exist in the response.

------------------------------------------------------------------------

## 6. Error Codes (Frozen)

1 → Invalid request  
2 → Payload exclusivity violation  
3 → Unknown or invalid format  
4 → Invalid payload  
5 → Length mismatch  
6 → Unsupported  

------------------------------------------------------------------------

## 7. Determinism Guarantees

- No implicit endian
- No fallback parsing
- No silent truncation
- No format guessing
- No word swapping

Ambiguity → error_code 3.

------------------------------------------------------------------------

End of PRISM_PROTOCOL.md V1