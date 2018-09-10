package errors

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

// Error represents internal errors
type Error struct {
	Code
	Message string
	Err     error
}

// Code represents type of internal error
type Code int

const (
	// UpstreamConnFailure represents failure to connect to the upstream gRPC service
	UpstreamConnFailure Code = 2
	// ServiceUnresolvable represents failure to resolve a gRPC service to its upstream FQDN
	ServiceUnresolvable Code = 3
	// ServiceNotFound represents a missing gRPC service in an upstream, even though the service resolved to that upstream
	ServiceNotFound Code = 4
	// MethodNotFound represents a missing gRPC method in an upstream
	MethodNotFound Code = 5
	// MessageTypeMismatch represents user provided JSON not matching the message's type
	MessageTypeMismatch Code = 6
	// Unknown represents an unknown internal error
	Unknown Code = 8
	// VersionNotSpecified represents the user not specifying the upstream version when it is required.
	VersionNotSpecified Code = 9
	// VersionUndecidable represents there being multiple upstreams that match the specified (service, version) pair
	VersionUndecidable Code = 10
)

func (e *Error) Error() string {
	switch e.Code {
	case UpstreamConnFailure:
		return "could not connect to backend gRPC service"
	case ServiceUnresolvable:
		return "could not resolve service"
	case ServiceNotFound:
		return "service not found; service discovery error"
	case MethodNotFound:
		return "no such gRPC method"
	case MessageTypeMismatch:
		return "message type mismatch"
	case VersionNotSpecified:
		return "multiple versions of this service exist. specify version in request"
	case VersionUndecidable:
		return "multiple backends exist. add version annotations"
	default:
		return "unknown failure"
	}
}

// GRPCError is an error returned by gRPC upstream
type GRPCError struct {
	StatusCode int    `json:"code"`
	Message    string `json:"message"`
}

// HTTPStatusCode converts gRPC status codes to HTTP status codes
func (e *GRPCError) HTTPStatusCode() int {
	c := codes.Code(e.StatusCode)
	switch c {
	case codes.OK:
		return http.StatusOK
	case codes.Canceled:
		return http.StatusRequestTimeout
	case codes.Unknown:
		return http.StatusInternalServerError
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.DeadlineExceeded:
		return http.StatusRequestTimeout
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.ResourceExhausted:
		return http.StatusServiceUnavailable
	case codes.FailedPrecondition:
		return http.StatusPreconditionFailed
	case codes.Aborted:
		return http.StatusConflict
	case codes.OutOfRange:
		return http.StatusBadRequest
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.Internal:
		return http.StatusInternalServerError
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	case codes.DataLoss:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

func (e *GRPCError) Error() string {
	return e.Message
}
