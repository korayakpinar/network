package client

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/korayakpinar/network/src/handler"
	"github.com/korayakpinar/network/src/proxy"
	"github.com/korayakpinar/network/src/utils"
	kaddht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
)

type Client struct {
	Host      host.Host
	PubSub    *pubsub.PubSub
	DHT       *kaddht.IpfsDHT
	Handler   *handler.Handler
	Proxy     *proxy.Proxy
	operators []peer.ID
	cfg       utils.Config
}

type discoveryNotifee struct {
	h   host.Host
	ctx context.Context
}

func NewClient(h host.Host, dht *kaddht.IpfsDHT, cfg utils.Config) *Client {
	return &Client{h, nil, dht, nil, nil, make([]peer.ID, 0), cfg}
}

func (cli *Client) Start(ctx context.Context, topicName string) {
	var topicHandle *pubsub.Topic
	errChan := make(chan error, 4)

	// Start the discovery
	go cli.startDiscovery(ctx, topicName, errChan)

	// Start the pubsub
	topicHandle = cli.startPubsub(ctx, topicName, errChan)

	// Subscribe to the topic
	sub, err := topicHandle.Subscribe()
	if err != nil {
		panic(err)
	}

	// TODO: Remove this, it's just for testing
	go streamConsoleTo(ctx, topicHandle)

	// Initialize the handler and start it
	handler := handler.NewHandler(sub, topicHandle)
	cli.Handler = handler
	go handler.Start(ctx, errChan)

	// Start the proxy server
	proxy := proxy.NewProxy(handler, "http://localhost:8545", "8080")
	cli.Proxy = proxy
	go proxy.Start()

	select {
	case err := <-errChan:
		log.Fatal("Error:", err)
	default:
		// do nothing
	}

}

func (cli *Client) startDiscovery(ctx context.Context, topicName string, errChan chan error) {
	//Start mDNS discovery
	notifee := &discoveryNotifee{h: cli.Host, ctx: ctx}
	mdns := mdns.NewMdnsService(cli.Host, "", notifee)
	if err := mdns.Start(); err != nil {
		errChan <- err
		return
	}

	// Start the DHT
	err := initDHT(ctx, cli)
	if err != nil {
		errChan <- err
		return
	}

	// Advertisement
	routingDiscovery := drouting.NewRoutingDiscovery(cli.DHT)
	dutil.Advertise(ctx, routingDiscovery, topicName)

	// Look for others who have announced and attempt to connect to them
	anyConnected := false
	for !anyConnected {
		fmt.Println("Searching for peers...")
		peerChan, err := routingDiscovery.FindPeers(ctx, topicName)
		if err != nil {
			errChan <- err
			return
		}
		for peer := range peerChan {
			if peer.ID == cli.Host.ID() {
				continue // No self connection
			}
			err := cli.Host.Connect(ctx, peer)
			if err != nil {
				fmt.Printf("Failed connecting to %s, error: %s\n", peer.ID, err)
			} else {
				fmt.Println("Connected to:", peer.ID)
				anyConnected = true
			}
		}
	}
	fmt.Println("Peer discovery complete")
	return
}

func (cli *Client) startPubsub(ctx context.Context, topicName string, errChan chan error) (topic *pubsub.Topic) {

	inspector := func(pid peer.ID, rpc *pubsub.RPC) error {

		if utils.IsOperator(pid) {
			return nil
		} else {
			return errors.New("not a operator")
		}

	}

	// Create a new PubSub service using the GossipSub router
	opts := []pubsub.Option{
		pubsub.WithMessageAuthor(cli.Host.ID()),
		pubsub.WithStrictSignatureVerification(true),
		pubsub.WithAppSpecificRpcInspector(inspector),
	}

	ps, err := pubsub.NewGossipSub(ctx, cli.Host, opts...)
	if err != nil {
		errChan <- err
		return
	}
	cli.PubSub = ps

	topicHandle, err := ps.Join(topicName)
	if err != nil {
		errChan <- err
		return
	}

	return topicHandle
}

func (m *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	if m.h.Network().Connectedness(pi.ID) != network.Connected {
		fmt.Printf("Found %s!\n", pi.ID.ShortString())
		m.h.Connect(m.ctx, pi)
	}
}

func initDHT(ctx context.Context, cli *Client) error {
	// Start a DHT, for use in peer discovery. We can't just make a new DHT
	// client because we want each peer to maintain its own local copy of the
	// DHT, so that the bootstrapping node of the DHT can go down without
	// inhibiting future peer discovery.
	kademliaDHT := cli.DHT

	if err := kademliaDHT.Bootstrap(ctx); err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	for _, peerAddr := range kaddht.DefaultBootstrapPeers {
		peerinfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := cli.Host.Connect(ctx, *peerinfo); err != nil {
				fmt.Println("Bootstrap warning:", err)
			}
		}()
	}
	wg.Wait()

	return nil
}

func streamConsoleTo(ctx context.Context, topic *pubsub.Topic) {
	reader := bufio.NewReader(os.Stdin)
	for {
		s, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		if err := topic.Publish(ctx, []byte(s)); err != nil {
			fmt.Println("### Publish error:", err)
		}
	}
}

func (cli *Client) AddOperator(p peer.ID) {
	cli.operators = append(cli.operators, p)
}

func (cli *Client) RemoveOperator(p peer.ID) {
	for i, pub := range cli.operators {
		if pub == p {
			cli.operators = append(cli.operators[:i], cli.operators[i+1:]...)
			break
		}
	}
}

func (cli *Client) GetOperators() []peer.ID {
	return cli.operators
}

func (cli *Client) IsPublisher(p peer.ID) bool {
	for _, pub := range cli.operators {
		if pub == p {
			return true
		}
	}
	return false
}

func (cli *Client) GetHandler() *handler.Handler {
	return cli.Handler
}
