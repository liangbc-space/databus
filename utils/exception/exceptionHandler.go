package exception

import (
	"fmt"
	"reflect"
	"runtime/debug"
)

type Exception struct {
	Message string
	Code int
	StackTrace []byte
}



type catchHandler struct {
	hasCatch bool
	catch interface{}
}

type CatchHandler interface {
	Catch(interface{}, func(ex interface{})) *catchHandler
	FinalHandler
}

func Throw(message string,code int, ex interface{}) {
	/*b := make([]byte, 1<<16)
	runtime.Stack(b, true)*/

	if ex == nil {
		ex = new(Exception)
	}

	rValue := reflect.ValueOf(ex)
	fmt.Println(rValue.Type())

	//rValue.FieldByName("Message").SetString(message)
	fmt.Println(rValue.FieldByName("Message"))
	rValue.FieldByName("Code").Set(reflect.ValueOf(code))
	rValue.FieldByName("StackTrace").SetBytes(debug.Stack())
	fmt.Println(ex)

	panic(ex)
}

type FinalHandler interface {
	Finally(handler func())
}

func Try(handler func()) *catchHandler {
	ch := new(catchHandler)

	defer func() {
		defer func() {
			r := recover()
			if r != nil {
				ch.catch = r
				ch.hasCatch = true
			}
		}()
		handler()
	}()

	return ch
}

func (t *catchHandler) Catch(ex interface{}, handler func(interface{})) *catchHandler {
	//<4>如果传入的error类型和发生异常的类型一致，则执行异常处理器，并将hasCatch修改为true代表已捕捉异常
	if reflect.TypeOf(ex) == reflect.TypeOf(t.catch) {
		handler(t.catch)
		t.hasCatch = false
	}
	return t
}

func (t *catchHandler) Finally(handler func()) {
	defer handler()

	if t != nil && t.hasCatch {

	}
}
