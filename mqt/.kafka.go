package mqt

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

func init() {
	Plugins["kafka"] = NewKafka
}

type KafkaCli struct {
	Conn *kafka.Conn
	*Options
}

// 127.0.0.1:9092, fix：只支持域名
// notice: kafka一个topic对应一个conn
func NewKafka(opt *Options) (Interface, error) {
	conn, err := kafka.Dial("tcp", opt.Address)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	parts, err := conn.ReadPartitions()
	if err != nil {
		return nil, err
	}
	for _, v := range parts {
		fmt.Printf("%s PartitionID %d\n", v.Topic, v.ID)
	}
	c := &KafkaCli{Options: opt}
	return c, nil
}

func (o *KafkaCli) dialLeader(topic string) error {
	if o.Conn != nil {
		return nil
	}
	// o.Writer = &kafka.Writer{
	// 	Addr: kafka.TCP(o.Address),
	// 	Topic: topic,
	// 	Balancer: &kafka.LeastBytes{},
	// }
	conn, err := kafka.DialLeader(context.TODO(), "tcp", o.Address, topic, 0)
	if err == nil {
		o.Conn = conn
	}
	return nil
}

func (o *KafkaCli) Publish(topic string, b []byte) error {
	if err := o.dialLeader(topic); err != nil {
		return err
	}
	_, err := o.Conn.WriteMessages(kafka.Message{Value: b})
	return err

}

func (o *KafkaCli) Subscribe(dest string, handler func([]byte) error) error {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{o.Address},
		Topic:     dest,
		Partition: 0,
	})
	go func() {
		defer r.Close()
		ctx := context.Background()
		for {
			m, err := r.FetchMessage(ctx)
			if err != nil {
				break
			}
			// fmt.Printf("offset %d %s=%s\n", m.Offset, m.Key, m.Value)
			if err := handler(m.Value); err == nil {
				r.CommitMessages(ctx, m)
			}
		}
	}()
	return nil
}

func (o *KafkaCli) Release() {
	if o.Conn != nil {
		o.Conn.Close()
		o.Conn = nil
	}
}
