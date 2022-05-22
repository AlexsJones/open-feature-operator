package service

type IServiceConfiguration interface {
}
type IService interface {
	Serve() error
}
