package info

import (
	"net"
)

// IPServicer
type IPServicer interface {
	lookupIP(host string) ([]net.IP, error)
}

type IPService struct {
}

func NewIPService() *IPService {
	return &IPService{}
}

func (service *IPService) lookupIP(host string) ([]net.IP, error) {
	return net.LookupIP(host)
}
