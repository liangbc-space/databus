package exception

import (
	"runtime"
)

type caller struct {
	pc   uintptr
	file string
	line int
}

type Exception struct {
	message    string
	code       int
	stackTrace []byte
	*caller
}

type catch struct {
	hasCatch  bool
	exception Exception
}

type finalHandler interface {
	Finally(handler func())
}

type exceptionHandler interface {
	StackTrace() string
	Message() string
	Code() int
	File() string
	Line() int
	Pc() uintptr
}

type catchHandler interface {
	Catch(func(ex Exception)) finalHandler
	finalHandler
}

func Throw(message string, code int) {
	b := make([]byte, 1<<16)
	runtime.Stack(b, false)

	pc, file, line, _ := runtime.Caller(1)
	panic(Exception{
		message,
		code,
		b,
		&caller{pc, file, line},
	})
}

func Try(handler func()) catchHandler {
	ch := new(catch)

	defer func() {
		defer func() {
			r := recover()
			if r == nil {
				return
			}

			ex := Exception{}

			if exception, ok := r.(Exception); ok {
				ex = exception
			} else {
				pc, file, line, _ := runtime.Caller(2)
				if err, ok := r.(runtime.Error); ok {
					ex.message = err.Error()
					pc, file, line, _ = runtime.Caller(3)
				} else if err, ok := r.(error); ok {
					ex.message = err.Error()
				} else if message, ok := r.(string); ok {
					ex.message = message
				} else {
					panic(r)
				}

				b := make([]byte, 1<<16)
				runtime.Stack(b, false)
				ex.stackTrace = b
				ex.caller = &caller{pc, file, line}

			}

			ch.exception = ex
			ch.hasCatch = true
		}()
		handler()
	}()

	return ch
}

func (t *catch) Catch(handler func(exception Exception)) finalHandler {
	if t.hasCatch {
		handler(t.exception)
		t.hasCatch = false
	}

	return t
}

func (t *catch) Finally(handler func()) {
	defer handler()

	if t != nil && t.hasCatch {
		panic(t.exception)
	}
}

func (ex Exception) Pc() uintptr {
	return ex.pc
}

func (ex Exception) File() string {
	return ex.file
}

func (ex Exception) Line() int {
	return ex.line
}

func (ex Exception) StackTrace() string {
	return string(ex.stackTrace)
}

func (ex Exception) Message() string {
	return ex.message
}

func (ex Exception) Code() int {
	return ex.code
}
