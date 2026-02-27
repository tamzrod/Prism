# ARCHITECTURE.md

PRISM --- Deterministic Binary Transformation Engine Architecture (V1)

------------------------------------------------------------------------

## Philosophy

PRISM is:

-   Deterministic
-   Stateless
-   Byte-aligned
-   Semantics-free

PRISM transforms representation only.

(payload + payload_format + target_format) → exactly ONE output value

Errors are returned as numeric error_code defined in PRISM_PROTOCOL.md.

------------------------------------------------------------------------

## Core Library Responsibilities

-   Validate request structure
-   Enforce payload exclusivity
-   Normalize payload\_\* to \[\]byte
-   Decode source representation (explicit be/le only)
-   Encode target representation
-   Return exactly one output value or error_code

Constraints:

-   Go stdlib only
-   No protocol knowledge
-   No word swapping
-   No implicit defaults

------------------------------------------------------------------------

## Internal Value Model

Internal primitives may use:

-   int64
-   uint64
-   float64
-   string
-   \[\]uint8

External response shape is defined strictly by PRISM_PROTOCOL.md.

------------------------------------------------------------------------

## Format Registry

Formats are explicit string codes.

Examples:

-   i32be
-   f32le
-   unix64be
-   u16be
-   hex

Unknown format → error_code 3.

------------------------------------------------------------------------

## Determinism Guarantees

-   No implicit endian
-   No fallback parsing
-   No silent truncation
-   No word swapping

Ambiguity → error_code.

------------------------------------------------------------------------

PRISM is a pure codec layer.
