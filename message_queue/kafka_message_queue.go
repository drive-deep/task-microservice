package message_queue

import (
	"context"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
	"github.com/drive-deep/task-microservice/services"
)

type KafkaMessageQueue struct {
	producer      sarama.SyncProducer
	consumerGroup sarama.ConsumerGroup
	taskService   *services.TaskService
}

func NewKafkaMessageQueue(taskService *services.TaskService) *KafkaMessageQueue {
	return &KafkaMessageQueue{
		taskService: taskService,
	}
}

func (kmq *KafkaMessageQueue) Connect(brokers []string, groupID string) (MessageQueue, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, err
	}

	kmq.producer = producer
	kmq.consumerGroup = consumerGroup
	return kmq, nil
}

func (kmq *KafkaMessageQueue) SendMessage(topic string, message []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(message),
	}

	_, _, err := kmq.producer.SendMessage(msg)
	if err != nil {
		log.Printf("Failed to send message to Kafka: %v", err)
		return err
	}

	return nil
}

func (kmq *KafkaMessageQueue) StartConsuming(topics []string) {
	consumer := KafkaConsumer{
		ready:       make(chan bool),
		taskService: kmq.taskService,
	}

	go func() {
		for {
			if err := kmq.consumerGroup.Consume(context.Background(), topics, &consumer); err != nil {
				log.Printf("Error from consumer: %v", err)
			}
			// Check if context was cancelled, signaling that the consumer should stop
			if context.Background().Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready // Await till the consumer has been set up
	log.Println("Kafka consumer up and running!")
}

func (kmq *KafkaMessageQueue) Close() error {
	if err := kmq.producer.Close(); err != nil {
		return err
	}
	return kmq.consumerGroup.Close()
}

type KafkaConsumer struct {
	ready       chan bool
	taskService *services.TaskService
}

func (consumer *KafkaConsumer) Setup(_ sarama.ConsumerGroupSession) error {
	close(consumer.ready)
	return nil
}

func (consumer *KafkaConsumer) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *KafkaConsumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
		switch message.Topic {
		case "task_create":
			consumer.handleTaskCreate(message.Value)
		case "task_update":
			consumer.handleTaskUpdate(message.Value)
		case "task_delete":
			consumer.handleTaskDelete(message.Value)
		}
		sess.MarkMessage(message, "")
	}
	return nil
}

func (consumer *KafkaConsumer) handleTaskCreate(message []byte) {
	var task services.Task
	if err := json.Unmarshal(message, &task); err != nil {
		log.Printf("Failed to unmarshal task create message: %v", err)
		return
	}

	if err := consumer.taskService.CreateTask(&task); err != nil {
		log.Printf("Failed to create task: %v", err)
	}
}

func (consumer *KafkaConsumer) handleTaskUpdate(message []byte) {
	var task services.Task
    if err := json.Unmarshal(message, &task); err != nil {
        log.Printf("Failed to unmarshal task create message: %v", err)
        return
    }

    if err := consumer.taskService.UpdateTask(&task); err != nil {
        log.Printf("Failed to create task: %v", err)
    }
}

func (consumer *KafkaConsumer) handleTaskDelete(message []byte) {
	var task services.Task
	if err := json.Unmarshal(message, &task); err != nil {
		log.Printf("Failed to unmarshal task create message: %v", err)
		return
	}

	if err := consumer.taskService.DeleteTask(task.ID); err != nil {
		log.Printf("Failed to create task: %v", err)
	}
}

