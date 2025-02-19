package svc

import (
	"log"
	"os"

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

func ErrorOutput(filename string) {
	// go1.23 debug.SetCrashOutput(f, debug.CrashOptions{})
	if file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); err == nil {
		//50K文件
		if fi, err := file.Stat(); err == nil && fi.Size() > 51200 {
			file.Truncate(0)
		}
		crashDup(file)
	}
}

func Run(p *Program) {
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
		log.Println(app, os.Args[1], err)
		return
	}
	err = s.Run()
	if err != nil {
		log.Fatal(err)
	}
}
