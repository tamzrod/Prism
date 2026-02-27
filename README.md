# PRISM

Deterministic binary transformation engine.

PRISM converts structured payloads between defined source and target formats.

It transforms representation --- nothing more.

------------------------------------------------------------------------

## Core Responsibility

(payload + payload_format + target_format) → exactly ONE output value (per protocol)

Numeric outputs are arrays.
String output is scalar.
Errors return numeric error_code.

------------------------------------------------------------------------

## Example 1: u16 → f32be

Request:

{ "payload_format": "u16be", "target_format": "f32be", "payload_u16": [17142, 58982] }

Response:

{ "value_f32be": [123.45] }

------------------------------------------------------------------------

## Example 2: hex → u16be

Request:

{ "payload_format": "bytes", "target_format": "u16be", "payload_hex": "42F6E666" }

Response:

{ "value_u16be": [17142, 58982] }

------------------------------------------------------------------------

## Error Example

{ "error_code": 3 }

------------------------------------------------------------------------

Protocol file is authoritative: PRISM_PROTOCOL.md

License: Apache 2.0