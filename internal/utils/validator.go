package utils

import (
	"errors"

	"buf.build/go/protovalidate"
	"github.com/fahrillrizal/ecommerce-grpc/pb/common"
	"google.golang.org/protobuf/proto"
)

func CheckValidation(req proto.Message) ([]*common.ValidationError, error) {
	if err := protovalidate.Validate(req); err != nil {
		var validationError *protovalidate.ValidationError
		if errors.As(err, &validationError) {
			var validationErrorResponses []*common.ValidationError = make([]*common.ValidationError, 0)
			for _, violation := range validationError.Violations {
				validationErrorResponses = append(validationErrorResponses, &common.ValidationError{
					Field:   *violation.Proto.Field.Elements[0].FieldName,
					Message: *violation.Proto.Message,
				})
			}

			return validationErrorResponses, nil
		}
		return nil, err
	}
	return nil, nil
}