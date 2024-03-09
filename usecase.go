package kafka

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

var mu sync.Mutex

func NewBroker() BrokerInterface {
	return &Broker{}
}

func (r *Broker) AddTopic(name string) error {
	for _, t := range r.topics {
		if t.name == name {
			return errors.New("topic already present")
		}
	}
	t := Topic{name: name}
	r.topics = append(r.topics, t)
	return nil
}

func (r *Broker) DeleteTopic(name string) error {
	key := -1
	for i, t := range r.topics {
		if t.name == name {
			key = i
		}
	}
	if key == -1 {
		return errors.New("topic not found")
	}
	r.topics[key] = r.topics[len(r.topics)-1]
	r.topics = r.topics[:len(r.topics)-1]
	return nil
}

func (r *Broker) Publish(topic string, message string) error {
	tp := &Topic{}
	for i, t := range r.topics {
		if t.name == topic {
			tp = &r.topics[i]
		}
	}

	if tp.name == "" {
		return errors.New("topic not found")
	}

	mu.Lock()
	defer mu.Unlock()

	tp.messages = append(tp.messages, message)
	return nil
}

func (r *Broker) AddConsumer(topic string) (string, error) {
	c := Consumer{name: uuid.NewString(), offset: 0}
	found := false
	for i, t := range r.topics {
		if t.name == topic {
			found = true
			r.topics[i].consumers = append(r.topics[i].consumers, c)
		}
	}
	if !found {
		return "", errors.New("topic not found")
	}
	return c.name, nil
}

func (r *Broker) NextMessage(cName string) (string, bool, error) {
	topicKey := -1
	con := &Consumer{}
	found := false
	for i, _ := range r.topics {
		for j, c := range r.topics[i].consumers {
			if c.name == cName {
				found = true
				topicKey = i
				con = &r.topics[i].consumers[j]
			}
		}
	}
	if !found {
		return "", false, errors.New("incorrect topic / consumer. Either of them might have been removed")
	}

	if con.offset >= len(r.topics[topicKey].messages) {
		return "", false, nil
	}

	msg := r.topics[topicKey].messages[con.offset]

	mu.Lock()
	defer mu.Unlock()
	con.offset += 1

	return msg, con.offset < len(r.topics[topicKey].messages), nil
}

func (r *Broker) BlockingConsume(cName string) {
	t := time.NewTicker(500 * time.Millisecond)
	for {
		select {
		case <-t.C:
			tp := &Topic{}
			co := &Consumer{}
			found := false
			for i, t := range r.topics {
				for j, c := range t.consumers {
					if c.name == cName {
						found = true
						tp = &r.topics[i]
						co = &r.topics[i].consumers[j]
					}
				}
			}
			if !found {
				break
			}
			if co.offset < len(tp.messages) {
				msg := tp.messages[co.offset]

				mu.Lock()

				co.offset += 1
				mu.Unlock()
				fmt.Printf("Msg received: %s", msg)
			}
		}
	}
}

func (r *Broker) SetOffset(cName string, val int) error {
	con := &Consumer{}
	found := false
	for i, _ := range r.topics {
		for j, c := range r.topics[i].consumers {
			if c.name == cName {
				found = true
				con = &r.topics[i].consumers[j]
			}
		}
	}
	if !found {
		return errors.New("incorrect topic / consumer. Either of them might have been removed")
	}

	mu.Lock()
	defer mu.Unlock()
	con.offset = val

	return nil
}

func (r *Broker) GetOffset(cName string) (int, error) {
	con := &Consumer{}
	found := false
	for i, _ := range r.topics {
		for j, c := range r.topics[i].consumers {
			if c.name == cName {
				found = true
				con = &r.topics[i].consumers[j]
			}
		}
	}
	if !found {
		return 0, errors.New("incorrect topic / consumer. Either of them might have been removed")
	}
	return con.offset, nil
}

func (r *Broker) TopicCount() int {
	return len(r.topics)
}

func (r *Broker) MsgsCount(topic string) int {
	for _, t := range r.topics {
		if t.name == topic {
			return len(t.messages)
		}
	}
	return 0
}

func (r *Broker) ConsumerCount(topic string) int {
	for _, t := range r.topics {
		if t.name == topic {
			return len(t.consumers)
		}
	}
	return 0
}
