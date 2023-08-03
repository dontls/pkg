package mqt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func init() {
	Plugins["mqtt"] = NewMqtt
}

type MqttCli struct {
	Conn    mqtt.Client
	connOpt *mqtt.ClientOptions
	*Options
}

func (o *MqttCli) connectServe() error {
	o.Conn = mqtt.NewClient(o.connOpt)
	token := o.Conn.Connect()
	token.Wait()
	return token.Error()
}

func NewMqtt(opt *Options) (Interface, error) {
	cli := &MqttCli{Options: opt}
	cli.connOpt = mqtt.NewClientOptions()
	cli.connOpt.AddBroker(opt.Address)
	if opt.User != "" {
		cli.connOpt.SetUsername(opt.User)
		cli.connOpt.SetPassword(opt.Pswd)
	}
	if err := cli.connectServe(); err != nil {
		return cli, err
	}
	return cli, nil
}

func (o *MqttCli) Publish(topic string, b []byte) error {
	if !o.Conn.IsConnected() {
		if err := o.connectServe(); err != nil {
			return err
		}
	}
	return o.Conn.Publish("topic/"+topic, 0, false, b).Error()
}

func (o *MqttCli) Subscribe(dest string, handler func([]byte) error) error {
	o.Conn.Subscribe("topic/"+dest, 1, func(c mqtt.Client, m mqtt.Message) {
		if err := handler(m.Payload()); err == nil {
			m.Ack()
		}
	})
	return nil
}

func (o *MqttCli) Release() {
	if o.Conn != nil {
		o.Conn.Disconnect(1)
	}
}
