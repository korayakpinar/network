package utils

import (
	"github.com/BurntSushi/toml"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	crypto "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
)

type Config struct {
	TopicName string `toml:"topicName"`
	PrivKey   string `toml:"privKey"`
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

func IdFromPubKey(privKey string) (*peer.ID, error) {
	publicKeyECDSA, err := ethCrypto.HexToECDSA(privKey)
	if err != nil {
		return nil, err
	}

	pubKeyBytes := ethCrypto.FromECDSAPub(&publicKeyECDSA.PublicKey)
	pub, err := crypto.UnmarshalSecp256k1PublicKey(pubKeyBytes)
	if err != nil {
		return nil, err
	}

	id, err := peer.IDFromPublicKey(pub)
	if err != nil {
		return nil, err
	}

	return &id, nil
}

func IdToEthAddress(id peer.ID) (string, error) {
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

	ethAddress := ethCrypto.PubkeyToAddress(*etherPubKey).Hex()
	return ethAddress, nil
}

func IsOperator(id peer.ID) bool {
	addr, err := IdToEthAddress(id)
	if err != nil {
		panic(err)
	}

	// check if the address is staked enough amount
	_ = (addr) // just to remove the error
	return true
}

func LoadConfig(path string) (Config, error) {
	var cfg Config
	_, err := toml.DecodeFile(path, &cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}
