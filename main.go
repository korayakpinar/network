package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/korayakpinar/network/src/client"
	"github.com/korayakpinar/network/src/pinata"
	"github.com/korayakpinar/network/src/utils"
	"github.com/libp2p/go-libp2p"
	kaddht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/routing"
	"github.com/libp2p/go-libp2p/p2p/net/connmgr"
)

var (
	topicName      = flag.String("topic", "test", "Topic name")
	privKey        = flag.String("privKey", "", "Private key in hex format")
	apiPort        = flag.String("apiPort", "8081", "Port for the API server")
	proxyPort      = flag.String("proxyPort", "8082", "Port for the proxy server")
	rpcURL         = flag.String("rpcURL", "", "URL of the RPC server")
	contractAddr   = flag.String("contractAddr", "", "Address of the smart contract")
	committeeSize  = flag.Uint64("committeeSize", 32, "Size of the committee")
	ipfsGatewayURL = flag.String("IpfsGatewayURL", "", "URL of the IPFS gateway server")
	bearerToken    = flag.String("bearerToken", "", "Bearer token for the IPFS gateway server")
)

func main() {
	flag.Parse()
	if *privKey == "" || *rpcURL == "" || *contractAddr == "" || *ipfsGatewayURL == "" || *bearerToken == "" {
		fmt.Println("Please provide all the required arguments")
		return
	}

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

	// Initialize the IPFS service
	ipfsService := pinata.NewIPFSService(*bearerToken, *ipfsGatewayURL)

	// Initialize the client
	client := client.NewClient(h, dht, ipfsService, *apiPort, *proxyPort, *rpcURL, *contractAddr, *privKey, *committeeSize)

	fmt.Println("Host created, ID:", h.ID())
	ethAddr := utils.IdToEthAddress(h.ID())

	fmt.Println("Pub Addr:", ethAddr)

	go client.Start(ctx, *topicName)

	select {}
}
