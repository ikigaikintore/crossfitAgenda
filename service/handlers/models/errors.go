package models

import (
	"fmt"
	"net/http"
)

type common struct {
	msg  string
	code int
}

func (c common) Wrap() string {
	return fmt.Sprintf("err: %s: %d", c.msg, c.code)
}

type ErrNotAuth struct {
	common
}

func (e *ErrNotAuth) Code() int {
	return e.code
}

type ErrInternal struct {
	common
}

func (e *ErrInternal) Code() int {
	return e.code
}

type AgendaError interface {
	Error() string
	Code() int
}

func NewErrNotAuth() AgendaError {
	return &ErrNotAuth{common{
		msg:  "user not authenticated",
		code: http.StatusUnauthorized,
	}}
}

func (e *ErrNotAuth) Error() string {
	return e.Wrap()
}

func NewErrInternal() AgendaError {
	return &ErrInternal{common{
		msg:  "internal error",
		code: http.StatusInternalServerError,
	}}
}

func (e *ErrInternal) Error() string {
	return e.Wrap()
}
