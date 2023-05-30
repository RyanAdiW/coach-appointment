package models

type DefaultResponse struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ResponseWithData struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func SuccessOperationDefault(status, message string) DefaultResponse {
	return DefaultResponse{
		200,
		status,
		message,
	}
}

func SuccessOperationWithData(status, message string, data interface{}) ResponseWithData {
	return ResponseWithData{
		200,
		status,
		message,
		data,
	}
}

func BadRequest(status, message string) DefaultResponse {
	return DefaultResponse{
		400,
		status,
		message,
	}
}

func BadRequestWithData(status, message string, data interface{}) ResponseWithData {
	return ResponseWithData{
		400,
		status,
		message,
		data,
	}
}

func UnauthorizedRequest(status, message string) DefaultResponse {
	return DefaultResponse{
		401,
		status,
		message,
	}
}

func NotFound(status, message string) DefaultResponse {
	return DefaultResponse{
		404,
		status,
		message,
	}
}

func InternalServerError(status, message string) DefaultResponse {
	return DefaultResponse{
		500,
		status,
		message,
	}
}
