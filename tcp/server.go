package tcp

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

type IClient interface {
	OnClose(error)
	OnHandler([]byte) (int, error)
	Request(interface{}) error
}

type AdapterHandler func(*Conn) IClient

type Server struct {
	listener net.Listener
	wg       sync.WaitGroup
	err      error
	adapters []AdapterHandler
	port     uint16
	timeout  float64
	mpClis   map[string]*Conn
	lock     sync.RWMutex
}

// ListenTCPAndServe start server
func NewServer(port uint16) *Server {
	s := &Server{port: port, timeout: 60}
	tcpAddr, _ := net.ResolveTCPAddr("tcp4", fmt.Sprintf(":%d", port)) //获取一个tcpAddr
	s.listener, s.err = net.ListenTCP("tcp", tcpAddr)
	return s
}

func (s *Server) SetReadTimeout(timeout float64) {
	s.timeout = timeout
}

func (s *Server) addConn(c *Conn) {
	s.lock.Lock()
	s.mpClis[c.remoteAddr] = c
	s.lock.Unlock()
}

func (s *Server) deleteConn(c *Conn) {
	s.lock.Lock()
	delete(s.mpClis, c.remoteAddr)
	s.lock.Unlock()
}

func (s *Server) accept() {
	defer s.listener.Close()
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Println(err)
			break
		}
		if c := s.newConn(conn); c != nil {
			c.s.addConn(c)
			s.wg.Add(1)
			go func() {
				c.start()
				c.s.deleteConn(c)
				s.wg.Done()
			}()
		}
	}
}

func (s *Server) newConn(c net.Conn) *Conn {
	return &Conn{
		conn:       c,
		ioReader:   io.Reader(c),
		ioWriter:   io.Writer(c),
		recvBytes:  make([]byte, 2048),
		s:          s,
		remoteAddr: c.RemoteAddr().String(),
	}
}

func (s *Server) ListenTCP(handler ...AdapterHandler) error {
	if s.err != nil {
		return s.err
	}
	s.adapters = append(s.adapters, handler...)
	go s.accept()
	return nil
}

func (s *Server) Shutdown() error {
	if s.listener != nil {
		s.listener.Close()
	}
	s.lock.RLock()
	for _, c := range s.mpClis {
		c.Close()
	}
	s.lock.RLock()
	s.wg.Wait()
	return nil
}
