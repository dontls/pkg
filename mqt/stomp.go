package mqt

import (
	"log"
	"net"
	"strings"
	"time"

	"github.com/go-stomp/stomp/v3"
	"github.com/go-stomp/stomp/v3/frame"
)

func init() {
	Plugins["stomp"] = NewStomp
}

type StompCli struct {
	Conn    *stomp.Conn
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
	conn, err := net.Dial("tcp", opt.Address)
	if err != nil {
		return cli, err
	}
	defer conn.Close()
	cli.ConnOpt = append(cli.ConnOpt, stomp.ConnOpt.HeartBeat(10*time.Second, 10*time.Second))
	// 保障链接不read timeout
	cli.ConnOpt = append(cli.ConnOpt, stomp.ConnOpt.HeartBeatError(360*time.Second))
	if opt.User != "" {
		cli.ConnOpt = append(cli.ConnOpt, stomp.ConnOpt.Login(opt.User, opt.Pswd))
	}
	return cli, nil
}

func (o *StompCli) Publish(dest string, b []byte) (err error) {
	if o.Conn == nil {
		if o.Conn, err = stomp.Dial("tcp", o.Address, o.ConnOpt...); err != nil {
			return
		}
	}
	if err = o.Conn.Send(dest, "text/plain", b); err == nil {
		return
	}
	o.Release()
	if o.Conn, err = stomp.Dial("tcp", o.Address, o.ConnOpt...); err != nil {
		return
	}
	return o.Conn.Send(dest, "text/plain", b)
}

// dest使用@分割，后面为ClientID. 如topic@ClientID
// SUBSCRIB
// destination:/topic/your-topic
// id:your-subscription-id
// activemq.subscriptionName:your-subscription-name
// ack:auto
func (o *StompCli) Subscribe(dest string, handler func([]byte) error) (err error) {
	opts := o.ConnOpt
	var subOpts []func(*frame.Frame) error
	arrs := strings.Split(dest, "@")
	if len(arrs) > 1 {
		opts = append(opts, stomp.ConnOpt.Header("client-id", arrs[1]))
		subOpts = append(subOpts, stomp.SubscribeOpt.Id(arrs[1]))
		subOpts = append(subOpts, stomp.SubscribeOpt.Header("activemq.subscriptionName", arrs[1]))
	}
	if o.Conn, err = stomp.Dial("tcp", o.Address, opts...); err != nil {
		return err
	}
	s, err := o.Conn.Subscribe(arrs[0], stomp.AckClient, subOpts...)
	if err != nil {
		return err
	}
	go func() {
		for v := range s.C {
			if v.Err != nil {
				log.Println(dest, err)
				break
			}
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
