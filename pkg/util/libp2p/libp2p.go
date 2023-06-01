package libp2p

import (
	"fmt"
	mrand "math/rand"
	"net"

	es "github.com/go-errors/errors"
	"github.com/libp2p/go-libp2p"
	libp2pcrypto "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"

	t "github.com/filecoin-project/mir/pkg/types"
)

// NewDummyHostWithPrivKey creates new dummy libp2p host with an identity
// determined by the private key given as an input.
func NewDummyHostWithPrivKey(listenAddr t.NodeAddress, privKey libp2pcrypto.PrivKey) (host.Host, error) {

	libp2pPeerID, err := peer.AddrInfoFromP2pAddr(listenAddr)
	if err != nil {
		return nil, es.Errorf("failed to get own libp2p addr info: %w", err)
	}
	return libp2p.New(
		libp2p.Identity(privKey),
		libp2p.DefaultTransports,
		libp2p.ListenAddrs(libp2pPeerID.Addrs[0]),
	)
}

// NewDummyHost creates an insecure libp2p host for test and demonstration purposes.
func NewDummyHost(numericID int, sourceAddr t.NodeAddress) host.Host {
	// sourceMultiAddr, priv := NewDummyHostID(id, basePort)

	h, err := libp2p.New(
		// Use the keypair we generated
		libp2p.Identity(NewDummyHostKey(numericID)),
		// Multiple listen addresses
		libp2p.DefaultTransports,
		libp2p.ListenAddrs(sourceAddr),
	)
	if err != nil {
		panic(err)
	}

	return h
}

// NewDummyHostNoListen creates an insecure libp2p host for test and demonstration purposes without listening interface.
func NewDummyHostNoListen(numericID int, _ t.NodeAddress) host.Host {
	// sourceMultiAddr, priv := NewDummyHostID(id, basePort)

	h, err := libp2p.New(
		// Use the keypair we generated
		libp2p.Identity(NewDummyHostKey(numericID)),
		// Multiple listen addresses
		libp2p.DefaultTransports,
	)
	if err != nil {
		panic(err)
	}

	return h
}

func NewFreeHostAddr() (m multiaddr.Multiaddr) {
	var a *net.TCPAddr
	var err error
	a, err = net.ResolveTCPAddr("tcp", "localhost:0")
	if err == nil {
		var l *net.TCPListener
		if l, err = net.ListenTCP("tcp", a); err == nil {
			defer func() {
				if err := l.Close(); err != nil {
					panic(err)
				}
			}()
			addr, ok := l.Addr().(*net.TCPAddr)
			if !ok {
				panic("TCPAddr type assertion failed")
			}
			m, err = multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", addr.Port))
			if err != nil {
				panic(err)
			}
		}
	}
	return
}

// NewDummyHostAddr generates a libp2p host address for test and demonstration purposes.
func NewDummyHostAddr(id, basePort int) multiaddr.Multiaddr {
	sourceMultiAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", basePort+id))
	if err != nil {
		panic(err)
	}

	return sourceMultiAddr
}

func NewDummyHostKey(numericID int) libp2pcrypto.PrivKey {
	rand := mrand.New(mrand.NewSource(int64(numericID))) // nolint

	priv, _, err := libp2pcrypto.GenerateKeyPairWithReader(libp2pcrypto.Ed25519, -1, rand)
	if err != nil {
		panic(err)
	}

	return priv
}

// NewDummyMultiaddr generates a libp2p peer multiaddress for the host generated by NewDummyHost(id, basePort).
func NewDummyMultiaddr(numericID int, sourceAddr t.NodeAddress) multiaddr.Multiaddr {

	peerID, err := peer.IDFromPrivateKey(NewDummyHostKey(numericID))
	if err != nil {
		panic(err)
	}

	peerInfo := peer.AddrInfo{
		ID:    peerID,
		Addrs: []multiaddr.Multiaddr{sourceAddr},
	}
	addrs, err := peer.AddrInfoToP2pAddrs(&peerInfo)
	if err != nil {
		panic(err)
	}

	if len(addrs) != 1 {
		panic(es.Errorf("wrong number of addresses %d", len(addrs)))
	}

	return addrs[0]
}

// NewFakeMultiaddr generates a libp2p peer multiaddress for a non existed host.
func NewFakeMultiaddr(numericID, port int) multiaddr.Multiaddr {

	sourceAddr := NewDummyHostAddr(0, port)

	peerID, err := peer.IDFromPrivateKey(NewDummyHostKey(numericID))
	if err != nil {
		panic(err)
	}

	peerInfo := peer.AddrInfo{
		ID:    peerID,
		Addrs: []multiaddr.Multiaddr{sourceAddr},
	}
	addrs, err := peer.AddrInfoToP2pAddrs(&peerInfo)
	if err != nil {
		panic(err)
	}

	if len(addrs) != 1 {
		panic(es.Errorf("wrong number of addresses %d", len(addrs)))
	}

	return addrs[0]
}
