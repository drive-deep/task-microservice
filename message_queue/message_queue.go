package message_queue

type MessageQueue interface {
    Connect(brokers []string, groupID string) (MessageQueue, error)
    SendMessage(topic string, message []byte) error
    StartConsuming(topics []string)
    Close() error
}