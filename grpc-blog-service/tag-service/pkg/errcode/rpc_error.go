package errcode

import (
	pb "github.com/summerKK/go-code-snippet-library/grpc-blog-service/tag-service/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TogRPCError(err *Error) error {
	s, _ := status.New(TogRPCCode(err.Code()), err.Msg()).WithDetails(&pb.Error{
		Code:    int32(err.Code()),
		Message: err.Msg(),
	})

	return s.Err()
}

func TogRPCCode(code int) codes.Code {
	var statusCode codes.Code
	switch code {
	case Fail.Code():
		statusCode = codes.Internal
	case InvalidParams.Code():
		statusCode = codes.InvalidArgument
	case Unauthorized.Code():
		statusCode = codes.Unauthenticated
	case NotFound.Code():
		statusCode = codes.NotFound
	case Unknown.Code():
		statusCode = codes.Unknown
	case DeadlineExceeded.Code():
		statusCode = codes.DeadlineExceeded
	case AccessDenied.Code():
		statusCode = codes.PermissionDenied
	case LimitExceeded.Code():
		statusCode = codes.ResourceExhausted
	case MethodNotAllowed.Code():
		statusCode = codes.Unimplemented
	default:
		statusCode = codes.Unknown
	}

	return statusCode
}

type Status struct {
	*status.Status
}

func FromError(err error) *Status {
	s, _ := status.FromError(err)
	return &Status{s}
}

func TogPRCStatus(code int, msg string) *Status {
	s, _ := status.New(TogRPCCode(code), msg).WithDetails(&pb.Error{
		Code:    int32(code),
		Message: msg,
	})

	return &Status{s}
}
