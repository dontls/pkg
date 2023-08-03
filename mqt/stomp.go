package mqt

import (
	"net"
	"strings"
	"time"

	"github.com/go-stomp/stomp/v3"
)

func init() {
	Plugins["stomp"] = NewStomp
}

type StompCli struct {
	Conn    *stomp.Conn
	conn    net.Conn
	ConnOpt []func(*stomp.Conn) error
	*Options
}

// CONNECT
// login:username
// passcode:password
// host:/my-vhost
// client-id:your-client-id
// heart-beat:10000,20000

func NewStomp(opt *Options) (Interface, error) {
	cli := &StompCli{Options: opt}
	cli.ConnOpt = append(cli.ConnOpt, stomp.ConnOpt.HeartBeat(60*time.Second, 60*time.Second))
	if opt.User != "" {
		cli.ConnOpt = append(cli.ConnOpt, stomp.ConnOpt.Login(opt.User, opt.Pswd))
	}
	conn, err := net.Dial("tcp", opt.Address)
	if err == nil {
		cli.conn = conn
	}
	return cli, err
}

func (o *StompCli) connect(opts ...func(*stomp.Conn) error) (err error) {
	opts = append(o.ConnOpt, opts...)
	o.Conn, err = stomp.Connect(o.conn, opts...)
	return err
}

func (o *StompCli) Publish(dest string, b []byte) (err error) {
	if o.Conn == nil {
		if err := o.connect(); err != nil {
			return err
		}
	}
	if err := o.Conn.Send(dest, "text/plain", b); err == nil {
		return err
	}
	o.Release()
	o.Conn, err = stomp.Dial("tcp", o.Address, o.ConnOpt...)
	if err != nil {
		return err
	}
	return o.Conn.Send(dest, "text/plain", b)
}

// dest使用@分割，前面为id. 如topic@ClientID
// SUBSCRIB
// destination:/topic/your-topic
// id:your-subscription-id
// activemq.subscriptionName:your-subscription-name
// ack:auto
func (o *StompCli) Subscribe(dest string, handler func([]byte) error) (err error) {
	arrs := strings.Split(dest, "@")
	var s *stomp.Subscription
	if len(arrs) > 1 {
		if err = o.connect(stomp.ConnOpt.Header("client-id", arrs[1])); err != nil {
			return err
		}
		s, err = o.Conn.Subscribe(arrs[0], stomp.AckClient, stomp.SubscribeOpt.Id(arrs[1]), stomp.SubscribeOpt.Header("activemq.subscriptionName", arrs[1]))
	} else {
		if err := o.connect(); err != nil {
			return err
		}
		s, err = o.Conn.Subscribe(arrs[0], stomp.AckClient)
	}
	if err != nil {
		return err
	}
	go func() {
		for v := range s.C {
			if err := handler(v.Body); err == nil {
				o.Conn.Ack(v)
			}
		}
	}()
	return nil
}

func (o *StompCli) Release() {
	if o.Conn != nil {
		o.Conn.Disconnect()
		o.Conn = nil
	}
}
