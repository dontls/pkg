package tcp

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"sync/atomic"
	"time"
)

type Conn struct {
	conn      net.Conn
	localAddr string
	ioReader  io.Reader
	ioWriter  io.Writer
	recvBytes []byte
	recvBuf   bytes.Buffer
	done      int32
	client    IClient
	s         *Server
}

func (c *Conn) RemoteAddr() string {
	return c.conn.RemoteAddr().String()
}

func (c *Conn) Write(b []byte) error {
	_, err := c.conn.Write(b)
	return err
}

func (c *Conn) Close() {
	c.conn.Close()
}

func (c *Conn) Body() []byte {
	return c.recvBuf.Bytes()
}

func (c *Conn) SrvPort() uint16 {
	return c.s.port
}

// 接收并返回消息
func (c *Conn) cacheRead(t *time.Time) int {
	c.conn.SetReadDeadline(t.Add(time.Millisecond * 200)) // 读超时
	recvLen, err := c.ioReader.Read(c.recvBytes)
	if err == io.EOF {
		panic(err)
	}
	if recvLen > 0 {
		c.recvBuf.Write(c.recvBytes[:recvLen])
	}
	return recvLen
}

func (c *Conn) processData() {
	if c.client == nil {
		for _, h := range c.s.adapters {
			if c.client = h(c); c.client != nil {
				break
			}
		}
		if c.client == nil {
			panic(errors.New("can't adapta protocol"))
		}
		c.s.addConn(c)
	}
	for {
		n, err := c.client.OnHandler(c.recvBuf.Bytes())
		if err != nil {
			panic(err)
		}
		if n == 0 {
			break
		}
		c.recvBuf.Next(n)
	}
}

func (c *Conn) start() {
	defer func() {
		c.conn.Close()
		e := recover()
		fmt.Println(c.conn.RemoteAddr().String()+" closed.", e)
		if c.client != nil {
			c.s.deleteConn(c)
			c.client.OnClose(fmt.Errorf("%v", e))
		}
	}()
	fmt.Println(c.conn.RemoteAddr().String()+" connected.", c.s.port)
	c.localAddr = c.conn.LocalAddr().Network()
	rtm := time.Now()
	for {
		if atomic.LoadInt32(&c.done) > 0 {
			panic(errors.New("closed by user"))
		}
		now := time.Now()
		if n := c.cacheRead(&now); n > 0 {
			rtm = now
			c.processData()
		}
		sec := now.Sub(rtm).Seconds()
		if sec > c.s.timeout {
			panic("io/timeout by options")
		}
	}
}
