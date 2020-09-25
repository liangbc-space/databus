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

