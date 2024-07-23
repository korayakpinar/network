package mempool

import (
	"bytes"
	"sync"
	"sync/atomic"

	"github.com/korayakpinar/network/src/types"
)

type Mempool struct {
	encryptedTxs     sync.Map
	includedTxs      sync.Map
	decryptedTxs     sync.Map
	encryptedTxCount int64
	includedTxCount  int64
	decryptedTxCount int64
	partDecRegistry  Registry[string, map[uint64][]byte]
	orderSigRegistry Registry[uint64, []types.OrderSig]
}

func NewMempool() *Mempool {
	return &Mempool{
		encryptedTxs:     sync.Map{},
		decryptedTxs:     sync.Map{},
		encryptedTxCount: 0,
		includedTxCount:  0,
		decryptedTxCount: 0,
		partDecRegistry:  NewRegistry[string, map[uint64][]byte](),
		orderSigRegistry: NewRegistry[uint64, []types.OrderSig](),
	}
}

// Encrypted Transactions methods

func (m *Mempool) AddEncryptedTx(tx *types.EncryptedTransaction) {
	if _, loaded := m.encryptedTxs.LoadOrStore(tx.Header.Hash, tx); !loaded {
		atomic.AddInt64(&m.encryptedTxCount, 1)
	}
}

func (m *Mempool) IncludeEncryptedTxs(txHashes []string) {
	for _, hash := range txHashes {
		if tx, ok := m.encryptedTxs.Load(hash); ok {
			m.encryptedTxs.Delete(hash)
			atomic.AddInt64(&m.encryptedTxCount, -1)
			m.includedTxs.Store(hash, tx)
		}
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

// Included Transactions methods

func (m *Mempool) AddIncludedTx(tx *types.EncryptedTransaction) {
	if _, loaded := m.includedTxs.LoadOrStore(tx.Header.Hash, tx); !loaded {
		atomic.AddInt64(&m.includedTxCount, 1)
	}
}

func (m *Mempool) GetIncludedTransactions() []*types.EncryptedTransaction {
	txs := make([]*types.EncryptedTransaction, 0, atomic.LoadInt64(&m.includedTxCount))
	m.includedTxs.Range(func(_, value interface{}) bool {
		txs = append(txs, value.(*types.EncryptedTransaction))
		return true
	})
	return txs
}

func (m *Mempool) GetIncludedTransaction(hash string) *types.EncryptedTransaction {
	if tx, ok := m.includedTxs.Load(hash); ok {
		return tx.(*types.EncryptedTransaction)
	}
	return nil
}

func (m *Mempool) RemoveIncludedTransactions(txHashes []string) {
	for _, hash := range txHashes {
		if _, loaded := m.includedTxs.LoadAndDelete(hash); loaded {
			atomic.AddInt64(&m.includedTxCount, -1)
		}
	}
}

func (m *Mempool) GetIncludedTransactionCount() int {
	return int(atomic.LoadInt64(&m.includedTxCount))
}

func (m *Mempool) GetThreshold(hash string) uint64 {
	if tx, ok := m.includedTxs.Load(hash); ok {
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

func (m *Mempool) AddPartialDecryption(hash string, index uint64, partDec []byte) {
	partMapPtr := m.partDecRegistry.Load(hash)
	if partMapPtr == nil {
		// If no map exists for this hash, create a new one
		newMap := make(map[uint64][]byte)
		newMap[index] = partDec
		m.partDecRegistry.Store(hash, &newMap)
	} else {
		partMap := *partMapPtr
		// If a map already exists, check if this index already has a value
		if existingPartDec, exists := partMap[index]; exists {
			if bytes.Equal(existingPartDec, partDec) {
				return // Already exists, don't add
			}
		}
		// Add the partial decryption for this index
		partMap[index] = partDec
		m.partDecRegistry.Store(hash, &partMap)
	}
}

func (m *Mempool) GetPartialDecryptions(hash string) map[uint64][]byte {
	partMapPtr := m.partDecRegistry.Load(hash)
	if partMapPtr == nil {
		return nil
	}
	return *partMapPtr
}

func (m *Mempool) GetPartialDecryptionCount(hash string) int {
	partMapPtr := m.partDecRegistry.Load(hash)
	if partMapPtr == nil {
		return 0
	}
	return len(*partMapPtr)
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
