package svc

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
	"unsafe"

	"github.com/kardianos/service"
)

type Program struct {
	AppName     func() string
	Description string
	Run         func() error
	Shutdown    func() error
}

func (p *Program) Start(s service.Service) error {
	return p.Run()
}

func (p *Program) Stop(s service.Service) error {
	return p.Shutdown()
}

var tracefile = filepath.Dir(os.Args[0]) + "/" + time.Now().Format("20060102150405") + ".dump"

func Run(p *Program) {
	defer func() {
		logFile, err := os.OpenFile(tracefile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
		if err == nil {
			log.SetOutput(logFile) // 将文件设置为log输出的文件
			log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)
		}
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			log.Fatalf("runtime error: %v\ntraceback:\n%v\n", err, *(*string)(unsafe.Pointer(&buf)))
		}
	}()
	app := p.AppName()
	s, err := service.New(p, &service.Config{
		Name:        app,
		DisplayName: app,
		Description: p.Description,
	})
	if err != nil {
		log.Fatal(err)
	}
	if len(os.Args) > 1 {
		if os.Args[1] == "install" {
			err = s.Install()
		} else if os.Args[1] == "uninstall" {
			err = s.Uninstall()
		}
		log.Println(err)
		return
	}
	err = s.Run()
	if err != nil {
		log.Fatal(err)
	}
}
