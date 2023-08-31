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
		log.Println(err)
		return
	}
	err = s.Run()
	if err != nil {
		log.Fatal(err)
	}

}
