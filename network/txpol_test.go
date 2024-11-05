package network

import (
	"github.com/peng9808/BlockchainX/core"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"strconv"
	"testing"
)

func TestNewTxPool(t *testing.T) {
	p := NewTxPool()
	assert.Equal(t, p.Len(), 0)
}

func TestTxPoolAddTx(t *testing.T) {
	p := NewTxPool()
	tx := core.NewTransaction([]byte("fooo"))
	assert.Nil(t, p.Add(tx))
	assert.Equal(t, p.Len(), 1)

	p.Flush()
	assert.Equal(t, p.Len(), 0)
}

func TestSortTransactions(t *testing.T) {
	p := NewTxPool()
	txLen := 1000

	for i := 0; i < txLen; i++ {
		tx := core.NewTransaction([]byte(strconv.FormatInt(int64(i), 10)))
		tx.SetFirstSeen(int64(i * rand.Intn(10000)))
		assert.Nil(t, p.Add(tx))
	}

	assert.Equal(t, txLen, p.Len())

	txx := p.Transactions()
	for i := 0; i < len(txx)-1; i++ {
		assert.True(t, txx[i].FirstSeen() <= txx[i+1].FirstSeen())
	}
}
