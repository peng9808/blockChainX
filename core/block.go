package core

import (
	"bytes"
	"crypto/sha256"
	"github.com/peng9808/BlockchainX/types"
)

type Block struct {
	Header
	Transactions []Transaction

	//Cached of the header hash
	hash types.Hash
}

func (b *Block) Hash() types.Hash {

	buf := &bytes.Buffer{}
	b.Header.EncodeBinary(buf)
	if b.hash.IsZero() {
		b.hash = types.Hash(sha256.Sum256(buf.Bytes()))
	}

	return b.hash
}
