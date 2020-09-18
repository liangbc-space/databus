package utils

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/liangbc-space/databus/system"
	"os"
	"regexp"
	"strings"
)

//libkafka设置，请参见文档：	https://github.com/edenhill/librdkafka/tree/master/CONFIGURATION.md
type ConsumerConfig map[string]interface{}

//kafka消费者实例
type ConsumerInstance struct {
	*kafka.Consumer
}

func (consumerConfig ConsumerConfig) ConsumerInstance(groupId string, autoCommitOffset bool) (consumer *kafka.Consumer) {
	configMap := kafka.ConfigMap{
		"bootstrap.servers":     strings.Join(system.ApplicationCfg.KafkaConfig.Brokers, ","),
		"broker.address.family": system.ApplicationCfg.KafkaConfig.BrokerAddressFamily,
	}

	for key, value := range consumerConfig {
		configMap[key] = value
	}

	//	防止被consumerConfig覆盖
	if groupId != "" {
		configMap["group.id"] = groupId
	}
	configMap["enable.auto.commit"] = autoCommitOffset

	var err error
	if consumer, err = kafka.NewConsumer(&configMap); err != nil {
		fmt.Fprintf(os.Stderr, "创建消费者连接失败: %s\n", err)
		os.Exit(1)
	}

	return consumer
}

func (consumer *ConsumerInstance) GetTopics(reg *regexp.Regexp) (topics []string) {
	metadata, err := consumer.GetMetadata(nil, true, 100)

	if err != nil {
		fmt.Fprintf(os.Stderr, "获取meta信息失败: %s\n", err)
		os.Exit(1)
	}

	for _, topicMetadata := range metadata.Topics {
		if reg != nil {
			if reg.MatchString(topicMetadata.Topic) {
				topics = append(topics, topicMetadata.Topic)
			}
		} else {
			topics = append(topics, topicMetadata.Topic)
		}
	}

	return topics
}

func (consumer *ConsumerInstance) ConsumerMessage(signChan chan os.Signal, timeOutMs uint, callback func(message *kafka.Message)) {
	for {
		select {
		case sig := <-signChan:
			fmt.Printf("收到信号：%v\n", sig)
			consumer.Close()
			os.Exit(0)
		default:
			ev := consumer.Poll(int(timeOutMs))
			if ev == nil {
				//fmt.Println("消息获取超时")
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				callback(e)
				/*fmt.Printf("%% Message on %s:\n%s\n",
					e.TopicPartition, string(e.Value))
				if e.Headers != nil {
					fmt.Printf("%% Headers: %v\n", e.Headers)
				}*/
			case kafka.Error:
				// Errors should generally be considered
				// informational, the client will try to
				// automatically recover.
				// But in this example we choose to terminate
				// the application if all brokers are down.
				fmt.Fprintf(os.Stderr, "%% Error: %v: %v\n", e.Code(), e)
				if e.Code() == kafka.ErrAllBrokersDown {
					break
				}
			default:
				fmt.Printf("Ignored %v\n", e)
			}
		}

	}
}
