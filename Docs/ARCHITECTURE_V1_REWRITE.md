# ARCHITECTURE.md

PRISM --- Deterministic Binary Transformation Engine (V1)

------------------------------------------------------------------------

## Philosophy

PRISM is:

- Deterministic
- Stateless
- Byte-aligned
- Semantics-free

It transforms representation only.

(payload + payload_format + target_format) → exactly ONE output value

Errors are returned as numeric error_code (see PRISM_PROTOCOL.md).

------------------------------------------------------------------------

## Core Responsibilities

- Validate request structure
- Enforce payload exclusivity
- Normalize payload_* to []byte
- Decode representation (explicit be/le only)
- Encode target representation
- Return exactly one output value or error_code

Constraints:

- Go stdlib only
- No protocol knowledge
- No word swapping
- No implicit defaults
- No global state

------------------------------------------------------------------------

## Internal Value Model

Internal primitives may use:

- int64
- uint64
- float64
- string
- []uint8

External response shape is defined strictly by PRISM_PROTOCOL.md.

------------------------------------------------------------------------

## Determinism Rule

Ambiguity → error_code 3.

PRISM is a pure codec layer.