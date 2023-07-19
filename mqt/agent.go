package mqt

import (
	"fmt"
)

type Interface interface {
	Publish(string, []byte) error
	Subscribe(string, func([]byte) error) error
	Release()
}

type Options struct {
	Address string
	User    string
	Pswd    string
}

type Client struct {
	clients []Interface
	Options
	goi int // 当前使用的goroutine
}

func (o *Client) goIndex() int {
	o.goi++
	if o.goi >= len(o.clients) {
		o.goi = 0
	}
	return o.goi
}

func (o *Client) Release() {
	for _, v := range o.clients {
		if v != nil {
			v.Release()
		}
	}
}

// 订阅会自动分配connection对象
// 订阅数大于连接数，出现同一连接多次订阅，报错
func (o *Client) Subscribe(dest string, handler func([]byte) error) error {
	i := o.goIndex()
	return o.clients[i].Subscribe(dest, handler)
}

// 动态均衡，自动适配connection发送数据
func (o *Client) Publish(dest string, b []byte) error {
	i := o.goIndex()
	return o.clients[i].Publish(dest, b)
}

type newHandler func(*Options) (Interface, error)

var Plugins = map[string]newHandler{}

// 支持创建多goroutine发布
func News(name string, opt *Options, count int) (*Client, error) {
	handler, ok := Plugins[name]
	if !ok {
		return nil, fmt.Errorf("unsupported mq %s", name)
	}
	if count == 0 {
		count = 1
	}
	c := &Client{Options: *opt}
	for i := 0; i < count; i++ {
		cli, err := handler(opt)
		if err != nil {
			return nil, err
		}
		c.clients = append(c.clients, cli)
	}
	return c, nil
}
