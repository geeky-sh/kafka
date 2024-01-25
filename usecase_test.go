package flipkart

import (
	"fmt"
	"log"
	"testing"
)

func TestBroker(t *testing.T) {
	b := NewBroker()

	// test add topic
	b.AddTopic("t1")
	b.AddTopic("t2")
	b.AddTopic("t3")
	b.AddTopic("t4")

	want := 4
	c := b.TopicCount()
	if c != want {
		t.Errorf("Got %d Want %d\n", c, want)
	}

	// test delete topic
	b.DeleteTopic("t1")

	want = 3
	c = b.TopicCount()
	if c != want {
		t.Errorf("Got %d Want %d\n", c, want)
	}

	b.DeleteTopic("t4")

	want = 2
	c = b.TopicCount()
	if c != want {
		t.Errorf("Got %d Want %d\n", c, want)
	}

	// simulate multiple publishers
	// var wg sync.WaitGroup
	// for i := 0; i < 20; i++ {
	// 	wg.Add(1)

	// 	msg := fmt.Sprintf("msg%d", i+1)
	// 	go func() {
	// 		if err := b.Publish("t2", msg); err != nil {
	// 			log.Fatalf("Err %v", err)
	// 		}
	// 		wg.Done()
	// 	}()
	// }
	// wg.Wait()

	// test publishing messages
	for i := 0; i < 20; i++ {
		msg := fmt.Sprintf("msg%d", i+1)
		if err := b.Publish("t2", msg); err != nil {
			log.Fatalf("Err %v", err)
		}
	}

	want = 20
	got := b.MsgsCount("t2")
	if got != want {
		t.Errorf("Got %d Want %d\n", got, want)
	}

	// test adding consumer
	cname, err := b.AddConsumer("t2")
	if err != nil {
		log.Fatalf("Err %v", err)
	}

	want = 1
	got = b.ConsumerCount("t2")
	if got != want {
		t.Errorf("Got %d Want %d\n", got, want)
	}

	// test message consumption
	gots, hasMore, err := b.NextMessage(cname)
	if err != nil {
		log.Fatalf("Err %v", err)
	}
	wants := "msg1"
	if gots != wants {
		t.Errorf("Got %s Want %s\n", gots, wants)
	}
	if !hasMore {
		t.Error("Should have more messages", got, want)
	}

	gots, hasMore, err = b.NextMessage(cname)
	if err != nil {
		log.Fatalf("Err %v", err)
	}
	wants = "msg2"
	if gots != wants {
		t.Errorf("Got %s Want %s\n", gots, wants)
	}
	if !hasMore {
		t.Error("Should have more messages", got, want)
	}

	want = 2 // since two messages are consumed till now
	got, err = b.GetOffset(cname)
	if got != want {
		t.Errorf("Got %d Want %d\n", got, want)
	}

	b.SetOffset(cname, 0)

	gots, hasMore, err = b.NextMessage(cname)
	if err != nil {
		log.Fatalf("Err %v", err)
	}
	wants = "msg1"
	if gots != wants {
		t.Errorf("Got %s Want %s\n", gots, wants)
	}
	if !hasMore {
		t.Error("Should have more messages", got, want)
	}

	// b.BlockingConsume(cname) // run this if you want to simulation blocking call which will consume all the items
}
