package errors

import (
	"fmt"
	"io"
	"strconv"
)

type withCode struct {
	error
	code  int
	cause error
	*stack
}

func WithCode(code int, err error) error {
	return &withCode{
		error: err,
		code:  code,
		stack: callers(),
	}
}

func WithCodef(code int, format string, args ...any) error {
	return &withCode{
		error: fmt.Errorf(format, args...),
		code:  code,
		stack: callers(),
	}
}

func WrapC(cause error, code int, format string, args ...any) error {
	if cause == nil {
		return nil
	}

	return &withCode{
		error: fmt.Errorf(format, args...),
		code:  code,
		cause: cause,
		stack: callers(),
	}
}

func (w *withCode) Error() string {
	format := "%d:%s; cause:%s"

	return fmt.Sprintf(format, w.code, w.error.Error(), w.cause.Error())
}

// todo : 支持json格式输出
func (w *withCode) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			// fmt.Fprintf(s, "%+v\n", w.Error())
			io.WriteString(s, w.Error())
			w.stack.Format(s, verb)
			return
		} else if s.Flag('-') {
			io.WriteString(s, w.Error())
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, strconv.Itoa(w.code)+": "+w.error.Error())
	case 'q':
		fmt.Fprintf(s, "%q", strconv.Itoa(w.code))
	}
}

func (w *withCode) Cause() error { return w.cause }

func (w *withCode) Unwrap() error { return w.cause }
