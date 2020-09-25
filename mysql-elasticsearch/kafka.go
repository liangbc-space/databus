package mysql_elasticsearch

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/liangbc-space/databus/system"
	"github.com/liangbc-space/databus/utils"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"os"
	"regexp"
)

var logger *zap.Logger

func getTopics() (topics []string) {

	consumer := createConsumerInstance()
	defer consumer.Close()

	reg := regexp.MustCompile(`^cn01_db.z_goods_(\d{2,3})$`)
	topics = consumer.GetTopics(reg)

	return topics
}

func pullMessages(consumer *utils.ConsumerInstance) (messages *kafka.Message) {

	event := consumer.Poll(100)
	if event == nil {
		return nil
	}

	switch e := event.(type) {
	case *kafka.Message:
		if system.ApplicationCfg.KafkaConfig.ConsumerLogs {
			if logger == nil {
				logger = getLogger()
			}

			logger.Info(consumer.String()+"成功消费到数据", zap.String("message", e.String()))
		}

		return e
		//consumerMessages(e)
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

	return messages
}

func createConsumerInstance() (consumer *utils.ConsumerInstance) {
	configMap := utils.ConsumerConfig{}
	configMap["session.timeout.ms"] = 6000
	configMap["auto.offset.reset"] = "earliest"

	consumer = new(utils.ConsumerInstance)
	consumer.Consumer = configMap.ConsumerInstance("binlog-canal-elasticsearch", false)

	return consumer
}

func getLogger() *zap.Logger {

	loggerCfg := utils.LoggerCfg{
		Level: zap.InfoLevel,
		Hook: lumberjack.Logger{
			Filename:   "logs/kafka-messages.log",
			MaxAge:     5,
			MaxBackups: 10,
			MaxSize:    512,
			Compress:   true,
			LocalTime:  true,
		},
		WithCaller: false,
	}

	return loggerCfg.NewLogger()
}
