package kafka

type Broker struct {
	topics []Topic
}

type Topic struct {
	name      string
	messages  []string
	consumers []Consumer
}

type Consumer struct {
	name   string // unique name
	offset int
}

type BrokerInterface interface {
	AddTopic(name string) error
	DeleteTopic(name string) error
	Publish(topic string, message string) error
	AddConsumer(topic string) (string, error)
	NextMessage(string) (string, bool, error)
	BlockingConsume(cName string)
	SetOffset(cName string, val int) error
	GetOffset(cName string) (int, error)

	TopicCount() int
	MsgsCount(topic string) int
	ConsumerCount(topic string) int
}
