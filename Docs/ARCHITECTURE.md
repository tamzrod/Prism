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

------------------------------------------------------------------------

## Programming Guidelines

To keep the codebase maintainable and consistent the following rules
must be observed when authoring Go source files:

1. **File header comment.** Every `.go` file begins with a single
   comment indicating its repository path, e.g.
   ```go
   // internal/engine/codec.go
   ```
   This makes it easy to locate the file when browsing or reviewing diffs.

2. **File length limit.** No file may exceed **300 lines** of actual code.
   If a package or type grows larger it should be split into multiple files
   (for example `validator.go` / `encoder.go` instead of one 400‑line file).

These guidelines are not enforced by tooling but they should be followed by
all contributors to keep the project readable and navigable.


------------------------------------------------------------------------

## Configuration

- **Format:** YAML file containing a top‑level `port` field. Example:
  ```yaml
  port: 12345
  ```
  Only one configuration file is required; if the file is absent the server
  will fall back to the `PRISM_PORT` environment variable. If neither is
  provided the default listening port is **12345**.
- Other configuration options may be added later, but port remains the only
  setting enforced by the implementation at this time.