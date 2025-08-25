package errors

import (
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	common "github.com/quiz_be/services/core/application/common/dto"
	"github.com/quiz_be/services/core/context"
	"github.com/quiz_be/services/core/i18n"
)

func InternalDefault(ctx context.Context) error {
	return WithCode(ctx, codes.Internal, LocKeyInternalServerError)
}

func WithCode(ctx context.Context, code codes.Code, locKey LocKey, args ...interface{}) error {
	message := i18n.Lt(ctx, locKey, args...)
	details := &common.AppError{Code: locKey.String(), LocalizedMessage: message}
	return detailError(ctx, code, message.Text, details)
}

func InvalidArgument(ctx context.Context, locKey LocKey, args ...interface{}) error {
	if locKey == "" {
		locKey = LocKeyInvalidArgumentError
	}
	return WithCode(ctx, codes.InvalidArgument, locKey, args...)
}

func NotFound(ctx context.Context, locKey LocKey, args ...interface{}) error {
	if locKey == "" {
		locKey = LocKeyNotFoundError
	}
	return WithCode(ctx, codes.NotFound, locKey, args...)
}

func FailedPreCondition(ctx context.Context, locKey LocKey, args ...interface{}) error {
	if locKey == "" {
		locKey = LocKeyFailedPreCondition
	}
	return WithCode(ctx, codes.FailedPrecondition, locKey, args...)
}

func DatabaseUnavailable(ctx context.Context, debug string) error {
	st := status.New(codes.Unavailable, "database_unavailable")
	appErr := &common.AppError{
		Code: "database_unavailable",
		LocalizedMessage: &common.LocalizedText{
			Text: "Không thể kết nối đến cơ sở dữ liệu. Vui lòng thử lại sau.",
			Key:  "database_unavailable",
		},
	}
	stWithDetails, err := st.WithDetails(appErr)
	if err != nil {
		// fallback nếu thêm details lỗi
		return st.Err()
	}
	return stWithDetails.Err()
}

func detailError(ctx context.Context, code codes.Code, message string, details proto.Message) error {
	detailStatus, err := status.New(code, message).WithDetails(details)
	if err != nil {
		return status.Error(codes.Internal, LocKeyInternalServerError.String())
	}

	return detailStatus.Err()
}

func Code(err error) codes.Code {
	codeStatus, _ := status.FromError(err)
	return codeStatus.Code()
}

func InvalidArgumentCode() codes.Code {
	return codes.InvalidArgument
}
