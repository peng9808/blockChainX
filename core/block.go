package core

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/peng9808/BlockchainX/crypto"
	"github.com/peng9808/BlockchainX/types"
)

type Header struct {
	Version       uint32
	PrevBlockHash types.Hash
	Timestamp     int64
	Height        uint32
	Nonce         uint64
}

type Block struct {
	*Header
	Transactions []Transaction
	Validator    crypto.PublicKey
	Signature    *crypto.Signature
	//Cached of the header hash
	hash types.Hash
}

func NewBlock(h *Header, txx []Transaction) *Block {
	return &Block{
		Header:       h,
		Transactions: txx,
	}
}

func (b *Block) Sign(privKey crypto.PrivateKey) error {

	sign, err := privKey.Sign(b.HeaderData())
	if err != nil {
		return err
	}
	b.Validator = privKey.PublicKey()
	b.Signature = sign
	return nil
}

func (b *Block) Verify() error {
	if b.Signature == nil {
		return fmt.Errorf("block has no signature")
	}
	if !b.Signature.Verify(b.Validator, b.HeaderData()) {
		return fmt.Errorf("block has invalid signature")
	}
	return nil
}

func (b *Block) HeaderData() []byte {
	buf := &bytes.Buffer{}
	encoder := gob.NewEncoder(buf)
	encoder.Encode(b.Header)
	return buf.Bytes()
}

func (b *Block) Hash(hasher Hasher[*Block]) types.Hash {

	if b.hash.IsZero() {
		b.hash = hasher.Hash(b)
	}
	return b.hash
}
