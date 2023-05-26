package util

import "github.com/gofiber/fiber/v2"

type Error struct {
	Status  int    	`json:"status"`
	ErrCode int 	`json:"code"`
	Message string 	`json:"message"`
	Data    *fiber.Ctx
}

func (e *Error) Error() string {
	return e.Message
}

func HttpNotFound(m string) *Error {
	return &Error{Status: fiber.StatusNotFound, ErrCode: 404, Message: m}
}

func HttpBad(m string) *Error {
	return &Error{Status: fiber.StatusBadRequest, ErrCode: 400, Message: m}
}

func HttpError(m string) *Error {
	return &Error{Status: fiber.StatusInternalServerError, ErrCode: 500,  Message: m}
}

func HttpOk(m string, data *fiber.Ctx) *Error {
	return &Error{Status: fiber.StatusOK, ErrCode: 200, Message: m, Data: data}
}

func HttpForbidden(m string) *Error {
	return &Error{Status: fiber.StatusForbidden, ErrCode: 403, Message: m}
}

func HttpUnauthorized(m string) *Error {
	return &Error{Status: fiber.StatusUnauthorized, ErrCode: 403, Message: m}
}