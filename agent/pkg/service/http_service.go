package service

import (
	"errors"
)

type HttpServiceConfiguration struct {
	Port int32
}
type HttpService struct {
	HttpServiceConfiguration *HttpServiceConfiguration
}

func (h *HttpService) Serve() error {
	if h.HttpServiceConfiguration == nil {
		return errors.New("http service configuration has not been initialised")
	}

	return nil
}
