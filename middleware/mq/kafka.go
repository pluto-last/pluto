package mq

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"os"
)

type Producer struct {
	producer *kafka.Producer
}

func InitProducer() (*Producer, error) {

	broker := "127.0.0.1:9092"

	var producer Producer
	var err error
	producer.producer, err = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker})

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}
	return &producer, err
}

func (p *Producer) PushList(topic, value string, partition int32) {

	doneChan := make(chan bool)

	go func() {
		defer close(doneChan)
		for e := range p.producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				m := ev
				if m.TopicPartition.Error != nil {
					fmt.Println(partition)
					fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
				} else {
					fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
						*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
				}
				return

			default:
				fmt.Printf("Ignored event: %s\n", ev)
			}
		}
	}()

	p.producer.ProduceChannel() <- &kafka.Message{TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny}, Value: []byte(value)}

	// wait for delivery report goroutine to finish
	_ = <-doneChan
}

func InitConsumer(topics []string, groupID int) (*kafka.Consumer, error) {

	broker := "127.0.0.1:9092"

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":               broker,
		"group.id":                        groupID,
		"session.timeout.ms":              6000,
		"go.events.channel.enable":        true,
		"go.application.rebalance.enable": true,
		"enable.partition.eof":            true,
		"auto.offset.reset":               "earliest"})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create consumer: %s\n", err)
	}

	fmt.Printf("Created Consumer %v\n", c)

	err = c.SubscribeTopics(topics, nil)

	return c, err
}
