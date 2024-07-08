package mempool

import (
	"slices"
	"sync"

	"github.com/korayakpinar/p2pclient/src/types"
)

type Mempool struct {
	encryptedTxs     []*types.EncryptedTransaction
	decryptedTxs     []*types.DecryptedTransaction
	partDecRegistry  Registry[string, [][]byte]         // TxHash -> []PartialDecryption
	orderSigRegistry Registry[uint32, []types.OrderSig] // BlockNum -> []OrderSig
}

func NewMempool() *Mempool {
	partDecRegistry := NewRegistry[string, [][]byte]()
	orderSigRegistry := NewRegistry[uint32, []types.OrderSig]()
	return &Mempool{
		encryptedTxs:     []*types.EncryptedTransaction{},
		decryptedTxs:     []*types.DecryptedTransaction{},
		partDecRegistry:  partDecRegistry,
		orderSigRegistry: orderSigRegistry,
	}
}

func (m *Mempool) AddEncryptedTx(tx *types.EncryptedTransaction) {
	//TODO: Check if the transaction is already in the mempool
	m.encryptedTxs = append(m.encryptedTxs, tx)
}

func (m *Mempool) GetTransactions() []*types.EncryptedTransaction {
	return m.encryptedTxs
}

func (m *Mempool) GetTransaction(hash string) *types.EncryptedTransaction {
	for _, tx := range m.encryptedTxs {
		if tx.Header.Hash == hash {
			return tx
		}
	}
	return nil
}

func (m *Mempool) RemoveTransactions(txHashes []string) {
	newTxs := []*types.EncryptedTransaction{}
	for _, tx := range m.encryptedTxs {
		if !slices.Contains(txHashes, tx.Header.Hash) {
			newTxs = append(newTxs, tx)
		}
	}
	m.encryptedTxs = newTxs
}

func (m *Mempool) AddDecryptedTx(tx *types.DecryptedTransaction) {
	//TODO: Check if the transaction is already in the mempool
	m.decryptedTxs = append(m.decryptedTxs, tx)
}

func (m *Mempool) GetThreshold(hash string) uint32 {
	for _, tx := range m.encryptedTxs {
		if tx.Header.Hash == hash {
			return tx.Body.Threshold
		}
	}
	return 0
}

func (m *Mempool) AddPartialDecryption(hash string, partDec *[]byte) {
	//TODO: Check if the partial decryption is already in the mempool
	arr := m.partDecRegistry.Load(hash)
	*arr = append(*arr, *partDec)
}

func (m *Mempool) GetPartialDecryptions(hash string) *[][]byte {
	partDecs := m.partDecRegistry.Load(hash)
	return partDecs
}

func (m *Mempool) GetPartialDecryptionCount(hash string) uint32 {
	partDecs := (*m).partDecRegistry.Load(hash)
	return uint32(len(*partDecs))
}

func (m *Mempool) AddOrderSig(blockNum uint32, orderSig types.OrderSig) {
	//TODO: Check if the order signature is already in the mempool
	orderSigs := m.orderSigRegistry.Load(blockNum)
	*orderSigs = append(*orderSigs, orderSig)
}

type Registry[K comparable, V any] struct {
	cache *sync.Map
}

func NewRegistry[K comparable, V any]() Registry[K, V] {
	return Registry[K, V]{
		cache: &sync.Map{},
	}
}

func (m Registry[K, V]) Load(k K) *V {
	if val, ok := m.cache.Load(k); ok {
		return val.(*V)
	}
	return nil
}

func (m Registry[K, V]) Store(k K, v *V) {
	m.cache.Store(k, v)
}

func (m *Registry[K, V]) Delete(k K) {
	m.cache.Delete(k)
}
