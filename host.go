package DHTNet

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-kad-dht"
	"io"
	"log"
	mrand "math/rand"
	dsync "github.com/ipfs/go-datastore/sync"
	ds "github.com/ipfs/go-datastore"
	rhost "github.com/libp2p/go-libp2p/p2p/host/routed"
	ma "github.com/multiformats/go-multiaddr"
)

func makeRoutedHost(listenPort int, randseed int64, bootstrapPeers []peer.AddrInfo) (host.Host, error) {

	// If the seed is zero, use real cryptographic randomness. Otherwise, use a
	// deterministic randomness source to make generated keys stay the same
	// across multiple runs
	var r io.Reader
	if randseed == 0 {
		r = rand.Reader
	} else {
		r = mrand.New(mrand.NewSource(randseed))
	}
	//externalIP := config.Config.GetP2PExternalIP()
	//var extMultiAddr ma.Multiaddr
	//if externalIP == "" {
	//	log.Println("External IP not defined, Peers might not be able to resolve this node if behind NAT\n")
	//} else {
	//	extMultiAddr, err := ma.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", externalIP, listenPort))
	//	if err != nil {
	//		log.Printf("Error creating multiaddress: %v\n", err)
	//		return nil, err
	//	}
	//}
	//addressFactory := func(addrs []ma.Multiaddr) []ma.Multiaddr {
	//	if extMultiAddr != nil {
	//		addrs = append(addrs, extMultiAddr)
	//	}
	//	return addrs
	//}
	// Generate a key pair for this host. We will use it at least
	// to obtain a valid host ID.
	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		return nil, err
	}
	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", listenPort)),
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip6/::/tcp/%d", listenPort)),
		libp2p.Identity(priv),
		libp2p.DefaultTransports,
		libp2p.DefaultMuxers,
		libp2p.DefaultSecurity,
		//libp2p.DefaultEnableRelay,
		//libp2p.EnableRelay(),
		//libp2p.EnableAutoRelay(),
		libp2p.NATPortMap(),
		//libp2p.AddrsFactory(addressFactory),
	}

	ctx := context.Background()

	basicHost, err := libp2p.New(ctx, opts...)
	if err != nil {
		return nil, err
	}

	// Construct a datastore (needed by the DHT). This is just a simple, in-memory thread-safe datastore.
	dstore := dsync.MutexWrap(ds.NewMapDatastore())

	// Make the DHT
	dht := dht.NewDHT(ctx, basicHost, dstore)

	// Make the routed host
	routedHost := rhost.Wrap(basicHost, dht)

	// connect to the chosen ipfs nodes
	err = bootstrapConnect(ctx, routedHost, bootstrapPeers)
	if err != nil {
		return nil, err
	}

	// Bootstrap the host
	err = dht.Bootstrap(ctx)
	if err != nil {
		return nil, err
	}

	// Build host multiaddress
	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", routedHost.ID().Pretty()))

	// Now we can build a full multiaddress to reach this host
	// by encapsulating both addresses:
	// addr := routedHost.Addrs()[0]
	addrs := routedHost.Addrs()
	log.Println("I can be reached at:")
	for _, addr := range addrs {
		log.Println(addr.Encapsulate(hostAddr))
	}
	log.Printf("Now I am on %s", PeerIdToDomain(routedHost.ID().Pretty()))
	//log.Printf("Now run \"./AirNet -l %d -d %s\" on a different terminal\n", listenPort+1, routedHost.ID().Pretty())

	return routedHost, nil
}

