package flipkart

type Broker struct {
	topics []Topic
}

type Topic struct {
	name      string
	messages  []string
	consumers []Consumer
}

type Consumer struct {
	name    string
	offset  int
	workCh  chan string
	handler func(msg string)
}

type BrokerInterface interface {
	AddTopic(name string) error
	DeleteTopic(name string) error
	Publish(topic string, message string) error
	AddConsumer(topic string, fn func(msg string)) (string, error)
	NextMessage(string) (string, bool, error)
	SetOffset(cName string, val int) error
	GetOffset(cName string) (int, error)

	TopicCount() int
	MsgsCount(topic string) int
	ConsumerCount(topic string) int
}
