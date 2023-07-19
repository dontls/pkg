package mqt

import (
	"time"

	"github.com/go-stomp/stomp/v3"
)

func init() {
	Plugins["stomp"] = NewStomp
}

type StompCli struct {
	Conn    *stomp.Conn
	ConnOpt []func(*stomp.Conn) error
	*Option
}

func NewStomp(opt *Options) (Interface, error) {
	cli := &StompCli{Options: opt}
	cli.ConnOpt = append(cli.ConnOpt, stomp.ConnOpt.HeartBeat(60*time.Second, 60*time.Second))
	if opt.User != "" {
		cli.ConnOpt = append(cli.ConnOpt, stomp.ConnOpt.Login(opt.User, opt.Pswd))
	}
	conn, err := stomp.Dial("tcp", opt.Address, cli.ConnOpt...)
	if err == nil {
		cli.Conn = conn
	}
	return cli, nil
}

func (o *StompCli) Publish(dest string, b []byte) error {
	err := o.Conn.Send(dest, "text/plain", b)
	if err == nil {
		return err
	}
	o.Conn, err = stomp.Dial("tcp", o.Address, o.ConnOpt...)
	if err != nil {
		return err
	}
	return o.Conn.Send(dest, "text/plain", b)
}

func (o *StompCli) Subscribe(dest string, handler func([]byte) error) error {
	s, err := o.Conn.Subscribe(dest, stomp.AckClient)
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
	}
}
