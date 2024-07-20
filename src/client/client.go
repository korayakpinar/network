package client

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"math/big"
	"os"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/korayakpinar/network/src/contracts"
	api "github.com/korayakpinar/network/src/crypto"
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
	Host    host.Host
	PubSub  *pubsub.PubSub
	DHT     *kaddht.IpfsDHT
	Handler *handler.Handler
	Proxy   *proxy.Proxy

	signers       *[]handler.Signer
	proxyPort     string
	rpcUrl        string
	contractAddr  string
	privKey       string
	apiPort       string
	committeeSize uint64
}

type discoveryNotifee struct {
	h   host.Host
	ctx context.Context
}

func NewClient(h host.Host, dht *kaddht.IpfsDHT, apiPort, proxyPort, rpcUrl, contractAddr, privKey string, committeSize uint64) *Client {
	signerArr := make([]handler.Signer, 0)
	return &Client{h, nil, dht, nil, nil, &signerArr, proxyPort, rpcUrl, contractAddr, privKey, apiPort, committeSize}
}

func (cli *Client) Start(ctx context.Context, topicName string) {

	// Initialize ethClient and register the operator
	ethClient, err := ethclient.Dial(cli.rpcUrl)
	if err != nil {
		log.Fatal(err)
		return
	}

	operatorsAddr := common.HexToAddress(cli.contractAddr)
	operatorsContract, err := contracts.NewOperators(operatorsAddr, ethClient)
	if err != nil {
		log.Fatal(err)
		return
	}

	ecdsaPrivKey, err := crypto.HexToECDSA(cli.privKey)
	if err != nil {
		log.Fatal(err)
		return
	}

	chainID, err := ethClient.ChainID(context.Background())
	if err != nil {
		log.Fatal(err)
		return
	}

	auth, err := bind.NewKeyedTransactorWithChainID(ecdsaPrivKey, chainID) // 1 yerine doğru chain ID'yi kullanın
	if err != nil {
		return
	}

	tx, err := operatorsContract.RegisterOperator(auth)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Printf("Operator registered with tx: %s\n", tx.Hash().Hex())

	// Wait for the operator to be registered
	waitForTxConfirmation(ethClient, tx, 2)

	ourIndex, err := operatorsContract.GetOperatorIndex(nil, auth.From)
	if err != nil {
		log.Fatal(err)
		return
	}

	api := api.NewCrypto(cli.apiPort)

	// Get and deploy BLS Public Key
	blsPubKey, err := api.GetPK(ourIndex.Uint64(), cli.committeeSize)
	if err != nil {
		log.Fatal(err)
		return
	}

	gasPrice, err := ethClient.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
		return
	}

	nonce, err := ethClient.PendingNonceAt(context.Background(), auth.From)
	if err != nil {
		log.Fatal(err)
		return
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = 3000000
	auth.GasPrice = gasPrice

	keystoreABI, err := contracts.BLSKeystoreMetaData.GetAbi()
	if err != nil {
		log.Fatal(err)
		return
	}
	keystoreAddr, tx, _, err := bind.DeployContract(auth, *keystoreABI, []byte(contracts.BLSKeystoreMetaData.Bin), ethClient, blsPubKey)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Printf("Keystore deployed with tx: %s\n", tx.Hash().Hex())

	// Wait for the keystore to be deployed
	log.Printf("Waiting for keystore to be deployed...\n")
	waitForTxConfirmation(ethClient, tx, 2)

	tx, err = operatorsContract.SubmitBLSPubKey(auth, keystoreAddr)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Printf("BLS Public Key submitted with tx: %s\n", tx.Hash().Hex())
	log.Printf("Waiting for BLS Public Key to be submitted...\n")
	waitForTxConfirmation(ethClient, tx, 2)

	var operatorCount *big.Int = big.NewInt(0)
	for operatorCount.Int64() < int64(cli.committeeSize) {
		operatorCount, err = operatorsContract.GetOperatorCount(nil)
		if err != nil {
			log.Fatal(err)
			return
		}
		log.Println("Waiting for all operators to be registered...")
		time.Sleep(3 * time.Second)
	}

	var signers []handler.Signer
	for i := 0; i < int(operatorCount.Int64()); i++ {
		operator, err := operatorsContract.Operators(nil, big.NewInt(int64(i)))
		if err != nil {
			log.Fatal(err)
			return
		}
		signerPeerID, err := utils.IdFromPubKey(operator.Operator.Hex())
		if err != nil {
			log.Fatal(err)
			return
		}

		signerKeystore, err := contracts.NewBLSKeystore(keystoreAddr, ethClient)
		if err != nil {
			log.Fatal(err)
			return
		}

		signerKey, err := signerKeystore.GetPubKey(nil)
		if err != nil {
			log.Fatal(err)
			return
		}

		newSigner := handler.NewSigner(*signerPeerID, operator.Operator, signerKey)
		signers = append(signers, newSigner)
	}

	fmt.Println("All operators registered successfully, ", len(signers))
	fmt.Println("Starting the client...")
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
	handler := handler.NewHandler(sub, topicHandle, &signers, cli.privKey, cli.rpcUrl, cli.apiPort, cli.committeeSize, ourIndex.Uint64(), cli.committeeSize/2)
	cli.Handler = handler
	go handler.Start(ctx, errChan)

	// Start the proxy server
	proxy := proxy.NewProxy(handler, cli.rpcUrl, cli.proxyPort)
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

		if cli.Handler.IsSigner(pid) {
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

func (cli *Client) GetHandler() *handler.Handler {
	return cli.Handler
}

func waitForTxConfirmation(client *ethclient.Client, tx *types.Transaction, blockConfirmations uint64) error {
	ctx := context.Background()
	receipt, err := waitForTxReceipt(client, ctx, tx.Hash(), 2*time.Minute)
	if err != nil {
		return fmt.Errorf("error waiting for transaction receipt: %v", err)
	}

	if receipt.Status == 0 {
		return fmt.Errorf("transaction failed")
	}

	fmt.Printf("Transaction included in block %d\n", receipt.BlockNumber.Uint64())

	currentBlock, err := client.BlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("error getting current block number: %v", err)
	}

	confirmations := currentBlock - receipt.BlockNumber.Uint64()
	for confirmations < blockConfirmations {
		time.Sleep(15 * time.Second)
		currentBlock, err = client.BlockNumber(ctx)
		if err != nil {
			return fmt.Errorf("error getting current block number: %v", err)
		}
		confirmations = currentBlock - receipt.BlockNumber.Uint64()
		fmt.Printf("Current confirmations: %d\n", confirmations)
	}

	fmt.Printf("Transaction confirmed with %d block confirmations\n", confirmations)
	return nil
}

func waitForTxReceipt(client *ethclient.Client, ctx context.Context, txHash common.Hash, timeout time.Duration) (*types.Receipt, error) {
	timeoutCh := time.After(timeout)
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		receipt, err := client.TransactionReceipt(ctx, txHash)
		if err == nil {
			return receipt, nil
		}

		select {
		case <-timeoutCh:
			return nil, fmt.Errorf("timeout waiting for transaction receipt")
		case <-ticker.C:
			continue
		}
	}
}
