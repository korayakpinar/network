package utils

import (
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/BurntSushi/toml"
	"github.com/ethereum/go-ethereum/common"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"

	crypto "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
)

type Config struct {
	TopicName string `toml:"topicName"`
	PrivKey   string `toml:"privKey"`
	RpcURL    string `toml:"rpc"`
	Port      string `toml:"port"`
}

// ExecuteTransactionWithRetry attempts to execute a transaction with retry logic
func ExecuteTransactionWithRetry(client *ethclient.Client, auth *bind.TransactOpts, txFunc func(*bind.TransactOpts) (*types.Transaction, error)) (*types.Transaction, error) {
	var tx *types.Transaction
	var err error

	for i := 0; i < 5; i++ { // 5 attempts
		// Check the current nonce and set it in the auth
		nonce, err := client.PendingNonceAt(context.Background(), auth.From)
		if err != nil {
			return nil, fmt.Errorf("failed to get nonce: %w", err)
		}
		auth.Nonce = big.NewInt(int64(nonce))

		// Get the suggested gas price
		gasPrice, err := client.SuggestGasPrice(context.Background())
		if err != nil {
			return nil, fmt.Errorf("failed to get gas price: %w", err)
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

		fmt.Printf("Transaction failed, retrying with updated nonce and higher gas price... (Attempt %d)\n", i+1)
		time.Sleep(time.Second * 2) // Short wait before retry
	}

	return nil, fmt.Errorf("failed to execute transaction after multiple attempts: %w", err)
}

func WaitForTxConfirmationWithRetry(ctx context.Context, client *ethclient.Client, tx *types.Transaction, confirmations uint64) error {
	minDelay := time.Second
	maxDelay := time.Second * 10
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			err := WaitForTxConfirmation(ctx, client, tx, confirmations)
			if err == nil {
				return nil
			}
			if strings.Contains(err.Error(), "transaction failed") {
				return err // If the transaction failed, don't retry
			}
			fmt.Printf("Error waiting for transaction confirmation: %v. Retrying...\n", err)
			delay := time.Duration(rand.Int63n(int64(maxDelay-minDelay))) + minDelay
			time.Sleep(delay)

			// Increase the delay range for the next attempt, up to a maximum
			if maxDelay < time.Minute {
				maxDelay += time.Second * 5
			}
			if minDelay < time.Second*30 {
				minDelay += time.Second
			}
		}
	}
}

// WaitForTxConfirmation waits for a transaction to be mined and confirmed
func WaitForTxConfirmation(ctx context.Context, client *ethclient.Client, tx *types.Transaction, confirmations uint64) error {
	receipt, err := waitForTxReceipt(ctx, client, tx.Hash(), 2*time.Minute)
	if err != nil {
		return fmt.Errorf("error waiting for transaction receipt: %w", err)
	}

	if receipt.Status == 0 {
		return fmt.Errorf("transaction failed")
	}

	fmt.Printf("Transaction included in block %d\n", receipt.BlockNumber.Uint64())

	currentBlock, err := client.BlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("error getting current block number: %w", err)
	}

	for currentBlock-receipt.BlockNumber.Uint64() < confirmations {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(15 * time.Second):
			currentBlock, err = client.BlockNumber(ctx)
			if err != nil {
				return fmt.Errorf("error getting current block number: %w", err)
			}
			fmt.Printf("Waiting for confirmations... Current: %d, Target: %d\n",
				currentBlock-receipt.BlockNumber.Uint64(), confirmations)
		}
	}

	fmt.Printf("Transaction confirmed with %d block confirmations\n", currentBlock-receipt.BlockNumber.Uint64())
	return nil
}

func waitForTxReceipt(ctx context.Context, client *ethclient.Client, txHash common.Hash, timeout time.Duration) (*types.Receipt, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		receipt, err := client.TransactionReceipt(ctx, txHash)
		if receipt != nil {
			return receipt, nil
		}
		if err != nil {
			if err != ethereum.NotFound {
				return nil, fmt.Errorf("error getting transaction receipt: %w", err)
			}
		}

		select {
		case <-timeoutCtx.Done():
			return nil, fmt.Errorf("timeout waiting for transaction receipt")
		case <-ticker.C:
			continue
		}
	}
}

func EthToLibp2pPrivKey(key string) (crypto.PrivKey, error) {
	privKey, err := ethCrypto.HexToECDSA(key)
	if err != nil {
		return nil, err
	}

	privKeyBytes := ethCrypto.FromECDSA(privKey)

	priv, err := crypto.UnmarshalSecp256k1PrivateKey(privKeyBytes)
	if err != nil {
		return nil, err
	}

	return priv, nil
}

func IdToEthAddress(id peer.ID) common.Address {
	pubKey, err := id.ExtractPublicKey()
	if err != nil {
		panic(err)
	}

	rawBytes, err := pubKey.Raw()
	if err != nil {
		panic(err)
	}

	etherPubKey, err := ethCrypto.DecompressPubkey(rawBytes)
	if err != nil {
		panic(err)
	}

	ethAddress := ethCrypto.PubkeyToAddress(*etherPubKey)
	return ethAddress
}

func LoadConfig(path string) (Config, error) {
	var cfg Config
	_, err := toml.DecodeFile(path, &cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}
