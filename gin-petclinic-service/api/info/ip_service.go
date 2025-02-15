package info

import (
	"go.uber.org/zap"
	"net"
)

// IPServicer
type IPServicer interface {
	lookupIP(host string) ([]net.IP, error)
}

type IPService struct {
	logger *zap.Logger
}

func NewIPService(logger *zap.Logger) *IPService {
	return &IPService{logger}
}

func (service *IPService) lookupIP(host string) ([]net.IP, error) {
	return net.LookupIP(host)
}
