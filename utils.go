package DHTNet

import (
	"bytes"
	"encoding/binary"
	"io"
	"sync"
)


//整形转换成字节
func IntToBytes(n int) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	_ = binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}
//字节转换成整形
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var x int32
	_ =binary.Read(bytesBuffer, binary.BigEndian, &x)
	return int(x)
}

func FullCopy(a,b io.ReadWriter) chan struct{}{
	stop := make(chan struct{})
	o := &sync.Once{}

	go func() {
		io.Copy(a, b)
		o.Do(func() {
			close(stop)
		})
	}()
	go func() {
		io.Copy(b, a)
		o.Do(func() {
			close(stop)
		})
	}()
	return stop
}

func NetworkToBytes(network string) []byte {
	switch network {
	case "tcp":
		return []byte{0}
	case "tcp4":
		return []byte{1}
	case "tcp6":
		return []byte{2}
	case "udp":
		return []byte{3}
	case "udp4":
		return []byte{4}
	case "udp6":
		return []byte{5}
	default:
		return []byte{255}
	}
}
func BytesToNetwork(network []byte ) string {
	switch network[0] {
	case byte(0):
		return  "tcp"
	case byte(1):
		return  "tcp4"
	case byte(2):
		return  "tcp6"
	case byte(3):
		return  "udp"
	case byte(4):
		return  "udp4"
	case byte(5):
		return  "udp6"
	default:
		return "tcp"
	}
}
func PeerIdToDomain(peerId string) string {
	res := ""
	for _, c := range peerId{
		if c > 'Z' || c < 'A'{
			res += string(c)
		}else {
			res += "-"
			res += string(c+32)
		}
	}
	return res
}
func DomainToPeerId(domain string) string {
	res := ""
	needToTuen := false
	for _, c := range domain{
		if c == '-'{
			needToTuen = true
		}else if needToTuen {
			res += string(c-32)
			needToTuen = false
		}else {
			res += string(c)
		}
	}
	return res
}