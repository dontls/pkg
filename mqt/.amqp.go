package mqt

import (
	"github.com/streadway/amqp"
)

func init() {
	Plugins["amqp"] = NewAmqp
}

type AmqpCli struct {
	Conn     *amqp.Connection
	channel  *amqp.Channel
	exchange string
	routing  string
	*Options
}

func NewAmqp(opt *Options) (Interface, error) {
	cli := &AmqpCli{Options: opt}
	conn, err := amqp.Dial(cli.Address)
	if err != nil {
		return cli, err
	}
	cli.Conn = conn
	cli.channel, err = conn.Channel()
	if err != nil {
		return cli, err
	}
	err = cli.channel.ExchangeDeclare(cli.exchange, "topic", true, false, false, false, nil)
	return cli, err
}

func (o *AmqpCli) Publish(dest string, b []byte) error {
	if _, err := o.channel.QueueDeclare(dest, true, false, false, false, nil); err != nil {
		return err
	}
	o.channel.Publish(o.exchange, o.routing, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        b,
	})
	return nil
}

func (o *AmqpCli) Subscribe(dest string, handler func([]byte) error) error {
	if _, err := o.channel.QueueDeclare(dest, true, false, false, false, nil); err != nil {
		return err
	}
	msgChan, err := o.channel.Consume(dest, "", true, false, false, true, nil)
	if err != nil {
		return err
	}
	go func() {
		for v := range msgChan {
			if err := handler(v.Body); err == nil {
				v.Ack(true)
			}
		}
	}()
	return nil
}

func (o *AmqpCli) Release() {
	if o.Conn != nil {
		o.Conn.Close()
		o.channel.Close()
	}
}
