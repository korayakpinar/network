package mempool

import (
	"bytes"
	"sync"
	"sync/atomic"

	"github.com/korayakpinar/network/src/types"
)

type Mempool struct {
	encryptedTxs     sync.Map
	decryptedTxs     sync.Map
	encryptedTxCount int64
	decryptedTxCount int64
	partDecRegistry  Registry[string, [][]byte]
	orderSigRegistry Registry[uint64, []types.OrderSig]
}

func NewMempool() *Mempool {
	return &Mempool{
		encryptedTxs:     sync.Map{},
		decryptedTxs:     sync.Map{},
		encryptedTxCount: 0,
		decryptedTxCount: 0,
		partDecRegistry:  NewRegistry[string, [][]byte](),
		orderSigRegistry: NewRegistry[uint64, []types.OrderSig](),
	}
}

// Encrypted Transactions methods

func (m *Mempool) AddEncryptedTx(tx *types.EncryptedTransaction) {
	if _, loaded := m.encryptedTxs.LoadOrStore(tx.Header.Hash, tx); !loaded {
		atomic.AddInt64(&m.encryptedTxCount, 1)
	}
}

func (m *Mempool) GetEncryptedTransactions() []*types.EncryptedTransaction {
	txs := make([]*types.EncryptedTransaction, 0, atomic.LoadInt64(&m.encryptedTxCount))
	m.encryptedTxs.Range(func(_, value interface{}) bool {
		txs = append(txs, value.(*types.EncryptedTransaction))
		return true
	})
	return txs
}

func (m *Mempool) GetEncryptedTransaction(hash string) *types.EncryptedTransaction {
	if tx, ok := m.encryptedTxs.Load(hash); ok {
		return tx.(*types.EncryptedTransaction)
	}
	return nil
}

func (m *Mempool) RemoveEncryptedTransactions(txHashes []string) {
	for _, hash := range txHashes {
		if _, loaded := m.encryptedTxs.LoadAndDelete(hash); loaded {
			atomic.AddInt64(&m.encryptedTxCount, -1)
		}
	}
}

func (m *Mempool) GetEncryptedTransactionCount() int {
	return int(atomic.LoadInt64(&m.encryptedTxCount))
}

func (m *Mempool) GetThreshold(hash string) uint64 {
	if tx, ok := m.encryptedTxs.Load(hash); ok {
		return tx.(*types.EncryptedTransaction).Body.Threshold
	}
	return 0
}

// Decrypted Transactions methods

func (m *Mempool) AddDecryptedTx(tx *types.DecryptedTransaction) {
	if _, loaded := m.decryptedTxs.LoadOrStore(tx.Header.Hash, tx); !loaded {
		atomic.AddInt64(&m.decryptedTxCount, 1)
	}
}

func (m *Mempool) GetDecryptedTransactions() []*types.DecryptedTransaction {
	txs := make([]*types.DecryptedTransaction, 0, atomic.LoadInt64(&m.decryptedTxCount))
	m.decryptedTxs.Range(func(_, value interface{}) bool {
		txs = append(txs, value.(*types.DecryptedTransaction))
		return true
	})
	return txs
}

func (m *Mempool) GetDecryptedTransaction(hash string) *types.DecryptedTransaction {
	if tx, ok := m.decryptedTxs.Load(hash); ok {
		return tx.(*types.DecryptedTransaction)
	}
	return nil
}

func (m *Mempool) RemoveDecryptedTransactions(txHashes []string) {
	for _, hash := range txHashes {
		if _, loaded := m.decryptedTxs.LoadAndDelete(hash); loaded {
			atomic.AddInt64(&m.decryptedTxCount, -1)
		}
	}
}

func (m *Mempool) GetDecryptedTransactionCount() int {
	return int(atomic.LoadInt64(&m.decryptedTxCount))
}

// Partial Decryption methods

func (m *Mempool) AddPartialDecryption(hash string, partDec *[]byte) {
	arr := m.partDecRegistry.Load(hash)
	if arr == nil {
		newArr := [][]byte{*partDec}
		m.partDecRegistry.Store(hash, &newArr)
	} else {
		for _, existingPartDec := range *arr {
			if bytes.Equal(existingPartDec, *partDec) {
				return // Already exists, don't add
			}
		}
		newArr := append(*arr, *partDec)
		m.partDecRegistry.Store(hash, &newArr)
	}
}

func (m *Mempool) GetPartialDecryptions(hash string) *[][]byte {
	return m.partDecRegistry.Load(hash)
}

func (m *Mempool) GetPartialDecryptionCount(hash string) uint64 {
	partDecs := m.partDecRegistry.Load(hash)
	if partDecs == nil {
		return 0
	}
	return uint64(len(*partDecs))
}

// Order Signature methods

func (m *Mempool) AddOrderSig(blockNum uint64, orderSig types.OrderSig) {
	orderSigs := m.orderSigRegistry.Load(blockNum)
	if orderSigs == nil {
		m.orderSigRegistry.Store(blockNum, &[]types.OrderSig{orderSig})
	} else {
		newOrderSigs := append(*orderSigs, orderSig)
		m.orderSigRegistry.Store(blockNum, &newOrderSigs)
	}
}

func (m *Mempool) GetOrderSigs(blockNum uint64) *[]types.OrderSig {
	return m.orderSigRegistry.Load(blockNum)
}

// Registry type and methods

type Registry[K comparable, V any] struct {
	cache *sync.Map
}

func NewRegistry[K comparable, V any]() Registry[K, V] {
	return Registry[K, V]{
		cache: &sync.Map{},
	}
}

func (r Registry[K, V]) Load(k K) *V {
	if val, ok := r.cache.Load(k); ok {
		return val.(*V)
	}
	return nil
}

func (r Registry[K, V]) Store(k K, v *V) {
	r.cache.Store(k, v)
}

func (r *Registry[K, V]) Delete(k K) {
	r.cache.Delete(k)
}
