package broker

import (
	"github.com/segmentio/kafka-go"
	"net"
	"strconv"
	"wb_tech_l0/internal/infrastructure/config"
)

func InitKafka() error {
	conn, err := kafka.Dial("tcp", config.C.KafkaConfig.Brokers[0])
	if err != nil {
		return err
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return err
	}

	controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return err
	}
	defer controllerConn.Close()

	topicConfig := []kafka.TopicConfig{
		{
			Topic:             config.C.KafkaConfig.Topic,
			NumPartitions:     config.C.KafkaConfig.NumPartitions,
			ReplicationFactor: 1,
		},
	}

	err = controllerConn.CreateTopics(topicConfig...)
	if err != nil {
		return err
	}
	return nil
}
