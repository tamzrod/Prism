# ARCHITECTURE.md

PRISM -- Deterministic Binary Transformation Engine Architecture (V1)

------------------------------------------------------------------------

## 1. Architectural Philosophy

PRISM is a **representation transformation engine**.

It is intentionally:

-   Deterministic
-   Stateless
-   Byte-aligned
-   Library-first
-   Semantics-free

PRISM transforms representation only.

(payload + payload_format + target_format) → exactly ONE output value

No side effects. No interpretation beyond binary representation.

------------------------------------------------------------------------

## 2. Layered Structure

PRISM has two logical layers.

### 2.1 Core Library (Primary)

This is the authoritative engine.

Responsibilities:

-   Validate request structure
-   Enforce payload exclusivity
-   Normalize payload\_\* to \[\]byte
-   Decode source representation
-   Encode target representation
-   Return exactly one value or structured error

Constraints:

-   Go stdlib only
-   No reflection-heavy logic
-   No external dependencies
-   No protocol knowledge
-   No semantic awareness

The core must be fully usable without HTTP.

------------------------------------------------------------------------

### 2.2 Microservice Wrapper (Optional)

Thin JSON-over-HTTP adapter.

Responsibilities:

-   JSON parsing
-   Request validation
-   Call core.Transform()
-   JSON encoding of response

The microservice MUST NOT:

-   Contain transformation logic
-   Duplicate decoding logic
-   Apply formatting rules not defined by target_format
-   Add defaults

All behavior must remain inside the core library.

------------------------------------------------------------------------

## 3. Data Flow Pipeline

The internal transformation pipeline is fixed.

1.  Validate request
2.  Validate exactly one payload\_\* field
3.  Normalize payload to \[\]byte
4.  Interpret bytes according to payload_format
5.  Convert to internal primitive representation
6.  Encode according to target_format
7.  Return one value

Any failure at any stage → deterministic error.

------------------------------------------------------------------------

## 4. Internal Value Model

PRISM returns exactly one output value.

Supported output categories:

-   Scalar number (int64, uint64, float64)
-   String
-   Array (e.g. \[\]uint16 or \[\]uint8)

Implementation recommendation:

Use a tagged union style struct:

-   Kind
-   Numeric
-   String
-   U16Array
-   U8Array

No reflection required.

------------------------------------------------------------------------

## 5. Format Registry Strategy

Formats are explicit string codes.

Examples:

-   i32be
-   f32le
-   unix64be
-   u16_array_be
-   hex

Implementation approach:

-   Use switch statements or static maps
-   No dynamic plugin loading in V1
-   Unknown formats → ERR_UNKNOWN_FORMAT

------------------------------------------------------------------------

## 6. Determinism Guarantees

PRISM must guarantee:

-   No implicit endian defaults
-   No fallback parsing
-   No automatic trimming
-   No silent truncation
-   No partial decoding
-   No format guessing

If input is ambiguous → error.

------------------------------------------------------------------------

## 7. Validation Rules

Strict fail-fast philosophy.

Examples:

-   Multiple payload\_\* fields → ERR_PAYLOAD_EXCLUSIVE
-   Invalid hex string → ERR_INVALID_PAYLOAD
-   Incorrect byte length → ERR_LENGTH_MISMATCH
-   Unsupported format → ERR_UNSUPPORTED

No warnings. No soft failures.

------------------------------------------------------------------------

## 8. Extension Policy (Future Versions)

New formats may be added if:

-   They remain byte-aligned
-   They preserve one-transform-per-call rule
-   They do not introduce semantics
-   They do not require multi-field decoding

Out of scope for PRISM:

-   Bit packing
-   Word swap logic (unless explicitly added in a future version)
-   Composite struct decoding
-   Protocol frame parsing

------------------------------------------------------------------------

## 9. Non-Goals (Architectural Guardrails)

PRISM must never evolve into:

-   A protocol decoder
-   A historian
-   A semantic interpreter
-   A business logic engine
-   A scaling/unit engine
-   A streaming processor

If functionality requires meaning, interpretation, or domain context, it
belongs outside PRISM.

------------------------------------------------------------------------

## 10. Architectural Integrity Principle

PRISM is a pure codec layer.

It exists to transform representation safely and deterministically.

Nothing more. Nothing less.

------------------------------------------------------------------------

End of ARCHITECTURE.md V1
