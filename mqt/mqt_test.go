package mqt

import (
	"encoding/json"
	"log"
	"testing"
	"time"
)

var index = 0

func getPushData() []byte {
	index++
	v := map[string]any{
		"device": "20198002",
		"now":    time.Now().Format("2006-01-02 15:04:05"),
		"index":  index,
	}
	data, _ := json.Marshal(v)
	return data
}

func TestMqtt(t *testing.T) {
	c, err := News("mqtt", &Options{Address: "127.0.0.1:35003"}, 3)
	if err != nil {
		log.Fatalln(err)
	}
	defer c.Release()
	topic := "test/mqtt"
	c.Subscribe(topic, func(b []byte) error {
		log.Printf("%s: %s\n", topic, b)
		return nil
	})
	for {
		time.Sleep(2 * time.Second)
		c.Publish(topic, getPushData())
	}
}

// Fixme 区分publisher, subscriber
func TestStomp(t *testing.T) {
	c, err := News("stomp", &Options{Address: "127.0.0.1:35002"}, 3)
	if err != nil {
		log.Fatalln(err)
	}
	defer c.Release()
	dest := "/topic/test/stomp"
	c.Subscribe(dest+"@client-1", func(b []byte) error {
		log.Printf("%s:1 %s\n", dest, b)
		return nil
	})
	c.Subscribe(dest, func(b []byte) error {
		log.Printf("%s:2 %s\n", dest, b)
		return nil
	})
	for {
		time.Sleep(2 * time.Second)
		c.Publish(dest, getPushData())
	}
}

func TestNats(t *testing.T) {
	c, err := News("nats", &Options{Address: "nats://127.0.0.1:4222"}, 3)
	if err != nil {
		log.Fatalln(err)
	}
	defer c.Release()
	dest := "/test/nats"
	c.Subscribe(dest, func(b []byte) error {
		log.Printf("%s:1 %s\n", dest, b)
		return nil
	})
	c.Subscribe(dest, func(b []byte) error {
		log.Printf("%s:2 %s\n", dest, b)
		return nil
	})
	for {
		time.Sleep(2 * time.Second)
		c.Publish(dest, getPushData())
	}
}

func TestKafka(t *testing.T) {
	c, err := News("kafka", &Options{Address: "172.16.60.219:9092"}, 1)
	if err != nil {
		log.Fatalln(err)
	}
	defer c.Release()
	dest := "test-topic"
	c.Subscribe(dest, func(b []byte) error {
		log.Printf("%s:1 %s\n", dest, b)
		return nil
	})
	for {
		time.Sleep(2 * time.Second)
		log.Printf("Publish %v\n", c.Publish(dest, getPushData()))
	}
}
