// internal/errors/codes.go

package errors

// Error codes as defined in the protocol document. These constants are
// exported so that various packages can return appropriate errors without
// importing protocol knowledge directly.

const (
    CodeInvalidRequest      = 1
    CodePayloadExclusivity  = 2
    CodeUnknownFormat       = 3
    CodeInvalidPayload      = 4
    CodeLengthMismatch      = 5
    CodeUnsupported         = 6
)
