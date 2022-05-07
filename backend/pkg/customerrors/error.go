package customerrors

import (
	"errors"
	"fmt"
	"runtime"
	"strings"

	"go.uber.org/zap/zapcore"
)

type ErrorCode string

const (
	ErrCodeNotFound        = ErrorCode("NotFound")
	ErrCodeInvalidArgument = ErrorCode("InvalidArgument")
	ErrUnexpected          = ErrorCode("Unexpected")
	ErrForbidden           = ErrorCode("Forbidden")
	ErrUnAuthorized        = ErrorCode("UnAuthorized")
	ErrCodeOptimisticLock  = ErrorCode("OptimisticLock")
	ErrCodeConflict        = ErrorCode("Conflict")
	ErrPayloadTooLarge     = ErrorCode("PayloadTooLarge")
	ErrRangeNotSatisfiable = ErrorCode("ErrRangeNotSatisfiable")
)

func IsServerError(errorCode ErrorCode) bool {
	return errorCode == ErrUnexpected
}

type trace struct {
	fileName string
	funcName string
	line     int
}

type field struct {
	key, value string
}

type CustomError struct {
	code           ErrorCode
	msg            string
	orgErr         error
	useTrace       bool
	trace          *trace
	adminAccountId string
	organizationId string
	loginAccountId string
	fields         []*field
}

// returns caller function file name, function name, and line number
// see https://stackoverflow.com/questions/25927660/how-to-get-the-current-function-name/46289376#46289376
func traceFunction() (string, string, int) {
	pc := make([]uintptr, 15)
	n := runtime.Callers(3, pc) // get caller of logging function
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return frame.File, frame.Function, frame.Line
}

func (e *CustomError) Error() string {
	if e.msg != "" {
		return fmt.Sprintf("code: %s, msg: %s", e.Code(), e.msg)
	}
	return fmt.Sprintf("code: %s", e.Code())
}

// zapcore.ObjectMarshalerを実装して、zap.ObjectでMarshalできるようにする
// https://pkg.go.dev/go.uber.org/zap/zapcore#ObjectMarshaler
func (e *CustomError) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("code", string(e.Code()))
	if e.msg != "" {
		enc.AddString("msg", e.msg)
	}
	if e.orgErr != nil {
		enc.AddString("cause", e.Cause())
	}
	if e.trace != nil {
		enc.AddString("traceFileName", e.trace.fileName)
		enc.AddString("traceFuncName", e.trace.funcName)
		enc.AddInt("traceLine", e.trace.line)
	}
	if e.organizationId != "" {
		enc.AddString("organizationId", e.organizationId)
	}
	if e.loginAccountId != "" {
		enc.AddString("loginAccountId", e.loginAccountId)
	}
	for _, v := range e.fields {
		enc.AddString(v.key, v.value)
	}
	return nil
}

func (e *CustomError) Cause() string {
	var sb strings.Builder
	for orgErr := e.orgErr; orgErr != nil; {
		sb.WriteString(orgErr.Error())
		sb.WriteString(": ")
		var customErr *CustomError
		if errors.As(orgErr, &customErr) {
			orgErr = customErr.orgErr
		} else {
			orgErr = nil
		}
	}
	return sb.String()
}

func (e *CustomError) Code() ErrorCode {
	return e.code
}

type ErrorOption func(*CustomError)

func WithMsg(msg string) ErrorOption {
	return func(e *CustomError) {
		e.msg = msg
	}
}

func WithError(err error) ErrorOption {
	return func(e *CustomError) {
		e.orgErr = err
	}
}

func WithUseTrace(use bool) ErrorOption {
	return func(e *CustomError) {
		e.useTrace = use
	}
}

func WithAdminAccountId(adminAccountId string) ErrorOption {
	return func(e *CustomError) {
		e.adminAccountId = adminAccountId
	}
}

func WithOrganizationId(organizationId string) ErrorOption {
	return func(e *CustomError) {
		e.organizationId = organizationId
	}
}

func WithLoginAccountId(accountId string) ErrorOption {
	return func(e *CustomError) {
		e.loginAccountId = accountId
	}
}

func WithField(key, value string) ErrorOption {
	return func(e *CustomError) {
		e.fields = append(e.fields, &field{key: key, value: value})
	}
}

func NewError(code ErrorCode, opts ...ErrorOption) error {
	e := &CustomError{
		code:     code,
		useTrace: true, // default true
	}
	for _, optionFunc := range opts {
		optionFunc(e)
	}

	if e.useTrace {
		fileName, funcName, line := traceFunction()
		e.trace = &trace{
			fileName: fileName,
			funcName: funcName,
			line:     line,
		}
	}

	return e
}
