package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/korayakpinar/network/src/client"
	"github.com/korayakpinar/network/src/utils"
	"github.com/libp2p/go-libp2p"
	kaddht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/routing"
	"github.com/libp2p/go-libp2p/p2p/net/connmgr"
)

var (
	topicName     = flag.String("topic", "test", "Topic name")
	privKey       = flag.String("privKey", "", "Private key in hex format")
	proxyPort     = flag.String("proxyPort", "8082", "Port for the proxy server")
	rpcURL        = flag.String("rpcURL", "", "URL of the RPC server")
	contractAddr  = flag.String("contractAddr", "", "Address of the smart contract")
	committeeSize = flag.Uint64("committeeSize", 32, "Size of the committee")
)

func main() {
	flag.Parse()
	ctx := context.Background()

	priv, err := utils.EthToLibp2pPrivKey(*privKey)
	if err != nil {
		panic(err)
	}

	connmgr, err := connmgr.NewConnManager(
		0,   // Lowwater
		250, // HighWater,
		connmgr.WithGracePeriod(time.Minute),
	)
	if err != nil {
		panic(err)
	}

	var dht *kaddht.IpfsDHT
	newDHT := func(h host.Host) (routing.PeerRouting, error) {
		var err error
		dht, err = kaddht.New(ctx, h)
		return dht, err
	}

	h, err := libp2p.New(
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"),
		libp2p.Identity(priv),
		libp2p.EnableHolePunching(),
		libp2p.NATPortMap(),
		libp2p.ConnectionManager(connmgr),
		libp2p.EnableNATService(),
		libp2p.Routing(newDHT),
	)
	if err != nil {
		panic(err)
	}

	// Initialize the client
	client := client.NewClient(h, dht, *proxyPort, *rpcURL, *contractAddr, *privKey, *committeeSize)

	fmt.Println("Host created, ID:", h.ID())
	ethAddr, err := utils.IdToEthAddress(h.ID())
	if err != nil {
		panic(err)

	}
	fmt.Println("Pub Addr:", ethAddr)

	go client.Start(ctx, *topicName)

	select {}
}
