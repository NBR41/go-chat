package crypto

import (
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	crypt, err := NewCrypto("fooofooofooofooo", "bar")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	var txt, enc, enc2, dec string
	txt = "some text"
	enc, err = crypt.Encrypt(txt)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	t.Log(enc)

	enc2, err = crypt.Encrypt(txt)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	t.Log(enc2)

	if enc == enc2 {
		t.Fatal("same encrypted text")
	}

	dec, err = crypt.Decrypt(enc)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if dec != txt {
		t.Fatalf("unexpected output: %s", dec)
	}
	dec, err = crypt.Decrypt(enc2)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if dec != txt {
		t.Fatalf("unexpected output: %s", dec)
	}
}

func TestHash(t *testing.T) {
	crypt, err := NewCrypto("fooofooofooofooo", "bar")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	var h, h2, h3 string
	h, err = crypt.Hash("password")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	t.Log(h)

	h2, err = crypt.Hash("password")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	t.Log(h2)
	if h != h2 {
		t.Fatal("hashes are differents")
	}

	h3, err = crypt.Hash("foo")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	t.Log(h3)

	if h == h3 {
		t.Fatal("hashes are equals")
	}
}
