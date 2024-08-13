package service

import "github.com/kardianos/service"

// ServiceWrapper is an interface that wraps the methods of the kardianos/service.Service.
type ServiceWrapper interface {
	Start() error
	Stop() error
	Install() error
	Uninstall() error
	Status() (service.Status, error)
	Run() error
	Logger(errs chan<- error) (service.Logger, error)
	Platform() string
	Restart() error
	String() string
	SystemLogger(errs chan<- error) (service.Logger, error)
}
