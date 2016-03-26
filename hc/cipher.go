// Use of this source code is governed by a license
// that can be found in the LICENSE file.

// The hc package implements the both stream ciphers
// HC-128 and HC-256 from the eSTREAM portfolio (software).
// Both ciphers were designed by Hongjun Wu.
// HC-128 and HC-256 are fast stream ciphers after[!]
// initialization. This may not be an issue, but if keys
// change too often, both ciphers spend a lot of time in
// initialization. In this case another cipher may perform
// better.
package hc

import (
	"crypto/cipher"
	"github.com/EncEve/crypto"
)

const (
	mod512  = 0x1FF
	mod1024 = 0x3FF
	mod2048 = 0x7FF
)

// A hc128 holds the both states P and Q, the counter,
// 4 byte of the keystream and the offset
type hc128 struct {
	p, q        []uint32
	ctr, stream uint32
	off         uint
}

// A hc256 holds the both states P and Q, the counter,
// 4 byte of the keystream and the offset
type hc256 struct {
	p, q        []uint32
	ctr, stream uint32
	off         uint
}

// New128 creates and returns a new cipher.Stream.
// The key argument must be 128 bit (16 byte).
// The nonce argument must be at least 128 bit (16 byte).
// The returned cipher.Stream implements the HC128 cipher.
// If the key is not 128 bit or the nonce is not 128 bit,
// this function returns an non-nil error.
func New128(key, nonce []byte) (cipher.Stream, error) {
	if k := len(key); k != 16 {
		return nil, crypto.KeySizeError(k)
	}
	if n := len(nonce); n != 16 {
		return nil, crypto.NonceSizeError(n)
	}
	c := &hc128{
		p:      make([]uint32, 512),
		q:      make([]uint32, 512),
		off:    4,
		ctr:    0,
		stream: 0,
	}
	c.initialize(key, nonce)

	return c, nil
}

// New256 creates and returns a new cipher.Stream.
// The key argument must be 256 bit (32 byte),
// The nonce argument must be at least 256 bit (32 byte),
// The returned cipher.Stream implements the HC256 cipher.
// If the key is not 256 bit or the nonce is 256 bit,
// this function returns an non-nil error.
func New256(key, nonce []byte) (cipher.Stream, error) {
	if k := len(key); k != 32 {
		return nil, crypto.KeySizeError(k)
	}
	if n := len(nonce); n != 32 {
		return nil, crypto.NonceSizeError(n)
	}
	c := &hc256{
		p:      make([]uint32, 1024),
		q:      make([]uint32, 1024),
		off:    4,
		ctr:    0,
		stream: 0,
	}
	c.initialize(key, nonce)

	return c, nil
}
