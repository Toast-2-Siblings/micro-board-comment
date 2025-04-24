package utils

import "github.com/gin-gonic/gin"

// 응답 성공 시에 사용하는 구조체
type response struct {
	StatusCode int `json:"status_code"`
	Message string `json:"message"`
	Data interface{} `json:"data,omitempty"`
}

// 응답 실패 시에 사용하는 구조체
type errorResponse struct {
	StatusCode int `json:"status_code"`
	Data interface{} `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
}

// 응답 성공 시에 사용하는 응답 함수
func APIResponse(c *gin.Context, StatusCode int, Message string, Data interface{}) {
	jsonResponse := response{
		StatusCode: StatusCode,
		Message: Message,
		Data: Data,
	}

	if StatusCode >= 400 {
		c.JSON(StatusCode, jsonResponse)
		defer c.AbortWithStatus(StatusCode)
	} else {
		c.JSON(StatusCode, jsonResponse)
	}
}

// 응답 실패 시에 사용하는 응답 함수
func ErrorResponse(c *gin.Context, StatusCode int, Error string, Data interface{}) {
	errorResponse := errorResponse{
		StatusCode: StatusCode,
		Error: Error,
		Data: Data,
	}

	c.JSON(StatusCode, errorResponse)
	defer c.AbortWithStatus(StatusCode)
}
