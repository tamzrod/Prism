# PRISM User Manual

This document explains how to build and run the PRISM binary and includes
numerous request/response examples demonstrating the supported format
combinations.

---

## Building

```bash
# Windows example
cd "D:\2026\Go Programming\Prism"
go build -o bin/prism.exe ./cmd/prism
```

The `bin/prism.exe` binary is produced.  You may also build on Linux or Mac
the same way (remove `.exe`).

## Running the Server

By default the server listens on port 8080, which may be overridden via
flag or configuration.

```bash
./bin/prism.exe -port 9459          # listen on TCP 9459
# or
# set Port: 9459 in config file and run without flags
```

The server uses a simple TCP protocol: each request is a single JSON object
terminated by newline, and the response is a JSON object sent back on the
same connection.

Example using `nc` (netcat) for quick manual testing:

```bash
printf '{"payload_format":"u16be","target_format":"hex","payload_u16":[17142,58982]}' | nc localhost 9459
# => {"value_hex":"42f6e666"}
```

## Request Format

A valid request JSON must include:

- `payload_format` (string)
- `target_format` (string)
- exactly one `payload_*` field corresponding to the chosen payload
  format.

The naming rule is `payload_<short>` where `<short>` is the base format
(`bytes`, `u8`, `u16`, etc.).

Numeric payloads use arrays of numbers.  Raw binary may be provided using
`payload_bytes` (base64 encoded) or `payload_hex` (lower‑case hex string).
`payload_format` may also be `hex`, an alias for `bytes` when you supply
`payload_hex`.

### Payload formats (examples)

```json
{ "payload_format":"bytes", "payload_bytes":"AQID" }            
# raw 0x01 0x02 0x03

{ "payload_format":"hex", "payload_hex":"01ff" }

{ "payload_format":"u8",   "payload_u8":[1,255] }

{ "payload_format":"i16be","payload_i16":[258] }

{ "payload_format":"unix32le","payload_u32":[1581231230] }
```

Unix- prefixed formats (`unix32be`, etc.) behave identically to the numeric
formats.  They are provided for semantic clarity; no timezone or range
checking is done.

### Target formats (examples)

- Numeric outputs: `u8`, `i16be`, `u32le`, `f64be`, `unix64le`, etc.
- String outputs: `hex`, `base64`, `rfc3339`.

`rfc3339` converts the first element of a `unixXX` payload into a UTC
timestamp string.

## Response Model

Successful numeric conversion:

```json
{ "value_u16be": [17142,58982] }
```

Successful string conversion:

```json
{ "value_hex": "deadbeef" }
{ "value_rfc3339": "2026-02-27T07:03:24Z" }
```

Error response:

```json
{ "error_code": 3 }
```

Error codes are defined in `Docs/PRISM_PROTOCOL.md` and include
`1`=invalid request, `2`=payload exclusivity, `3`=unknown format, etc.

## Example Scenarios

1. **Numeric→Numeric**
   - Convert two big-endian 16-bit values to unsigned 32-bit little-endian:
   ```json
   {"payload_format":"u16be","target_format":"u32le","payload_u16":[17142,58982]}
   ```
   Response: `{ "value_u32le": [1118531840, ...] }` (computed values.)

2. **Hex string input → conversion**
   ```json
   {"payload_format":"hex","target_format":"u16be","payload_hex":"42f6e666"}
   ```
   Response: `{ "value_u16be": [17142,58982] }`

3. **Bytes→Hex output**
   ```json
   {"payload_format":"bytes","target_format":"hex","payload_bytes":"Qvbm"}
   ```
   Response: `{ "value_hex":"42f6e666" }`

4. **Unix timestamp → RFC3339 UTC**
   ```json
   {"payload_format":"unix32be","target_format":"rfc3339","payload_u32":[1581231230]}
   ```
   Response: `{ "value_rfc3339":"2020-02-09T15:53:50Z" }`

5. **Invalid format**
   ```json
   {"payload_format":"foo","target_format":"u8","payload_bytes":"AQ=="}
   ```
   => `{ "error_code": 3 }`

6. **Multiple payloads**
   ```json
   {"payload_format":"u8","target_format":"hex","payload_u8":[1],"payload_bytes":"AQ=="}
   ```
   => `{ "error_code": 2 }` (exclusivity violation)

7. **Length mismatch**
   ```json
   {"payload_format":"i16be","target_format":"u8","payload_bytes":"AQ"}
   ```
   => `{ "error_code": 5 }`

## CLI Usage

```
prism.exe -port 8080            # specify listening port
```

No other CLI flags exist.  Configuration file load still supported but
not covered here.

## Testing

To run the unit tests:

```bash
go test ./...
```

All existing tests exercise a variety of conversions.  New tests cover
hex alias, unix aliases, rfc3339 UTC, etc.

---

Refer to `Docs/PRISM_PROTOCOL.md` for the authoritative protocol
specification and error codes.
