package p2p

import (
	"fmt"
	"net"
	"sync"
)

// TCPeer represents the remote node over a TCP established connection
type TCPeer struct {
	// conn is underlying connection of the peer
	conn net.Conn

	// if we accept and retrieve a conn => outbound == true
	outbound bool
}

func NewTCPeer(conn net.Conn, outbound bool) *TCPeer {
	return &TCPeer{
		conn:     conn,
		outbound: outbound,
	}
}

type TCPTransport struct {
	listenAddress string
	listener      net.Listener

	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

func NewTCPTransport(listenAddr string) *TCPTransport {
	return &TCPTransport{
		listenAddress: listenAddr,
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error

	t.listener, err = net.Listen("tcp", t.listenAddress)
	if err != nil {
		return err
	}
	go t.startAcceptLoop()

	return nil

}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Println("TCP accept error: %v\n", err)
		}
		go t.handleConn(conn)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn) {
	peer := NewTCPeer(conn, true)
	fmt.Printf("New incoming connection: %+v\n", peer)
}
