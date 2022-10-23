package web

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/thewindear/go-web-easy-kit/pkg"
)

type Resp struct {
	// HttpStatus http状态码
	HttpStatus int
	// Code 业务代码
	Code int `json:"code,omitempty"`
	// Message 错误消息
	Message string `json:"message,omitempty"`
	// Data 数据
	Data interface{} `json:"data,omitempty"`
}

type RespError struct {
	Resp
	Err error `json:"-"`
	// Errors 更多的错误细节
	Errors interface{} `json:"errors,omitempty"`
}

func (im *RespError) Unwrap() error {
	return im.Err
}

func (im *RespError) Error() string {
	errMessage := fmt.Sprintf("response error: http_status %d", im.HttpStatus)
	if im.Err != nil {
		errMessage += fmt.Sprintf("wrap error: %s", im.Err.Error())
	}
	if im.Code != 0 {
		errMessage += fmt.Sprintf(" custom code: %d", im.Code)
	}
	if im.Message != "" {
		errMessage += fmt.Sprintf(" message: %s", im.Message)
	}
	return errMessage
}

// ParseParamsError 解析参数失败
func ParseParamsError(oriErr error, message string) error {
	return &RespError{
		Resp: Resp{
			HttpStatus: fiber.StatusUnprocessableEntity,
			Message:    message,
		},
		Err: oriErr,
	}
}

// BadRequest 只有错误提示信息
func BadRequest(message string) error {
	return &RespError{
		Resp: Resp{
			HttpStatus: fiber.StatusBadRequest,
			Message:    message,
		},
	}
}

// BadRequestWithError 只有错误提示信息
func BadRequestWithError(message string, err error) error {
	return &RespError{
		Resp: Resp{
			HttpStatus: fiber.StatusBadRequest,
			Message:    message,
		},
		Err: err,
	}
}

// ValidationFailed 表单验证失败
func ValidationFailed(field []*pkg.InvalidField) error {
	return &RespError{
		Resp: Resp{
			HttpStatus: fiber.StatusBadRequest,
			Message:    "表单输入有误",
		},
		Errors: field,
	}
}

// Forbidden 无权限
func Forbidden(message string) error {
	return &RespError{
		Resp: Resp{
			HttpStatus: fiber.StatusForbidden,
			Message:    message,
		},
	}
}

// Unauthorized 未授权
func Unauthorized(message string) error {
	return &RespError{
		Resp: Resp{
			HttpStatus: fiber.StatusUnauthorized,
			Message:    message,
		},
	}
}

// Error 服务错误无法处理的错误
func Error(oriErr error) *RespError {
	return &RespError{
		Resp: Resp{
			HttpStatus: fiber.StatusInternalServerError,
			Message:    "服务异常,请稍后再试",
		},
		Err: oriErr,
	}
}

// DefaultNotFound 默认资源不存在提示
func DefaultNotFound() error {
	return NotFound("")
}

// Conflict 冲突
func Conflict(message string) error {
	return &RespError{
		Resp: Resp{
			HttpStatus: fiber.StatusConflict,
			Message:    message,
		},
	}
}

// NotFound 资源不存在
func NotFound(message string) error {
	if message == "" {
		message = "资源不存在"
	}
	return &RespError{
		Resp: Resp{
			HttpStatus: fiber.StatusNotFound,
			Message:    message,
		},
	}
}

type Pagination struct {
	TotalSize int         `json:"totalSize"`
	TotalPage int         `json:"totalPage"`
	Page      int         `json:"page"`
	Items     interface{} `json:"items"`
}
