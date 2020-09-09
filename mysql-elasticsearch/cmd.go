package mysql_elasticsearch

import (
	"databus/utils"
	"encoding/json"
	"fmt"
	"github.com/panjf2000/ants/v2"
	"os"
	"os/signal"
	"strings"
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
	if err := consumer.Subscribe(topic, nil); err != nil {
		panic(err)
	}

	defer consumer.Close()
	defer control.Done()

	for {

		select {
		case <-control.SignalChan:
			fmt.Println("收到退出信号")
			return
		default:

			message := pullMessages(consumer)
			if message == nil {
				continue
			}

			optionData := make(map[string]interface{})
			json.Unmarshal(message.Value, &optionData)
			data := optionData["data"].([]interface{})

			list := make([]map[string]interface{}, 0)
			for _, item := range data {
				item := item.(map[string]interface{})
				goodsData := make(map[string]interface{})

				goodsData["goods_id"] = item["id"]
				goodsData["store_id"] = item["store_id"]
				goodsData["operation_type"] = strings.ToUpper(optionData["type"].(string))
				list = append(list, goodsData)
			}

			fmt.Println(list)
		}

	}

}
