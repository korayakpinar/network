package mempool

import (
	"sync"

	"github.com/korayakpinar/p2pclient/src/types"
)

type Mempool struct {
	txs              *[]*types.EncryptedTransaction
	partDecRegistry  *Registry[types.TxHash, types.PartialDecryption]
	orderSigRegistry *Registry[types.BlockNum, types.OrderSig]
}

func NewMempool() *Mempool {
	partDecRegistry := NewRegistry[types.TxHash, types.PartialDecryption]()
	orderSigRegistry := NewRegistry[types.BlockNum, types.OrderSig]()
	return &Mempool{
		txs:              &[]*types.EncryptedTransaction{},
		partDecRegistry:  &partDecRegistry,
		orderSigRegistry: &orderSigRegistry,
	}
}

func (m *Mempool) AddTransaction(tx *types.EncryptedTransaction) {
	*m.txs = append(*m.txs, tx)
}

func (m *Mempool) AddPartialDecryption(hash types.TxHash, partDec types.PartialDecryption) {
	m.partDecRegistry.Store(hash, &partDec)
}

func (m *Mempool) AddOrderSig(blockNum types.BlockNum, orderSig types.OrderSig) {
	m.orderSigRegistry.Store(blockNum, &orderSig)
}

type Registry[K comparable, V any] struct {
	cache *sync.Map
}

func NewRegistry[K comparable, V any]() Registry[K, V] {
	return Registry[K, V]{
		cache: &sync.Map{},
	}
}

func (m Registry[K, V]) Load(k K) (*V, bool) {
	if val, ok := m.cache.Load(k); ok {
		return val.(*V), true
	}
	return nil, false
}

func (m Registry[K, V]) Store(k K, v *V) {
	m.cache.Store(k, v)
}

func (m *Registry[K, V]) Delete(k K) {
	m.cache.Delete(k)
}
