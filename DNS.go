package DHTNet

import (
	"context"
	"net"
	"strings"
)

type DNSResolver struct{}

func (d DNSResolver) Resolve(ctx context.Context, name string) (context.Context, net.IP, error) {
	domainlist := strings.Split(name, ".")
	if domainlist[len(domainlist)-1] == "qm" || domainlist[len(domainlist)-2] == "qm" {
		return ctx, net.IP{0,0,0,0}, nil
	}
	addr, err := net.ResolveIPAddr("ip", name)
	if err != nil {
		return ctx, nil, err
	}
	return ctx, addr.IP, err
}

