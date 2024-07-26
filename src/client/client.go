package client

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
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
		log.Panicln("Couldn't dial to rpc, error: ", err)
		return
	}

	operatorsAddr := common.HexToAddress(cli.contractAddr)
	operatorsContract, err := contracts.NewOperators(operatorsAddr, ethClient)
	if err != nil {
		log.Panicln("Couldn't create operators contract, error: ", err)
		return
	}

	ecdsaPrivKey, err := crypto.HexToECDSA(cli.privKey)
	if err != nil {
		log.Panicln("Couldn't get ECDSA private key from hex, error: ", err)
		return
	}

	chainID, err := ethClient.ChainID(context.Background())
	if err != nil {
		log.Panicln("Couldn't get chain id from the RPC, error: ", err)
		return
	}

	auth, err := bind.NewKeyedTransactorWithChainID(ecdsaPrivKey, chainID)
	if err != nil {
		log.Panicln("Couldn't get the keyed transactor, error: ", err)
		return
	}

	registered, err := operatorsContract.IsRegistered(nil, auth.From)
	if err != nil {
		log.Panicln("Is registered call went wrong, error: ", err)
		return
	}

	if !registered {

		api := api.NewCrypto(cli.apiPort)

		tx, err := executeTransactionWithRetry(ethClient, auth, func(auth *bind.TransactOpts) (*types.Transaction, error) {
			return operatorsContract.RegisterOperator(auth)
		})

		if err != nil {
			log.Panicln("Registration transaction couldn't be successful, error: ", err)
			return
		}

		// Wait for the operator to be registered
		waitForTxConfirmation(ethClient, tx, 2)

		log.Println("Operator registered successfully")

		ourIndex, err := operatorsContract.GetOperatorIndex(nil, auth.From)
		if err != nil {
			log.Panicln("Get operator by index call went wrong, error: ", err)
			return
		}

		// Get and deploy BLS Public Key
		blsPubKey, err := api.GetPK(ourIndex.Uint64(), cli.committeeSize)
		fmt.Println("BLS Public Key: ", blsPubKey)
		if err != nil {
			log.Panicln("API couldn't send the BLS Public Key, error: ", err)
			return
		}

		tx, err = executeTransactionWithRetry(ethClient, auth, func(auth *bind.TransactOpts) (*types.Transaction, error) {
			return operatorsContract.SubmitBlsKey(auth, blsPubKey)
		})
		if err != nil {
			log.Panicln("Submit BLS Public Key transaction couldn't be successful, error: ", err)
			return
		}

		// Wait for the BLS Public Key to be submitted
		waitForTxConfirmation(ethClient, tx, 2)

	}
	log.Println("Operator registered successfully or already registered")

	log.Println("BLS Public Key submitted successfully or already submitted")

	var operatorCount *big.Int = big.NewInt(0)
	for operatorCount.Int64() < int64(cli.committeeSize) {
		operatorCount, err = operatorsContract.GetOperatorCount(nil)
		if err != nil {
			log.Panicln("Couldn't get the operator count from the RPC, error: ", err)
			return
		}
		log.Println("Waiting for all operators to be registered...")
		time.Sleep(3 * time.Second)
	}

	var ourIndex *big.Int
	var signers []handler.Signer
	for i := 0; i < int(operatorCount.Int64()); i++ {
		operator, err := operatorsContract.Operators(nil, big.NewInt(int64(i)))
		if err != nil {
			log.Panicln("Couldn't get the operators from the RPC, error: ", err)
			return
		}

		signerKey, err := operatorsContract.GetBLSPubKeyByIndex(nil, big.NewInt(int64(i)))
		if err != nil {
			log.Panicln("Couldn't get the BLS Pub Key by Index from the contract, error: ", err)
			return
		}

		if operator.Operator == auth.From {
			ourIndex = big.NewInt(int64(i))
		}
		newSigner := handler.NewSigner(operator.Operator, signerKey)
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

	/* inspector := func(pid peer.ID, rpc *pubsub.RPC) error {
		ethAddr := utils.IdToEthAddress(pid)

		signers := *cli.Handler.GetSigners()
		for _, signer := range signers {
			if signer.GetAddress() == ethAddr {
				return nil
			}
		}

		return errors.New("not a operator")
	} */

	// Create a new PubSub service using the GossipSub router
	opts := []pubsub.Option{
		pubsub.WithMessageAuthor(cli.Host.ID()),
		pubsub.WithStrictSignatureVerification(true),
		/* pubsub.WithAppSpecificRpcInspector(inspector), */
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

func executeTransactionWithRetry(
	client *ethclient.Client,
	auth *bind.TransactOpts,
	txFunc func(*bind.TransactOpts) (*types.Transaction, error),
) (*types.Transaction, error) {
	var tx *types.Transaction
	var err error

	for i := 0; i < 5; i++ { // 5 tries
		// Check the current nonce and set it in the auth
		nonce, err := client.PendingNonceAt(context.Background(), auth.From)
		if err != nil {
			return nil, fmt.Errorf("failed to get nonce: %v", err)
		}
		auth.Nonce = big.NewInt(int64(nonce))

		// Take the suggested gas price
		gasPrice, err := client.SuggestGasPrice(context.Background())
		if err != nil {
			return nil, fmt.Errorf("failed to get gas price: %v", err)
		}

		// Increase gas price by 10%
		gasPrice = new(big.Int).Mul(gasPrice, big.NewInt(110))
		gasPrice = new(big.Int).Div(gasPrice, big.NewInt(100))
		auth.GasPrice = gasPrice

		// Call the transaction function
		tx, err = txFunc(auth)
		if err == nil {
			return tx, nil
		}

		if !strings.Contains(err.Error(), "replacement transaction underpriced") &&
			!strings.Contains(err.Error(), "nonce too low") {
			return nil, err
		}

		log.Printf("Transaction failed, retrying with updated nonce and higher gas price... (Attempt %d)\n", i+1)
		time.Sleep(time.Second * 2) // Kısa bir bekleme süresi
	}

	return nil, fmt.Errorf("failed to execute transaction after multiple attempts: %v", err)
}
