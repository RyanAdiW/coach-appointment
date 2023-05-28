package models

type DefaultResponse struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type SuccessResponseWithData struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// NewInternalServerErrorResponse default internal server error response
func SuccessOperationDefault(status, message string) DefaultResponse {
	return DefaultResponse{
		200,
		status,
		message,
	}
}

func SuccessOperationWithData(status, message string, data interface{}) SuccessResponseWithData {
	return SuccessResponseWithData{
		200,
		status,
		message,
		data,
	}
}

// NewBadRequestResponse default not found error response
func BadRequest(status, message string) DefaultResponse {
	return DefaultResponse{
		400,
		status,
		message,
	}
}

// NewInternalServerErrorResponse default internal server error response
func InternalServerError(status, message string) DefaultResponse {
	return DefaultResponse{
		500,
		status,
		message,
	}
}
