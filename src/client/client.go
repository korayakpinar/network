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
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/korayakpinar/network/src/contracts"
	"github.com/korayakpinar/network/src/utils"

	api "github.com/korayakpinar/network/src/crypto"
	"github.com/korayakpinar/network/src/handler"

	"github.com/korayakpinar/network/src/ipfs"
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
	Host              host.Host
	PubSub            *pubsub.PubSub
	DHT               *kaddht.IpfsDHT
	Handler           *handler.Handler
	Proxy             *proxy.Proxy
	IPFSService       *ipfs.IPFSService
	signers           *[]handler.Signer
	proxyPort         string
	rpcUrl            string
	contractAddr      string
	privKey           string
	apiPort           string
	committeeSize     uint64
	ethClient         *ethclient.Client
	operatorsContract *contracts.Operators
	auth              *bind.TransactOpts
}

func NewClient(h host.Host, dht *kaddht.IpfsDHT, ipfsService *ipfs.IPFSService, apiPort, proxyPort, rpcUrl, contractAddr, privKey string, committeSize uint64) *Client {
	signerArr := make([]handler.Signer, 0)
	return &Client{
		Host:          h,
		DHT:           dht,
		IPFSService:   ipfsService,
		signers:       &signerArr,
		proxyPort:     proxyPort,
		rpcUrl:        rpcUrl,
		contractAddr:  contractAddr,
		privKey:       privKey,
		apiPort:       apiPort,
		committeeSize: committeSize,
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

	ecdsaPrivKey, err := crypto.HexToECDSA(cli.privKey)
	if err != nil {
		return fmt.Errorf("couldn't get ECDSA private key from hex: %w", err)
	}

	var chainID *big.Int
	err = cli.executeRPCCallWithRetry(ctx, func() error {
		chainID, err = cli.ethClient.ChainID(ctx)
		return err
	})
	if err != nil {
		return fmt.Errorf("couldn't get chain id from the RPC: %w", err)
	}

	cli.auth, err = bind.NewKeyedTransactorWithChainID(ecdsaPrivKey, chainID)
	if err != nil {
		return fmt.Errorf("couldn't get the keyed transactor: %w", err)
	}

	return nil
}

func (cli *Client) Bootstrap(ctx context.Context) error {
	attempt := 0
	minDelay := time.Second * 5
	maxDelay := time.Second * 10

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
	if err := cli.registerOperator(ctx); err != nil {
		return fmt.Errorf("failed to register operator: %w", err)
	}

	if err := cli.submitBLSKey(ctx); err != nil {
		return fmt.Errorf("failed to submit BLS key: %w", err)
	}

	if err := cli.waitForAllOperators(ctx); err != nil {
		return fmt.Errorf("failed to wait for all operators: %w", err)
	}

	return nil
}

func (cli *Client) registerOperator(ctx context.Context) error {
	var registered bool
	var err error

	err = cli.executeRPCCallWithRetry(ctx, func() error {
		registered, err = cli.operatorsContract.IsRegistered(nil, cli.auth.From)
		return err
	})
	if err != nil {
		return fmt.Errorf("is registered call went wrong: %w", err)
	}

	if !registered {
		tx, err := utils.ExecuteTransactionWithRetry(cli.ethClient, cli.auth, func(auth *bind.TransactOpts) (*types.Transaction, error) {
			return cli.operatorsContract.RegisterOperator(auth)
		})
		if err != nil {
			return fmt.Errorf("registration transaction couldn't be successful: %w", err)
		}

		if err := utils.WaitForTxConfirmationWithRetry(ctx, cli.ethClient, tx, 2); err != nil {
			return fmt.Errorf("failed to wait for registration confirmation: %w", err)
		}

		log.Println("Operator registered successfully")
	} else {
		log.Println("Operator already registered")
	}

	return cli.waitForOperatorIndex(ctx)
}

func (cli *Client) waitForOperatorIndex(ctx context.Context) error {
	minDelay := time.Second
	maxDelay := time.Second * 5
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			var ourIndex *big.Int
			err := cli.executeRPCCallWithRetry(ctx, func() error {
				var err error
				ourIndex, err = cli.operatorsContract.GetOperatorIndex(nil, cli.auth.From)
				return err
			})
			if err == nil {
				log.Printf("Our operator index: %d\n", ourIndex)
				return nil
			}
			log.Println("Waiting for operator index to be available...")
			delay := time.Duration(rand.Int63n(int64(maxDelay-minDelay))) + minDelay
			time.Sleep(delay)

		}
	}
}

func (cli *Client) submitBLSKey(ctx context.Context) error {
	var submissions bool
	err := cli.executeRPCCallWithRetry(ctx, func() error {
		var err error
		submissions, err = cli.operatorsContract.HasSubmittedBLSKey(nil, cli.auth.From)
		return err
	})
	if err != nil {
		return fmt.Errorf("couldn't get the BLS Key submission status: %w", err)
	}

	if !submissions {
		var ourIndex *big.Int
		err := cli.executeRPCCallWithRetry(ctx, func() error {
			var err error
			ourIndex, err = cli.operatorsContract.GetOperatorIndex(nil, cli.auth.From)
			return err
		})
		if err != nil {
			return fmt.Errorf("get operator by index call went wrong: %w", err)
		}

		api := api.NewCrypto(cli.apiPort)
		blsPubKey, err := api.GetPK(ourIndex.Uint64(), cli.committeeSize)
		if err != nil {
			return fmt.Errorf("API couldn't send the BLS Public Key: %w", err)
		}

		cid, err := cli.uploadAndVerifyBLSKey(ctx, blsPubKey)
		if err != nil {
			return fmt.Errorf("failed to upload and verify BLS key: %w", err)
		}

		tx, err := utils.ExecuteTransactionWithRetry(cli.ethClient, cli.auth, func(auth *bind.TransactOpts) (*types.Transaction, error) {
			return cli.operatorsContract.SubmitBlsKeyCID(auth, cid)
		})
		if err != nil {
			return fmt.Errorf("submit BLS Public Key transaction couldn't be successful: %w", err)
		}

		if err := utils.WaitForTxConfirmationWithRetry(ctx, cli.ethClient, tx, 2); err != nil {
			return fmt.Errorf("failed to wait for BLS key submission confirmation: %w", err)
		}

		log.Println("BLS Public Key submitted successfully")
	} else {
		log.Println("BLS Public Key already submitted")
	}

	return cli.waitForBLSKeySubmission(ctx)
}

func (cli *Client) waitForBLSKeySubmission(ctx context.Context) error {

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			var submitted bool
			err := cli.executeRPCCallWithRetry(ctx, func() error {
				var err error
				submitted, err = cli.operatorsContract.HasSubmittedBLSKey(nil, cli.auth.From)
				return err
			})
			if err == nil && submitted {
				return nil
			}
			log.Println("Waiting for BLS key submission to be recognized...")

		}
	}
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
}

func (cli *Client) waitForAllOperators(ctx context.Context) error {
	var operatorCount *big.Int
	var err error
	for {
		err = cli.executeRPCCallWithRetry(ctx, func() error {
			operatorCount, err = cli.operatorsContract.GetOperatorCount(nil)
			return err
		})
		if err != nil {
			return fmt.Errorf("couldn't get the operator count: %w", err)
		}
		if operatorCount.Int64() >= int64(cli.committeeSize-1) {
			break
		}
		log.Printf("Waiting for all operators to be registered (%d/%d)...\n", operatorCount.Int64(), cli.committeeSize)
		time.Sleep(3 * time.Second)
	}

	var signers []handler.Signer
	for i := 0; i < int(operatorCount.Int64()); i++ {
		var operator struct {
			Operator     common.Address
			BlsPubKeyCID string
		}
		err = cli.executeRPCCallWithRetry(ctx, func() error {
			operator, err = cli.operatorsContract.Operators(nil, big.NewInt(int64(i)))
			return err
		})
		if err != nil {
			return fmt.Errorf("couldn't get the operators: %w", err)
		}

		var keyCID string
		err = cli.executeRPCCallWithRetry(ctx, func() error {
			keyCID, err = cli.operatorsContract.GetBLSPubKeyCIDByIndex(nil, big.NewInt(int64(i)))
			return err
		})
		if err != nil {
			return fmt.Errorf("couldn't get the BLS Pub Key by Index: %w", err)
		}

		key, err := cli.IPFSService.GetKeyByCID(keyCID)
		if err != nil {
			return fmt.Errorf("couldn't get the BLS Pub Key from IPFS: %w", err)
		}

		newSigner := handler.NewSigner(operator.Operator, key)
		signers = append(signers, newSigner)
	}

	cli.signers = &signers
	log.Printf("All operators registered successfully: %d\n", len(signers))

	return nil
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

	var ourIndex *big.Int
	err = cli.executeRPCCallWithRetry(ctx, func() error {
		var err error
		ourIndex, err = cli.operatorsContract.GetOperatorIndex(nil, cli.auth.From)
		return err
	})
	if err != nil {
		return fmt.Errorf("couldn't get our operator index: %w", err)
	}

	cli.Handler = handler.NewHandler(sub, topicHandle, cli.signers, cli.privKey, cli.rpcUrl, cli.apiPort, cli.committeeSize-1, ourIndex.Uint64(), cli.committeeSize/2)
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
