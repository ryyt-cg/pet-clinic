package info

import (
	"net"
)

// IPServicer defines the interface for looking up IP addresses by host name.
type IPServicer interface {
	LookupIP(host string) ([]net.IP, error)
}

type IPService struct {
}

func NewIPService() *IPService {
	return &IPService{}
}

func (service *IPService) LookupIP(host string) ([]net.IP, error) {
	return net.LookupIP(host)
}
