package DHTNet

import (
	"bufio"
	"fmt"
	"github.com/armon/go-socks5"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"golang.org/x/net/context"
	"log"
	"net"
	"strconv"
	"strings"
)


const Protocol = "/DHTNet/0.0.1"

func New(p2pport *int, seed *int64, socksPort *int){

	var bootstrapPeers []peer.AddrInfo
	log.Println("使用全球DHT节点")
	bootstrapPeers = IPFS_PEERS
	host, err := makeRoutedHost(*p2pport, *seed, bootstrapPeers)
	if err != nil {
		log.Fatal(err)
	}
	host.SetStreamHandler(Protocol,streamHandler)

	// Create a SOCKS5 server
	conf := &socks5.Config{
		Resolver: &DNSResolver{},
		Dial: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
			domain := strings.Split(addr, ":")[0]
			port, err := strconv.Atoi(strings.Split(addr, ":")[1])
			if err != nil {
				return nil, err
			}
			domainlist := strings.Split(domain, ".")
			if domain[len(domain)-1] == '.' {
				domainlist = domainlist[0 : len(domainlist)-1]
			}
			if domainlist[len(domainlist)-1] != "qm" {
				return net.Dial(network, addr)
			}
			destPeer := domainlist[len(domainlist)-2]
			destPeer = DomainToPeerId(destPeer)
			destPeerID, err := peer.IDB58Decode(destPeer)
			stream, err := host.NewStream(context.Background(), destPeerID, Protocol)
			if err != nil {
				log.Printf("[ERR] 无法连接到目标地址: %v", err)
				return nil, err
			}
			defer stream.Close()
			if _, err = stream.Write(IntToBytes(port)); err != nil {
				log.Printf("[ERR] 无法连接到目标地址的目标端口: %v", err)
				return nil, err
			}
			if _, err = stream.Write(NetworkToBytes(network)); err != nil {
				log.Printf("[ERR] 无法连接到目标地址的目标端口: %v", err)
				return nil, err
			}

			return &streamConn{stream}, nil
		},
	}
	server, err := socks5.New(conf)
	if err != nil {
		panic(err)
	}

	// Create SOCKS5 proxy on localhost port 8000
	if err := server.ListenAndServe("tcp", fmt.Sprintf("127.0.0.1:%d", *socksPort)); err != nil {
		panic(err)
	}
}

func streamHandler(stream network.Stream) {
	defer stream.Close()
	buf := bufio.NewReader(stream)
	portByte := []byte{0, 0, 0, 0}
	if _, err := buf.Read(portByte); err != nil {
		log.Printf("[ERR] 读取端口字节失败: %v", err)
		return
	}
	port := BytesToInt(portByte)
	if port > 65535 || port < 0 {
		err := fmt.Errorf("无效端口号: %d", port)
		log.Printf("[ERR] 转发: %v", err)
		return
	}
	networkBuf := []byte{0}
	if _, err := buf.Read(networkBuf); err != nil {
		log.Printf("[ERR] 读取连接方式失败: %v", err)
		return
	}
	c, err := net.Dial(BytesToNetwork(networkBuf), fmt.Sprintf("127.0.0.1:%d",port))
	if err != nil {
		log.Printf("[ERR] 尝试连接错误端口: %d", port)
		return
	}
	<- FullCopy(c, stream)
	return
}