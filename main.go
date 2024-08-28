package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/korayakpinar/network/src/client"
	"github.com/korayakpinar/network/src/ipfs"
	"github.com/korayakpinar/network/src/utils"
	"github.com/libp2p/go-libp2p"
	kaddht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/routing"
	"github.com/libp2p/go-libp2p/p2p/net/connmgr"
	"github.com/libp2p/go-libp2p/p2p/transport/tcp"
)

var (
	topicName      = flag.String("topic", "test", "Topic name")
	privKeyPath    = flag.String("privKey", "", "Private key path")
	apiPort        = flag.String("apiPort", "8081", "Port for the API server")
	proxyPort      = flag.String("proxyPort", "8082", "Port for the proxy server")
	rpcURL         = flag.String("rpcURL", "", "URL of the RPC server")
	contractAddr   = flag.String("contractAddr", "", "Address of the smart contract")
	committeeSize  = flag.Uint64("committeeSize", 32, "Size of the committee")
	ipfsGatewayURL = flag.String("ipfsGatewayURL", "", "URL of the IPFS gateway server")
)

func main() {
	flag.Parse()
	if *privKeyPath == "" || *rpcURL == "" || *contractAddr == "" || *ipfsGatewayURL == "" {
		fmt.Println("Please provide all the required arguments")
		return
	}

	privKey, err := readPrivateKeyFromFile(*privKeyPath)
	if err != nil {
		fmt.Printf("Failed to read private key: %v\n", err)
		return
	}

	ctx := context.Background()

	h, dht, err := setupLibp2p(ctx, privKey)
	if err != nil {
		fmt.Printf("Failed to setup libp2p: %v\n", err)
		return
	}

	fmt.Println("Host created, ID:", h.ID(), h.Addrs())
	ethAddr := utils.IdToEthAddress(h.ID())
	fmt.Println("Pub Addr:", ethAddr)

	// Initialize the IPFS service
	ipfsService := ipfs.NewIPFSService(*ipfsGatewayURL)

	// Initialize the client
	client := client.NewClient(h, dht, ipfsService, *apiPort, *proxyPort, *rpcURL, *contractAddr, privKey, *committeeSize)

	// Initialize the client
	if err := client.Initialize(ctx); err != nil {
		fmt.Printf("Failed to initialize client: %v\n", err)
		return
	}

	// Perform bootstrapping
	if err := client.Bootstrap(ctx); err != nil {
		fmt.Printf("Failed to bootstrap: %v\n", err)
		return
	}

	// Start the client
	if err := client.Start(ctx, *topicName); err != nil {
		fmt.Printf("Client error: %v\n", err)
		return
	}

	// Keep the main goroutine alive
	select {}
}

func readPrivateKeyFromFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("failed to open private key file: %w", err)
	}
	defer file.Close()

	privKeyBytes, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read private key file: %w", err)
	}

	return string(privKeyBytes), nil
}

func setupLibp2p(ctx context.Context, privKey string) (host.Host, *kaddht.IpfsDHT, error) {
	priv, err := utils.EthToLibp2pPrivKey(privKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to convert private key: %w", err)
	}

	connmgr, err := connmgr.NewConnManager(
		5,  // Lowwater
		10, // HighWater,
		connmgr.WithGracePeriod(time.Minute),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create connection manager: %w", err)
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
		libp2p.ConnectionManager(connmgr),
		libp2p.Routing(newDHT),
		libp2p.Transport(tcp.NewTCPTransport),
		libp2p.EnableNATService(),
		libp2p.EnableHolePunching(),
		libp2p.NATPortMap(),
		libp2p.DisableMetrics(),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create libp2p host: %w", err)
	}

	return h, dht, nil
}
