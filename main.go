package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/korayakpinar/network/src/client"
	"github.com/korayakpinar/network/src/ipfs"
	"github.com/korayakpinar/network/src/types"
	"github.com/libp2p/go-libp2p"
	kaddht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/routing"
	"github.com/libp2p/go-libp2p/p2p/net/connmgr"
	"github.com/libp2p/go-libp2p/p2p/transport/tcp"
)

var (
	topicName      = flag.String("topic", "test", "Topic name")
	keysPath       = flag.String("keys", "", "Path to directory containing keys")
	apiPort        = flag.String("apiPort", "8081", "Port for the API server")
	proxyPort      = flag.String("proxyPort", "8082", "Port for the proxy server")
	rpcURL         = flag.String("rpcURL", "", "URL of the RPC server")
	contractAddr   = flag.String("contractAddr", "", "Address of the smart contract")
	committeeSize  = flag.Uint64("committeeSize", 512, "Size of the committee")
	networkSize    = flag.Uint64("networkSize", 16, "Size of the network")
	networkIndex   = flag.Uint64("networkIndex", 0, "Index of the node in the network")
	ipfsGatewayURL = flag.String("ipfsGatewayURL", "", "URL of the IPFS gateway server")
	adminKey       = flag.String("adminKey", "", "Admin key for bootstrapping")
)

func main() {
	flag.Parse()
	if *keysPath == "" || *rpcURL == "" || *contractAddr == "" || *ipfsGatewayURL == "" || *adminKey == "" {
		fmt.Println("Please provide all the required arguments")
		return
	}

	keystore, err := readKeyPairsFromPath(*keysPath, *networkSize, *networkIndex, *committeeSize)
	if err != nil {
		fmt.Printf("Failed to read private key: %v\n", err)
		return
	}

	// It will use secure random number generator to generate a key pair by default.
	priv, _, err := crypto.GenerateEd25519Key(nil)
	if err != nil {
		fmt.Printf("failed to convert private key: %v", err)
		return
	}

	ctx := context.Background()
	h, dht, err := setupLibp2p(ctx, priv)
	if err != nil {
		fmt.Printf("Failed to setup libp2p: %v\n", err)
		return
	}

	fmt.Println("Host created, ID:", h.ID(), h.Addrs())

	// Initialize the IPFS service
	ipfsService := ipfs.NewIPFSService(*ipfsGatewayURL)

	// Initialize the client
	client := client.NewClient(h, dht, ipfsService, *apiPort, *proxyPort, *rpcURL, *contractAddr, *adminKey, &priv, keystore, *committeeSize, *networkSize, *networkIndex)

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

func readKeyPairsFromPath(path string, networkSize, networkIndex, committeeSize uint64) (*types.Keystore, error) {
	keystore := &types.Keystore{
		Keys: make(map[uint64]types.KeyPair),
	}

	keysPerNode := committeeSize / networkSize
	startIndex := networkIndex * keysPerNode
	var endIndex uint64

	// Adjust for the last node
	if networkIndex == networkSize-1 {
		endIndex = committeeSize - 1 // Because key indices start from 0
	} else {
		endIndex = (networkIndex + 1) * keysPerNode
	}

	for i := startIndex; i < endIndex; i++ {
		fileIndex := i + 1 // Because files are named starting from 1
		ecdsaPath := filepath.Join(path, fmt.Sprintf("%d-ecdsa", fileIndex))
		blsPrivPath := filepath.Join(path, fmt.Sprintf("%d-bls", fileIndex))
		blsPubPath := filepath.Join(path, fmt.Sprintf("%d-pk", fileIndex))

		ecdsaKey, err := os.ReadFile(ecdsaPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read ECDSA key for index %d: %w", fileIndex, err)
		}

		blsPrivKey, err := os.ReadFile(blsPrivPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read BLS private key for index %d: %w", fileIndex, err)
		}

		blsPubKey, err := os.ReadFile(blsPubPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read BLS public key for index %d: %w", fileIndex, err)
		}

		keystore.Keys[i] = types.KeyPair{
			ECDSAPrivateKey: string(ecdsaKey),
			BLSPrivateKey:   string(blsPrivKey),
			BLSPublicKey:    string(blsPubKey),
		}
	}

	return keystore, nil
}

func setupLibp2p(ctx context.Context, privKey crypto.PrivKey) (host.Host, *kaddht.IpfsDHT, error) {
	connmgr, err := connmgr.NewConnManager(
		20, // Lowwater
		40, // HighWater,
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
		libp2p.Identity(privKey),
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
