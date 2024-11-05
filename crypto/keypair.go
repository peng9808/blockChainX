package crypto

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"github.com/peng9808/BlockchainX/types"
	"math/big"
)

type PrivateKey struct {
	key *ecdsa.PrivateKey
}

type Signature struct {
	S, R *big.Int
}

type PublicKey struct {
	Key *ecdsa.PublicKey
}

func (pk *PublicKey) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	curveName := pk.Key.Curve.Params().Name
	if err := enc.Encode(curveName); err != nil {
		return nil, err
	}

	if err := enc.Encode(pk.Key.X); err != nil {
		return nil, err
	}
	if err := enc.Encode(pk.Key.Y); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (pk *PublicKey) GobDecode(buf []byte) error {
	dec := gob.NewDecoder(bytes.NewReader(buf))

	var curveName string
	if err := dec.Decode(&curveName); err != nil {
		return err
	}

	curve := elliptic.P256()
	if curveName != curve.Params().Name {
		return fmt.Errorf("not support: %s", curveName)
	}

	var x, y *big.Int
	if err := dec.Decode(&x); err != nil {
		return err
	}
	if err := dec.Decode(&y); err != nil {
		return err
	}

	pk.Key = &ecdsa.PublicKey{
		Curve: curve,
		X:     x,
		Y:     y,
	}
	return nil
}

func (k *PrivateKey) Sign(data []byte) (*Signature, error) {
	r, s, err := ecdsa.Sign(rand.Reader, k.key, data)
	if err != nil {
		return nil, err
	}
	return &Signature{
		R: r,
		S: s,
	}, nil
}

func (sig Signature) Verify(pubKey PublicKey, data []byte) bool {
	return ecdsa.Verify(pubKey.Key, data, sig.R, sig.S)
}

func GeneratePrivateKey() PrivateKey {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	return PrivateKey{
		key: key,
	}
}

func (k PrivateKey) PublicKey() PublicKey {
	return PublicKey{
		Key: &k.key.PublicKey,
	}
}

func (k PublicKey) ToSlice() []byte {
	return elliptic.MarshalCompressed(k.Key, k.Key.X, k.Key.Y)
}

func (k PublicKey) Address() types.Address {
	h := sha256.Sum256(k.ToSlice())

	return types.AddressFromBytes(h[len(h)-20:])
}
