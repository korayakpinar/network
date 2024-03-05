package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/korayakpinar/p2pclient/src/client"
	"github.com/korayakpinar/p2pclient/src/utils"
	"github.com/libp2p/go-libp2p"
	kaddht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/routing"
	"github.com/libp2p/go-libp2p/p2p/net/connmgr"
)

var (
	cfgPath = flag.String("config", "./config.toml", "Path to the config file")
)

func main() {

	flag.Parse()
	ctx := context.Background()

	cfg, err := utils.LoadConfig(*cfgPath)
	if err != nil {
		panic(err)
	}

	priv, err := utils.EthToLibp2pPrivKey(cfg.PrivKey)
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
	client := client.NewClient(h, dht)

	fmt.Println("Host created, ID:", h.ID())
	ethAddr, err := utils.IdToEthAddress(h.ID())
	if err != nil {
		panic(err)

	}
	fmt.Println("Pub Addr:", ethAddr)

	go client.Start(ctx, cfg.TopicName)

	select {}
}
