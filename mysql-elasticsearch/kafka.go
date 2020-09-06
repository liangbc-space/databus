package mysql_elasticsearch

import (
	"databus/utils"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"
)

func poolMessages() {



	configMap := utils.ConsumerConfig{}
	configMap["session.timeout.ms"] = 6000
	configMap["auto.offset.reset"] = "earliest"

	consumer := new(utils.ConsumerInstance)
	consumer.Consumer = configMap.ConsumerInstance("test", false)

	reg := regexp.MustCompile("^cn01_db.z_goods_(\\d{2})$")
	topics := consumer.GetTopics(reg)
	consumer.Close()

	done := make(chan bool)



	for _, item := range topics {


		topic := []string{item}
		go func(topic []string) {

			signalChan := make(chan os.Signal, 1)
			signal.Notify(signalChan,
				os.Kill,
				os.Interrupt,
				syscall.SIGHUP,
				syscall.SIGINT,
				syscall.SIGTERM,
				syscall.SIGQUIT,
			)

			consumer := new(utils.ConsumerInstance)
			consumer.Consumer = configMap.ConsumerInstance("test", false)
			defer consumer.Close()

			if err := consumer.SubscribeTopics(topic, nil); err != nil {
				panic(err)
			}

			consumer.ConsumerMessage(signalChan, 100, consumerMessages)

			done <- true

		}(topic)

	}

	for i := 1; i <= len(topics); i++ {
		<-done
		if len(done) == 0 {
			break
		}
	}
}

func consumerMessages(message *kafka.Message) {
	/*fmt.Println("准备提交偏移量")
	time.Sleep(time.Second * 2)*/
	//consumer.Commit()

	if *message.TopicPartition.Topic == "cn01_db_z_goods_01" {
		time.Sleep(time.Second * 5)
	} else {
		time.Sleep(time.Second * 2)
	}
	fmt.Println(*message.TopicPartition.Topic)

}
