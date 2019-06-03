package DHTNet

import (
	"github.com/libp2p/go-libp2p-core/network"
	"net"
)

type streamConn struct {
	network.Stream
}
func (s *streamConn)LocalAddr() net.Addr{
	return &net.TCPAddr{
		IP:   net.IP{127,0,0,1},
		Port: 0,
		Zone: "",
	}
}
	// RemoteAddr returns the remote network address.
func (s *streamConn)RemoteAddr() net.Addr{
	return &net.TCPAddr{
		IP:   net.IP{127,0,0,1},
		Port: 0,
		Zone: "",
	}
}
