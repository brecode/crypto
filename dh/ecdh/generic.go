// Use of this source code is governed by a license
// that can be found in the LICENSE file.

package ecdh

import (
	"crypto/elliptic"
	cryptorand "crypto/rand"
	"errors"
	"io"
	big "github.com/ncw/gmp"
)

// The same unmarshal as elliptic.Unmarshal but without
// point checking (see Check method)
func unmarshal(curve elliptic.Curve, data []byte) (x, y *big.Int) {
	byteLen := (curve.Params().BitSize + 7) >> 3
	if len(data) != 1+2*byteLen {
		return
	}
	if data[0] != 4 { // uncompressed form
		return
	}
	x = new(big.Int).SetBytes(data[1 : 1+byteLen])
	y = new(big.Int).SetBytes(data[1+byteLen:])
	return
}

// GenericCurve creates a new ecdh.KeyExchange with
// generic elliptic.Curve implementations.
func GenericCurve(c elliptic.Curve) KeyExchange {
	if c == nil {
		panic("ecdh: curve is nil")
	}
	return genericCurve{curve: c}
}

type genericCurve struct {
	curve elliptic.Curve
}

func (g genericCurve) GenerateKey(rand io.Reader) (private PrivateKey, public PublicKey, err error) {
	if rand == nil {
		rand = cryptorand.Reader
	}
	private, x, y, err := elliptic.GenerateKey(g.curve, rand)
	if err != nil {
		private = nil
		return
	}
	public = elliptic.Marshal(g.curve, x, y)
	return
}

func (g genericCurve) PublicKey(private PrivateKey) (public PublicKey) {
	N := g.curve.Params().N

	if new(big.Int).SetBytes(private).Cmp(N) >= 0 {
		panic("ecdh: private key cannot used with given curve")
	}

	x, y := g.curve.ScalarBaseMult(private)
	public = elliptic.Marshal(g.curve, x, y)
	return
}

func (g genericCurve) Check(peersPublic PublicKey) (err error) {
	x, y := unmarshal(g.curve, peersPublic)
	if !g.curve.IsOnCurve(x, y) {
		err = errors.New("peer's public key is not on curve")
	}
	return
}

func (g genericCurve) ComputeSecret(private PrivateKey, peersPublic PublicKey) (secret []byte) {
	x, y := unmarshal(g.curve, peersPublic)

	sX, _ := g.curve.ScalarMult(x, y, private)

	secret = sX.Bytes()
	return
}
