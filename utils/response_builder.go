package utils

import (
	"net/http"
	"strings"

	constant "github.com/sipkyjayaputra/ticketing-system/constants"
)

func BuildSuccessResponse(data interface{}) *ResponseContainer {
	return &ResponseContainer{
		Response: Response{
			StatusCode:      http.StatusOK,
			Success:         true,
			ResponseMessage: &constant.RESPONSE_MESSAGE_SUCCESS,
			Errors:          nil,
			Data:            data,
			Info:            nil,
		},
	}
}

func BuildSuccessResponseWithInfo(data interface{}, info *ResponseInfo) *ResponseContainer {
	return &ResponseContainer{
		Response: Response{
			StatusCode:      http.StatusOK,
			Success:         true,
			ResponseMessage: &constant.RESPONSE_MESSAGE_SUCCESS,
			Errors:          nil,
			Data:            data,
			Info:            info,
		},
	}
}

func BuildDataNotFoundResponse() *ErrorContainer {
	return &ErrorContainer{
		Response: Response{
			StatusCode:      http.StatusNotFound,
			Success:         false,
			ResponseMessage: &constant.RESPONSE_MESSAGE_DATA_NOT_FOUND,
			Errors:          nil,
			Data:            nil,
			Info:            nil,
		},
	}
}

func BuildDataNotFoundResponseWithMessage(msg string) *ErrorContainer {
	return &ErrorContainer{
		Response: Response{
			StatusCode:      http.StatusNotFound,
			Success:         false,
			ResponseMessage: &constant.RESPONSE_MESSAGE_DATA_NOT_FOUND,
			Errors:          strings.Split(msg, "\n"),
			Data:            nil,
			Info:            nil,
		},
	}
}

func BuildBadRequestResponse(errMessage, throwable string) *ErrorContainer {
	return &ErrorContainer{
		Response: Response{
			StatusCode:      http.StatusBadRequest,
			Success:         false,
			ResponseMessage: &errMessage,
			Errors:          strings.Split(throwable, "\n"),
			Data:            nil,
			Info:            nil,
		},
	}
}

func BuildBadRequestResponseWithData(errCode, respCode, errMessage, throwable string, data interface{}) *ErrorContainer {
	return &ErrorContainer{
		Response: Response{
			StatusCode:      http.StatusBadRequest,
			Success:         false,
			ResponseMessage: &errMessage,
			Errors:          strings.Split(throwable, "\n"),
			Data:            data,
			Info:            nil,
		},
	}
}

func BuildInternalErrorResponse(errMessage, throwable string) *ErrorContainer {
	return &ErrorContainer{
		Response: Response{
			StatusCode:      http.StatusInternalServerError,
			Success:         false,
			ResponseMessage: &errMessage,
			Errors:          strings.Split(throwable, "\n"),
			Data:            nil,
			Info:            nil,
		},
	}
}

func BuildRouteNotFoundResponse() *ErrorContainer {
	return &ErrorContainer{
		Response: Response{
			StatusCode:      http.StatusNotFound,
			Success:         false,
			ResponseMessage: &constant.RESPONSE_MESSAGE_ROUTE_NOT_FOUND,
			Errors:          nil,
			Data:            nil,
			Info:            nil,
		},
	}
}

func BuildEmptyBodyReqResponse(errMessage, throwable string) *ErrorContainer {
	return &ErrorContainer{
		Response: Response{
			StatusCode:      http.StatusBadRequest,
			Success:         false,
			ResponseMessage: &errMessage,
			Errors:          strings.Split(throwable, "\n"),
			Data:            nil,
		},
	}
}

func BuildInvalidTypeResponse(errMessage, throwable string) *ErrorContainer {
	return &ErrorContainer{
		Response: Response{
			StatusCode:      http.StatusBadRequest,
			Success:         false,
			ResponseMessage: &errMessage,
			Errors:          strings.Split(throwable, "\n"),
			Data:            nil,
		},
	}
}

func BuildUnauthorizedResponse(errMessage, throwable string) *ErrorContainer {
	return &ErrorContainer{
		Response: Response{
			StatusCode:      http.StatusUnauthorized,
			Success:         false,
			ResponseMessage: &errMessage,
			Errors:          strings.Split(throwable, "\n"),
			Data:            nil,
		},
	}
}

func BuildTimeoutResponse(throwable string) *ErrorContainer {
	return &ErrorContainer{
		Response: Response{
			StatusCode:      http.StatusRequestTimeout,
			Success:         false,
			ResponseMessage: &constant.RESPONSE_MESSAGE_TIMEOUT,
			Errors:          strings.Split(throwable, "\n"),
			Data:            nil,
		},
	}
}

func BuildForbiddenAccessResponse(throwable string) *ErrorContainer {
	return &ErrorContainer{
		Response: Response{
			StatusCode:      http.StatusForbidden,
			Success:         false,
			ResponseMessage: &constant.RESPONSE_MESSAGE_FORBIDDEN_ACCESS,
			Errors:          strings.Split(throwable, "\n"),
			Data:            nil,
		},
	}
}
