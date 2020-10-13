package utils

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/liangbc-space/databus/system"
	"github.com/liangbc-space/databus/utils/exception"
	"go.uber.org/zap"
	"regexp"
	"strings"
)

//libkafka设置，请参见文档：	https://github.com/edenhill/librdkafka/tree/master/CONFIGURATION.md
type ConsumerConfig map[string]interface{}

//kafka消费者实例
type ConsumerInstance struct {
	*kafka.Consumer
}

func (consumerConfig ConsumerConfig) ConsumerInstance(groupId string, autoCommitOffset bool) *kafka.Consumer {
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

	consumer, err := kafka.NewConsumer(&configMap)
	if err != nil {
		logger := NewDefaultLogger()
		defer logger.Sync()

		logger.Panic("创建消费者连接失败："+err.Error(), zap.Reflect("connConfig", configMap))
		exception.Throw("创建消费者连接失败："+err.Error(), 1)
		return nil
	}

	return consumer
}

func (consumer *ConsumerInstance) GetTopics(reg *regexp.Regexp) (topics []string) {
	metadata, err := consumer.GetMetadata(nil, true, 100)

	if err != nil {
		logger := NewDefaultLogger()
		defer logger.Sync()

		logger.Panic("获取meta信息失败："+err.Error(), zap.String("connInfo", consumer.String()))
		exception.Throw("获取meta信息失败："+err.Error(), 1)
		return nil
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
