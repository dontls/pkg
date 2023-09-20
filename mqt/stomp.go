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

func (o *StompCli) dialPublish(dest string, b []byte) (err error) {
	if o.Conn == nil {
		if o.Conn, err = stomp.Dial("tcp", o.Address, o.ConnOpt...); err != nil {
			return
		}
	}
	return o.Conn.Send(dest, "text/plain", b)
}

func (o *StompCli) Publish(dest string, b []byte) error {
	if err := o.dialPublish(dest, b); err == nil {
		return nil
	}
	if o.Conn != nil {
		o.Conn.Disconnect()
		o.Conn = nil
	}
	return o.dialPublish(dest, b)
}

func (o *StompCli) dialSubscribe(dest string, opts []func(*frame.Frame) error) (s *stomp.Subscription, err error) {
	if o.Conn, err = stomp.Dial("tcp", o.Address, o.ConnOpt...); err != nil {
		return nil, err
	}
	return o.Conn.Subscribe(dest, stomp.AckClient, opts...)
}

// dest使用@分割，后面为ClientID. 如topic@ClientID
// SUBSCRIB
// destination:/topic/your-topic
// id:your-subscription-id
// activemq.subscriptionName:your-subscription-name
// ack:auto
func (o *StompCli) Subscribe(dest string, handler func([]byte) error) (err error) {
	var subOpts []func(*frame.Frame) error
	arrs := strings.Split(dest, "@")
	if len(arrs) > 1 {
		o.ConnOpt = append(o.ConnOpt, stomp.ConnOpt.Header("client-id", arrs[1]))
		subOpts = append(subOpts, stomp.SubscribeOpt.Id(arrs[1]))
		subOpts = append(subOpts, stomp.SubscribeOpt.Header("activemq.subscriptionName", arrs[1]))
	}
	go func() {
		for o.Address != "" {
			s, err := o.dialSubscribe(arrs[0], subOpts)
			log.Println(dest, err)
			if err != nil {
				time.Sleep(2 * time.Second)
				continue
			}
			for v := range s.C {
				if v.Err != nil {
					break
				}
				if err := handler(v.Body); err == nil {
					o.Conn.Ack(v)
				}
			}
		}
	}()
	return nil
}

func (o *StompCli) Release() {
	if o.Conn != nil {
		o.Address = ""
		o.Conn.Disconnect()
	}
}
