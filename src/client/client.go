package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/korayakpinar/network/src/contracts"
	"github.com/korayakpinar/network/src/handler"
	"github.com/korayakpinar/network/src/types"

	"github.com/korayakpinar/network/src/ipfs"
	"github.com/korayakpinar/network/src/proxy"
	kaddht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	libp2pCrypto "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
)

type Client struct {
	Host              host.Host
	PubSub            *pubsub.PubSub
	DHT               *kaddht.IpfsDHT
	Handler           *handler.Handler
	Proxy             *proxy.Proxy
	IPFSService       *ipfs.IPFSService
	signers           map[uint64]*types.Signer
	proxyPort         string
	rpcUrl            string
	contractAddr      string
	privKey           *libp2pCrypto.PrivKey
	keystore          *types.Keystore
	apiPort           string
	committeeSize     uint64
	networkSize       uint64
	networkIndex      uint64
	ethClient         *ethclient.Client
	operatorsContract *contracts.Operators
	adminKey          string
}

const minDelay = time.Second * 5
const maxDelay = time.Second * 10

func NewClient(h host.Host, dht *kaddht.IpfsDHT, ipfsService *ipfs.IPFSService, apiPort, proxyPort, rpcUrl, contractAddr, adminKey string, privKey *libp2pCrypto.PrivKey, keystore *types.Keystore, committeSize, networkSize, networkIndex uint64) *Client {
	signers := make(map[uint64]*types.Signer, 0)
	return &Client{
		Host:          h,
		DHT:           dht,
		IPFSService:   ipfsService,
		signers:       signers,
		proxyPort:     proxyPort,
		rpcUrl:        rpcUrl,
		contractAddr:  contractAddr,
		privKey:       privKey,
		keystore:      keystore,
		apiPort:       apiPort,
		committeeSize: committeSize,
		networkSize:   networkSize,
		networkIndex:  networkIndex,
		adminKey:      adminKey,
	}
}

func (cli *Client) Initialize(ctx context.Context) error {
	var err error

	err = cli.executeRPCCallWithRetry(ctx, func() error {
		cli.ethClient, err = ethclient.Dial(cli.rpcUrl)
		return err
	})
	if err != nil {
		return fmt.Errorf("couldn't dial to rpc: %w", err)
	}

	operatorsAddr := common.HexToAddress(cli.contractAddr)
	err = cli.executeRPCCallWithRetry(ctx, func() error {
		cli.operatorsContract, err = contracts.NewOperators(operatorsAddr, cli.ethClient)
		return err
	})
	if err != nil {
		return fmt.Errorf("couldn't create operators contract: %w", err)
	}

	return nil
}

func (cli *Client) Bootstrap(ctx context.Context) error {
	attempt := 0

	for {
		attempt++
		err := cli.bootstrapAttempt(ctx)
		if err == nil {
			return nil
		}
		log.Printf("Bootstrap attempt %d failed: %v", attempt, err)

		delay := time.Duration(rand.Int63n(int64(maxDelay-minDelay))) + minDelay
		log.Printf("Retrying in %v...", delay)
		select {
		case <-time.After(delay):
			// Continue with the next iteration
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (cli *Client) bootstrapAttempt(ctx context.Context) error {

	contractSigners, err := cli.operatorsContract.GetAllOperators(&bind.CallOpts{})
	if err != nil {
		signerCount, err := cli.operatorsContract.GetOperatorCount(&bind.CallOpts{})
		if err != nil {
			return fmt.Errorf("couldn't get the operator count: %w", err)
		}
		if signerCount != big.NewInt(int64(cli.committeeSize-1)) {
			log.Panicf("Signer count mismatch: expected %d, got %d", cli.committeeSize-1, signerCount)
			return fmt.Errorf("failed to get all operators: %w", err)
		}

		// Take every operator that has been registered so far one by one
		for i := 0; i < int(signerCount.Int64()); i++ {
			operator, err := cli.operatorsContract.Operators(nil, big.NewInt(int64(i)))
			if err != nil {
				return fmt.Errorf("couldn't get operator: %w", err)
			}
			contractSigners = append(contractSigners, operator)
		}
	}

	if err := cli.InitializeSigners(contractSigners); err != nil {
		return fmt.Errorf("failed to initialize signers: %w", err)
	}

	return nil
}

/*
func (cli *Client) registerAllOperators(ctx context.Context) error {
	api := api.NewCrypto(cli.apiPort)

	ecdsaPrivKey, err := crypto.HexToECDSA(cli.adminKey)
	if err != nil {
		return fmt.Errorf("couldn't get ECDSA private key from hex: %w, key: %v", err, cli.adminKey)
	}

	chainID, err := cli.ethClient.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("couldn't get chain ID: %w", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(ecdsaPrivKey, chainID)
	if err != nil {
		return fmt.Errorf("couldn't create auth: %w", err)
	}

	operators := make([]contracts.OperatorsContractNode, 0, len(cli.keystore.Keys))
	for i, key := range cli.keystore.Keys {
		operatorPrivKey, err := crypto.HexToECDSA(key.ECDSAPrivateKey)
		if err != nil {
			return fmt.Errorf("couldn't get ECDSA private key from hex: %w", err)
		}

		index := (cli.committeeSize/cli.networkSize)*cli.networkIndex + uint64(i)
		blsPubKey, err := api.GetPK([]byte(key.BLSPrivateKey), index, cli.committeeSize)
		if err != nil {
			return fmt.Errorf("couldn't get BLS public key: %w", err)
		}

		if string(blsPubKey) != key.BLSPublicKey {
			return fmt.Errorf("BLS public key mismatch: expected %s, got %s", key.BLSPublicKey, blsPubKey)
		}

		operatorAddress := crypto.PubkeyToAddress(operatorPrivKey.PublicKey)
		cid, err := cli.uploadAndVerifyBLSKey(ctx, []byte(key.BLSPublicKey))
		if err != nil {
			return fmt.Errorf("couldn't upload and verify BLS key: %w", err)
		}

		operators = append(operators, contracts.OperatorsContractNode{
			Operator:     operatorAddress,
			BlsPubKeyCID: cid,
		})
	}

	tx, err := cli.operatorsContract.AppendMultipleOperators(auth, operators)
	if err != nil {
		return fmt.Errorf("failed to append multiple operators: %w", err)
	}

	if err := utils.WaitForTxConfirmationWithRetry(ctx, cli.ethClient, tx, 2); err != nil {
		return fmt.Errorf("failed to wait for AppendMultipleOperators confirmation: %w", err)
	}

	log.Println("All operators registered successfully")
	return nil
}

func (cli *Client) waitForAllOperators(ctx context.Context) error {
	for {
		operatorCount, err := cli.operatorsContract.GetOperatorCount(&bind.CallOpts{})
		if err != nil {
			return fmt.Errorf("couldn't get the operator count: %w", err)
		}
		if operatorCount.Cmp(big.NewInt(int64(cli.committeeSize-1))) >= 0 {
			break
		}
		log.Printf("Waiting for all operators to be registered (%d/%d)...\n", operatorCount, cli.committeeSize-1)
		time.Sleep(3 * time.Second)
	}

	log.Printf("All operators registered successfully: %d\n", cli.committeeSize-1)
	return nil
}

func (cli *Client) uploadAndVerifyBLSKey(ctx context.Context, blsPubKey []byte) (string, error) {
	cid, err := cli.IPFSService.UploadKey(blsPubKey)
	if err != nil {
		return "", fmt.Errorf("couldn't upload the BLS Public Key to IPFS: %w", err)
	}

	for i := 0; i < 5; i++ { // Try 5 times
		// Random delay between 1 and 5 seconds
		delay := time.Duration(rand.Intn(4)+1) * time.Second
		time.Sleep(delay)

		retrievedKey, err := cli.IPFSService.GetKeyByCID(cid)
		if err == nil && string(retrievedKey) == string(blsPubKey) {
			return cid, nil
		}
		log.Printf("Failed to verify uploaded key, retrying... (Attempt %d)\n", i+1)
	}

	return "", fmt.Errorf("failed to verify uploaded BLS key after multiple attempts")
} */

func (c *Client) InitializeSigners(contractSigners []contracts.OperatorsContractNode) error {
	c.signers = make(map[uint64]*types.Signer)

	for index, contractSigner := range contractSigners {
		key, err := c.IPFSService.GetKeyByCID(contractSigner.BlsPubKeyCID)
		if err != nil {
			return fmt.Errorf("couldn't get key by CID: %w", err)
		}

		signer := &types.Signer{
			Index:        uint64(index),
			Address:      contractSigner.Operator,
			BLSPublicKey: key,
			IsLocal:      false,
		}

		if keyPair, exists := c.keystore.Keys[uint64(index)]; exists {
			signer.IsLocal = true
			signer.KeyPair = &keyPair
		}

		c.signers[uint64(index)] = signer
	}

	return nil
}

func (c *Client) GetSignerInfo(index uint64) (*types.Signer, bool) {
	signer, exists := c.signers[index]
	return signer, exists
}

func (c *Client) GetLocalSigners() []*types.Signer {
	var localSigners []*types.Signer
	for _, signer := range c.signers {
		if signer.IsLocal {
			localSigners = append(localSigners, signer)
		}
	}
	return localSigners
}

func (c *Client) GetAllSigners() []*types.Signer {
	var allSigners []*types.Signer
	for _, signer := range c.signers {
		allSigners = append(allSigners, signer)
	}
	return allSigners
}

func (c *Client) UpdateSignerIndex(oldIndex, newIndex uint64) error {
	if signer, exists := c.signers[oldIndex]; exists {
		signer.Index = newIndex
		c.signers[newIndex] = signer
		delete(c.signers, oldIndex)
		return nil
	}
	return fmt.Errorf("signer with index %d not found", oldIndex)
}

func (cli *Client) Start(ctx context.Context, topicName string) error {
	log.Println("Starting the client...")

	errChan := make(chan error, 4)

	go cli.startDiscovery(ctx, topicName, errChan)

	topicHandle, err := cli.startPubsub(ctx, topicName)
	if err != nil {
		return fmt.Errorf("failed to start pubsub: %w", err)
	}

	sub, err := topicHandle.Subscribe()
	if err != nil {
		return fmt.Errorf("failed to subscribe to topic: %w", err)
	}

	cli.Handler = handler.NewHandler(sub, topicHandle, cli.signers, *cli.privKey, cli.rpcUrl, cli.apiPort, cli.committeeSize-1, cli.committeeSize/2, cli.networkIndex, cli.networkSize)
	log.Println("Index:", cli.networkIndex)
	go cli.Handler.Start(ctx, errChan)

	cli.Proxy = proxy.NewProxy(cli.Handler, cli.rpcUrl, cli.proxyPort)
	go cli.Proxy.Start()

	select {
	case err := <-errChan:
		return fmt.Errorf("error in client components: %w", err)
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (cli *Client) startDiscovery(ctx context.Context, topicName string, errChan chan<- error) {
	// Start mDNS discovery
	notifee := &discoveryNotifee{h: cli.Host, ctx: ctx}
	mdns := mdns.NewMdnsService(cli.Host, "", notifee)
	if err := mdns.Start(); err != nil {
		errChan <- fmt.Errorf("failed to start mDNS: %w", err)
		return
	}

	// Start the DHT
	if err := cli.initDHT(ctx); err != nil {
		errChan <- fmt.Errorf("failed to initialize DHT: %w", err)
		return
	}

	// Use a routing discovery to find peers
	routingDiscovery := drouting.NewRoutingDiscovery(cli.DHT)
	dutil.Advertise(ctx, routingDiscovery, topicName)

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	anyConnected := false
	for !anyConnected {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			peers, err := routingDiscovery.FindPeers(ctx, topicName)
			if err != nil {
				log.Printf("Error finding peers: %v", err)
				continue
			}
			for peer := range peers {
				if peer.ID == cli.Host.ID() {
					continue // Skip self
				}
				if err := cli.Host.Connect(ctx, peer); err != nil {
					log.Printf("Failed connecting to %s: %v", peer.ID, err)
				} else {
					anyConnected = true
				}
			}
		}
	}
}

func (cli *Client) startPubsub(ctx context.Context, topicName string) (*pubsub.Topic, error) {
	// Create a new PubSub service using the GossipSub router
	ps, err := pubsub.NewGossipSub(ctx, cli.Host)
	if err != nil {
		return nil, fmt.Errorf("failed to create pubsub: %w", err)
	}
	cli.PubSub = ps

	// Join the topic
	topic, err := ps.Join(topicName)
	if err != nil {
		return nil, fmt.Errorf("failed to join topic: %w", err)
	}

	return topic, nil
}

func (cli *Client) initDHT(ctx context.Context) error {
	if err := cli.DHT.Bootstrap(ctx); err != nil {
		return fmt.Errorf("failed to bootstrap DHT: %w", err)
	}

	var wg sync.WaitGroup
	for _, peerAddr := range kaddht.DefaultBootstrapPeers {
		peerinfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := cli.Host.Connect(ctx, *peerinfo); err != nil {
				log.Printf("Error connecting to bootstrap peer %s: %v", peerinfo.ID, err)
			}
		}()
	}
	wg.Wait()

	return nil
}

type discoveryNotifee struct {
	h   host.Host
	ctx context.Context
}

func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	if n.h.Network().Connectedness(pi.ID) != network.Connected {
		log.Printf("Discovered new peer %s\n", pi.ID.ShortString())
		if err := n.h.Connect(n.ctx, pi); err != nil {
			log.Printf("Failed to connect to peer %s: %v\n", pi.ID.ShortString(), err)
		}
	}
}

func (cli *Client) GetHandler() *handler.Handler {
	return cli.Handler
}

func (c *Client) executeRPCCallWithRetry(ctx context.Context, rpcCall func() error) error {
	for {
		err := rpcCall()
		if err == nil {
			return nil
		}

		if strings.Contains(err.Error(), "429 Too Many Requests") {
			var errorResponse struct {
				Error struct {
					Data struct {
						TryAgainIn string `json:"try_again_in"`
					} `json:"data"`
				} `json:"error"`
			}

			// Extract the JSON part from the error message
			parts := strings.SplitN(err.Error(), "{", 2)
			if len(parts) < 2 {
				return fmt.Errorf("unexpected error format: %v", err)
			}
			jsonPart := "{" + parts[1]

			if jsonErr := json.Unmarshal([]byte(jsonPart), &errorResponse); jsonErr == nil {
				waitDuration, parseErr := time.ParseDuration(errorResponse.Error.Data.TryAgainIn)
				if parseErr == nil {
					fmt.Printf("Rate limit exceeded. Will sleep for %v before retrying.\n", waitDuration)
					select {
					case <-time.After(waitDuration):
						continue
					case <-ctx.Done():
						return ctx.Err()
					}
				} else {
					return fmt.Errorf("failed to parse duration: %v", parseErr)
				}
			} else {
				return fmt.Errorf("failed to unmarshal JSON: %v", jsonErr)
			}
		}

		return err
	}
}
