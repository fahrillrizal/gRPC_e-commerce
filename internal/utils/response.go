package utils

import "github.com/fahrillrizal/ecommerce-grpc/pb/common"

func SuccessResponse(message string) *common.BaseResponse {
	return &common.BaseResponse{
		StatusCode: 200,
		Message:    message,
		IsError:    false,
	}
}

func BadRequestResponse(message string) *common.BaseResponse {
	return &common.BaseResponse{
		StatusCode: 400,
		Message:    message,
		IsError:    true,
	}
}

func UnauthorizedResponse(message string) *common.BaseResponse {
	return &common.BaseResponse{
		StatusCode: 401,
		Message:    message,
		IsError:    true,
	}
}

func ForbiddenResponse(message string) *common.BaseResponse {
	return &common.BaseResponse{
		StatusCode: 403,
		Message:    message,
		IsError:    true,
	}
}

func NotFoundResponse(message string) *common.BaseResponse {
	return &common.BaseResponse{
		StatusCode: 404,
		Message:    message,
		IsError:    true,
	}
}

func InternalServerErrorResponse(message string) *common.BaseResponse {
	return &common.BaseResponse{
		StatusCode: 500,
		Message:    message,
		IsError:    true,
	}
}

func ValidationErrorResponse(validationErrors []*common.ValidationError) *common.BaseResponse {
	return &common.BaseResponse{
		StatusCode:       400,
		Message:          "Validation failed",
		IsError:          true,
		ValidationErrors: validationErrors,
	}
}