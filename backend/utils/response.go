package utils

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func SuccessResponse(data interface{}, message string) Response {
	return Response{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func ErrorResponse(message string) Response {
	return Response{
		Success: false,
		Message: message,
		Data:    nil,
	}
}
