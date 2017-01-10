package crypto

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"

	"golang.org/x/crypto/scrypt"
)

// Crypto struct for crypto tools
type Crypto struct {
	block cipher.Block
	salt  []byte
}

// NewCrypto return new instance of Crypto
// The key argument should be the AES key, either 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256.
func NewCrypto(key, salt string) (*Crypto, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	return &Crypto{block, []byte(salt)}, nil
}

// Encrypt returns encrypted text
func (c *Crypto) Encrypt(text string) (string, error) {
	plaintext := []byte(text)

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	stream := cipher.NewCFBEncrypter(c.block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)
	return hex.EncodeToString(ciphertext), nil
}

// Decrypt returns decrypted text cripted by Encrypt
func (c *Crypto) Decrypt(text string) (string, error) {
	ciphertext, err := hex.DecodeString(text)
	if err != nil {
		return "", err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(c.block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)
	return fmt.Sprintf("%s", ciphertext), nil
}

// Hash return derived key of text
func (c *Crypto) Hash(text string) (string, error) {
	dk, err := scrypt.Key([]byte(text), c.salt, 16384, 8, 1, 32)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", dk), nil
}

// SaveEncryptedFile encrypts file and save it to disk
func SaveEncryptedFile(src multipart.File, destFile, key, destPath string) error {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return err
	}
	var iv = make([]byte, aes.BlockSize)
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return err
	}
	stream := cipher.NewCFBEncrypter(block, iv)

	var file *os.File
	file, err = os.Create(path.Join(destPath, destFile))
	if err != nil {

		return err
	}
	defer file.Close()

	_, err = io.Copy(&cipher.StreamWriter{S: stream, W: file}, bufio.NewReaderSize(src, 512))
	if err != nil {
		return err
	}
	return nil
}
