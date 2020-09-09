package utils

import (
	"os"
	"sync"
)

type SignalControl struct {
	SignalChan chan os.Signal //信道
	*sync.WaitGroup
}

func SignalEvent(signalChan chan os.Signal, signalHandleFuc func(control *SignalControl)) (signalControl *SignalControl) {
	signalControl = new(SignalControl)
	signalControl.WaitGroup = new(sync.WaitGroup)
	signalControl.SignalChan = signalChan

	go func() {
		signalHandleFuc(signalControl)

		close(signalControl.SignalChan)

		return

	}()

	return signalControl

}
