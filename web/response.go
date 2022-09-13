package web

import (
    "fmt"
    "github.com/gofiber/fiber/v2"
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
    OriError error `json:"-"`
    // Errors 更多的错误细节
    Errors interface{} `json:"errors,omitempty"`
}

func (im *RespError) Error() string {
    errMessage := fmt.Sprintf("response error: http_status %d", im.HttpStatus)
    if im.OriError != nil {
        errMessage += fmt.Sprintf("wrap error: %s", im.OriError.Error())
    }
    if im.Code != 0 {
        errMessage += fmt.Sprintf(" custom code: %d", im.Code)
    }
    if im.Message != "" {
        errMessage += fmt.Sprintf(" message: %s", im.Message)
    }
    return errMessage
}

// BadRequest 解析参数失败
func BadRequest(oriErr error) error {
    return &RespError{
        Resp: Resp{
            HttpStatus: fiber.ErrBadRequest.Code,
            Message:    fiber.ErrBadRequest.Message,
        },
        OriError: oriErr,
    }
}

// ValidationFailed 表单验证失败
func ValidationFailed(field []*InvalidField) error {
    return &RespError{
        Resp: Resp{
            HttpStatus: fiber.ErrUnprocessableEntity.Code,
            Message:    fiber.ErrUnprocessableEntity.Message,
        },
        OriError: nil,
        Errors:   field,
    }
}

// DefaultNotFound 默认资源不存在提示
func DefaultNotFound() error {
    return NotFound("")
}

func NotFound(message string) error {
    if message == "" {
        message = fiber.ErrNotFound.Message
    }
    return &RespError{
        Resp: Resp{
            HttpStatus: fiber.ErrNotFound.Code,
            Message:    message,
        },
        OriError: fiber.ErrNotFound,
    }
}

func Error(oriErr error) error {
    return &RespError{
        Resp: Resp{
            HttpStatus: fiber.ErrInternalServerError.Code,
            Message:    fiber.ErrInternalServerError.Message,
        },
        OriError: oriErr,
        Errors:   nil,
    }
}

type Pagination struct {
    TotalSize int         `json:"totalSize"`
    TotalPage int         `json:"totalPage"`
    Page      int         `json:"page"`
    Items     interface{} `json:"items"`
}
