package core

import (
	"fmt"
	"github.com/peng9808/BlockchainX/types"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBlockHash(t *testing.T) {
	b := &Block{
		Header: Header{
			Version:   1,
			PrevBlock: types.RandomHash(),
			Timestamp: time.Now().UnixNano(),
			Height:    10,
		},
		Transactions: []Transaction{},
	}

	h := b.Hash()
	fmt.Println(h)
	assert.False(t, h.IsZero())
}
