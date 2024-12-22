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

func WrapC(err error, code int, format string, args ...any) error {
	if err == nil {
		return nil
	}

	err = &withCode{
		error: fmt.Errorf(format, args...),
		code:  code,
		cause: err,
		stack: callers(),
	}

	return err
}

func (w *withCode) Error() string {
	format := "code(%d): %s."

	return fmt.Sprintf(format, w.code, w.error.Error())
}

// todo : 1.支持json格式输出  2.修改输出格式 3. 支持 %-v 输出最近的错误调用 4. 支持 %+v 输出所有的错误调用
func (w *withCode) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			io.WriteString(s, w.Error())
			w.stack.Format(s, verb)
			if err := w.Cause(); err != nil {
				fmt.Fprintf(s, "\n%+v", w.Cause())
			}
			return
		} else if s.Flag('-') {
			io.WriteString(s, w.Error())
			w.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, w.Error())
	case 'q':
		fmt.Fprintf(s, "%q", w.Error())
	case 'd':
		io.WriteString(s, strconv.Itoa(w.code))
	}
}

func (w *withCode) Cause() error { return w.cause }

func (w *withCode) Unwrap() error { return w.cause }
