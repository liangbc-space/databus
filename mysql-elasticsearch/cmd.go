package mysql_elasticsearch

import (
	"databus/utils"
	"fmt"
	"github.com/panjf2000/ants/v2"
	"os"
	"os/signal"
	"syscall"
)

func Run() {
	//	获取topics
	topics := getTopics()

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	control := utils.SignalEvent(signalChan, signalHandleFunc)

	p, err := ants.NewPoolWithFunc(len(topics), execute)
	if err != nil {
		panic(err)
	}

	for _, topic := range topics {
		control.Add(1)

		args := make(map[string]interface{}, 2)
		args["topic"] = topic
		args["signal_control"] = control

		if err := p.Invoke(args); err != nil {
			panic(err)
		}
	}

	control.Wait()

}

func signalHandleFunc(control *utils.SignalControl) {
	<-control.SignalChan
}

func execute(args interface{}) {
	argsMap := args.(map[string]interface{})

	topic := argsMap["topic"].(string)
	control := argsMap["signal_control"].(*utils.SignalControl)

	consumer := createConsumerInstance()
	fmt.Printf("%p\n", consumer)

	defer consumer.Close()
	defer control.Done()

	for {

		select {
		case <-control.SignalChan:
			fmt.Println("收到退出信号")
			return
		default:
			messages := pullMessages(consumer, topic, 10)

			for _, message := range messages {
				fmt.Println(message.Value)
			}
		}

	}

}
