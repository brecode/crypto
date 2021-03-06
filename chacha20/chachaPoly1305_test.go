// Use of this source code is governed by a license
// that can be found in the LICENSE file.

package chacha20

import "testing"

var recFunc = func(t *testing.T, msg string) {
	if recover() == nil {
		t.Fatalf("Expected error: %s", msg)
	}
}

func TestNewChaCha20Poly1305WithTagSize(t *testing.T) {
	var key [32]byte
	_, err := NewChaCha20Poly1305WithTagSize(&key, 0)
	if err == nil {
		t.Fatalf("NewChaCha20Poly1305WithTagSize accepted invalid tagsize: %d", 0)
	}

	_, err = NewChaCha20Poly1305WithTagSize(&key, 17)
	if err == nil {
		t.Fatalf("NewChaCha20Poly1305WithTagSize accepted invalid tagsize: %d", 0)
	}
}

func TestOverhead(t *testing.T) {
	var key [32]byte
	c := NewChaCha20Poly1305(&key)

	if o := c.Overhead(); o != TagSize {
		t.Fatalf("Expected %d but Overhead() returned %d", TagSize, o)
	}

	c, err := NewChaCha20Poly1305WithTagSize(&key, 12)
	if err != nil {
		t.Fatalf("Failed to create ChaCha20Poly1305 instance: %s", err)
	}
	if o := c.Overhead(); o != 12 {
		t.Fatalf("Expected %d but Overhead() returned %d", 12, o)
	}
}

func TestNonceSize(t *testing.T) {
	var key [32]byte
	c := NewChaCha20Poly1305(&key)
	if n := c.NonceSize(); n != NonceSize {
		t.Fatalf("Expected %d but NonceSize() returned %d", TagSize, n)
	}
}

func TestSeal(t *testing.T) {
	var key [32]byte
	c := NewChaCha20Poly1305(&key)

	var (
		nonce [NonceSize]byte
		src   [64]byte
		dst   [64 + TagSize]byte
	)

	mustFail := func(msg string, dst, nonce, src []byte) {
		defer recFunc(t, msg)
		c.Seal(dst, nonce, src, nil)
	}

	mustFail("nonce size is invalid", dst[:], nonce[:NonceSize-1], src[:])

	mustFail("dst length invalid", dst[:len(dst)-2], nonce[:], src[:])
}

func TestOpen(t *testing.T) {
	var key [32]byte
	c := NewChaCha20Poly1305(&key)

	var (
		nonce [NonceSize]byte
		src   [64]byte
		dst   [64 + TagSize]byte
	)

	_, err := c.Open(dst[:], nonce[:NonceSize-1], src[:], nil)
	if err == nil {
		t.Fatal("Open() accepted invalid nonce size")
	}

	_, err = c.Open(dst[:], nonce[:], src[:TagSize-1], nil)
	if err == nil {
		t.Fatal("Open() accepted invalid ciphertext length")
	}

	mustFail := func(msg string, dst, nonce, src []byte) {
		defer recFunc(t, msg)
		c.Open(dst, nonce, src, nil)
	}

	mustFail("dst length invalid", dst[:len(src)-TagSize-1], nonce[:], src[:])

	// Check tag verification
	c.Seal(dst[:], nonce[:], src[:], nil)
	dst[len(src)+1] += 1 // modify tag

	_, err = c.Open(src[:], nonce[:], dst[:], nil)
	if err == nil {
		t.Fatal("Open() accepted invalid auth. tag")
	}
}

func BenchmarkSeal64B(b *testing.B) {
	var key [32]byte
	var nonce [12]byte
	c := NewChaCha20Poly1305(&key)

	msg := make([]byte, 64)
	dst := make([]byte, len(msg)+TagSize)
	data := make([]byte, 32)

	b.SetBytes(int64(len(msg)))
	for i := 0; i < b.N; i++ {
		dst = c.Seal(dst, nonce[:], msg, data)
	}
}

func BenchmarkSeal1K(b *testing.B) {
	var key [32]byte
	var nonce [12]byte
	c := NewChaCha20Poly1305(&key)

	msg := make([]byte, 1024)
	dst := make([]byte, len(msg)+TagSize)
	data := make([]byte, 32)

	b.SetBytes(int64(len(msg)))
	for i := 0; i < b.N; i++ {
		dst = c.Seal(dst, nonce[:], msg, data)
	}
}

func BenchmarkSeal64K(b *testing.B) {
	var key [32]byte
	var nonce [12]byte
	c := NewChaCha20Poly1305(&key)

	msg := make([]byte, 64*1024)
	dst := make([]byte, len(msg)+TagSize)
	data := make([]byte, 32)

	b.SetBytes(int64(len(msg)))
	for i := 0; i < b.N; i++ {
		dst = c.Seal(dst, nonce[:], msg, data)
	}
}

func BenchmarkOpen64B(b *testing.B) {
	var key [32]byte
	var nonce [12]byte
	c := NewChaCha20Poly1305(&key)

	msg := make([]byte, 64)
	dst := make([]byte, len(msg))
	ciphertext := make([]byte, len(msg)+TagSize)
	data := make([]byte, 32)
	ciphertext = c.Seal(ciphertext, nonce[:], msg, data)

	b.SetBytes(int64(len(msg)))
	for i := 0; i < b.N; i++ {
		dst, _ = c.Open(dst, nonce[:], ciphertext, data)
	}
}

func BenchmarkOpen1K(b *testing.B) {
	var key [32]byte
	var nonce [12]byte
	c := NewChaCha20Poly1305(&key)

	msg := make([]byte, 1024)
	dst := make([]byte, len(msg))
	ciphertext := make([]byte, len(msg)+TagSize)
	data := make([]byte, 32)
	ciphertext = c.Seal(ciphertext, nonce[:], msg, data)

	b.SetBytes(int64(len(msg)))
	for i := 0; i < b.N; i++ {
		dst, _ = c.Open(dst, nonce[:], ciphertext, data)
	}
}

func BenchmarkOpen64K(b *testing.B) {
	var key [32]byte
	var nonce [12]byte
	c := NewChaCha20Poly1305(&key)

	msg := make([]byte, 64*1024)
	dst := make([]byte, len(msg))
	ciphertext := make([]byte, len(msg)+TagSize)
	data := make([]byte, 32)
	ciphertext = c.Seal(ciphertext, nonce[:], msg, data)

	b.SetBytes(int64(len(msg)))
	for i := 0; i < b.N; i++ {
		dst, _ = c.Open(dst, nonce[:], ciphertext, data)
	}
}
