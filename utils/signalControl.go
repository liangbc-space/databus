package utils

import (
	"os"
	"sync"
)

type SignalControl struct {
	SignalChan chan os.Signal //信道
	*sync.WaitGroup
}

func NewSignalListener(signalChan chan os.Signal) (listener *SignalControl) {
	listener = new(SignalControl)
	listener.WaitGroup = new(sync.WaitGroup)
	listener.SignalChan = signalChan

	return listener
}

func (listener *SignalControl) SignalEvent(fn func()) () {

	go func() {
		fn() //调用自定义信号处理方法

		close(listener.SignalChan)

		return

	}()
}
