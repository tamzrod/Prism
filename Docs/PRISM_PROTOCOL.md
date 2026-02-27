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

## 3. Supported `payload_format` Values

Payload values may be supplied in JSON using a field named `payload_<short>`
where `<short>` corresponds to the base format name shown below.  Exactly one
`payload_*` field must appear; other combinations produce
`error_code 2` (Payload exclusivity).

Formats:

| Format        | Description                         | JSON field type                  |
|---------------|-------------------------------------|----------------------------------|
| `bytes`       | Raw octets                          | `payload_bytes` (base64 binary)  |
| `hex`*        | Hex string for `bytes` payload     | `payload_hex` (string)           |
| `u8`          | Unsigned 8‑bit integers             | `payload_u8` (array of numbers)  |
| `i16be`/`i16le` | Signed 16‑bit ints (big/little) | `payload_i16`                   |
| `u16be`/`u16le` | Unsigned 16‑bit ints           | `payload_u16`                   |
| `i32be`/`i32le` | Signed 32‑bit ints              | `payload_i32`                   |
| `u32be`/`u32le` | Unsigned 32‑bit ints            | `payload_u32`                   |
| `i64be`/`i64le` | Signed 64‑bit ints              | `payload_i64`                   |
| `u64be`/`u64le` | Unsigned 64‑bit ints            | `payload_u64`                   |
| `f32be`/`f32le` | 32‑bit IEEE floats             | `payload_f32`                   |
| `f64be`/`f64le` | 64‑bit IEEE floats             | `payload_f64`                   |
| `unix32be`/`unix32le` | 32‑bit unix timestamp   | any of the above (see note)     |
| `unix64be`/`unix64le` | 64‑bit unix timestamp   | any of the above (see note)     |
| `str_utf8`    | UTF‑8 string                       | `payload_str` (string)          |




> **Note:** `unixXX` formats behave identically to the corresponding
> unsigned integer format.  They are accepted to emphasise the familiar
> semantic but no special timezone handling is performed; `unix` input is
> simply treated as an unsigned number during transformation.  Input may be
> provided as the numeric array or as raw bytes (or hex string when
> `bytes` format is selected).

Missing or unknown format values cause `error_code 3`.

Missing endian notation (e.g. `i16` without `be`/`le`) is also an error.

No implicit defaults.  No word swapping.  No byte‑order guessing.

------------------------------------------------------------------------

## 4. Supported `target_format` Values

Output formats follow the same naming conventions as payloads, with
`numeric` formats producing arrays and the three string formats returning a
scalar string.  The returned JSON field is always `value_<target_format>`.

| Format        | Output example                        |
|---------------|----------------------------------------|
| Numeric       | `[ ... ]` (array)                     |
| `u8`          | `"value_u8": [1,2,3]`               |
| `i16be`/...   | `"value_i16le": [...]`              |
| `u16be`/...   |                                       |
| `i32be`/...   |                                       |
| `u32be`/...   |                                       |
| `i64be`/...   |                                       |
| `u64be`/...   |                                       |
| `f32be`/...   |                                       |
| `f64be`/...   |                                       |
| `unix32be`/...| identical to `u32*` numeric outputs   |
| `unix64be`/...| identical to `u64*` numeric outputs   |

String outputs:
- `hex`       → raw bytes as lower‑case hex string
- `base64`    → raw bytes as base64
- `rfc3339`   → standard RFC‑3339 timestamp (only when input is
                interpreted as `time.Time`)

Examples (success):

```json
{ "value_f32be": [123.45] }
{ "value_u16be": [17142, 58982] }
{ "value_hex": "deadbeef" }
{ "value_base64": "SGVsbG8=" }
```

```text
# string output is always under key matching format
```
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