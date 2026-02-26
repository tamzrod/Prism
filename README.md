# PRISM

**PRISM** is a deterministic binary transformation engine.

It converts structured payloads between defined source and target
formats.

PRISM does not interpret meaning.\
It does not apply business logic.\
It does not apply scaling or semantics.

It transforms representation --- nothing more.

------------------------------------------------------------------------

## Philosophy

PRISM follows strict architectural boundaries:

-   Byte-aligned only
-   One data type per request
-   No semantic awareness (no "voltage", no "status")
-   No scaling rules
-   No alarm logic
-   No device knowledge
-   No implicit defaults
-   Deterministic behavior

PRISM separates:

Binary representation\
from\
Meaning

------------------------------------------------------------------------

## Core Responsibility

Input payload + source format + target format\
→\
One transformed value

That is the entire contract.

------------------------------------------------------------------------

## V1 Scope

### Supported Payload Encodings (Transport Layer)

Exactly one must be provided per request:

-   `payload_u16`
-   `payload_u8`
-   `payload_hex`
-   `payload_base64`
-   `payload_str`

If none or more than one are present → error.

------------------------------------------------------------------------

### Supported Data Types (Initial)

**Integer** - `int16` - `uint16` - `int32` - `uint32` - `int64` -
`uint64`

**Floating Point** - `float32` - `float64`

**Timestamp** - `unix32` - `unix64`

**Raw** - `bytes` - `u16_array_be` - `u16_array_le`

All multi-byte formats must explicitly define endian (`be` / `le`).

No silent defaults.

------------------------------------------------------------------------

## Example

### Example: u16 Array → float32

Request:

{ "payload_format": "u16be", "target_format": "f32", "payload_u16":
\[17142, 58982\] }

Response:

123.45

------------------------------------------------------------------------

### Example: hex → u16 Array

Request:

{ "payload_format": "bytes", "target_format": "u16_array_be",
"payload_hex": "42F6E666" }

Response:

\[17142, 58982\]

------------------------------------------------------------------------

## Architecture

PRISM is:

-   Library-first
-   Zero external dependencies (stdlib only)
-   Deterministic
-   Explicit in behavior

Microservice support (JSON over HTTP) is an adapter layer, not the core.

Core engine operates on:

\[\]byte

All other payload encodings are normalized before transformation.

------------------------------------------------------------------------

## Non-Goals

PRISM does not:

-   Apply scaling
-   Interpret units
-   Perform validation beyond structural correctness
-   Store data
-   Log events
-   Manage memory
-   Handle protocol framing
-   Perform bit-level packing (byte-aligned only)

If you need semantics, use a higher layer.

------------------------------------------------------------------------

## Design Principle

PRISM behaves like a prism in optics:

Raw input enters.\
Structure is revealed.\
Meaning is not assigned.

------------------------------------------------------------------------

## License

Apache 2.0
