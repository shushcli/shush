package lib

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/vault/shamir"
)

const (
	keySize    = 32
	nonceSize  = 12
	encryptExt = ".shush"
	shardExt   = ".shard"
)

var (
	errInvalidKey        = errors.New("invalid key file provided")
	errNotShushEncrypted = errors.New("provided file isn't shush encrypted")
	errNotEnoughShards   = errors.New("You must supply at least 2 shards to attempt to combine them into a secret")
	errFileExists        = func(path string) error { return fmt.Errorf("cannot write \"%s\"; file already exists", path) }
)

// Gen creates a new aes key and writes it to disk
func Gen(keyName string) error {
	key := make([]byte, keySize)
	_, err := rand.Read(key)
	if err != nil {
		return err
	}

	err = safeWrite(keyName, base64encode(key), 0600)
	if err != nil {
		return err
	}

	fmt.Printf("Wrote key to %s\n\n", keyName)

	return nil
}

// Split reads the fileName, and writes the shards to disk
func Split(file string, parts int, threshold int) error {
	secret, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	shards, err := shamir.Split(secret, parts, threshold)
	if err != nil {
		return err
	}

	shardNames, err := writeShards(file, shards)
	if err != nil {
		return err
	}

	fmt.Println("Successfully wrote shards:")
	for _, f := range shardNames {
		fmt.Println(" ", f)
	}
	fmt.Print("\n")

	return nil
}

// Merge reads the files, and writes the recovered secret
func Merge(files []string) error {
	if len(files) < 2 {
		return errNotEnoughShards
	}

	shards, err := readFiles(files)
	if err != nil {
		return err
	}

	fmt.Println("Merging shards:")
	for _, f := range files {
		fmt.Println(" ", filepath.Base(f))
	}
	fmt.Print("\n")

	result, err := shamir.Combine(shards)
	if err != nil {
		return err
	}

	parts := strings.Split(files[0], ".")
	dst := strings.Join(parts[0:len(parts)-1], ".")

	err = safeWrite(dst, result, 0600)
	if err != nil {
		return err
	}

	fmt.Printf("Wrote result to %s\n\n", dst)

	return nil
}

// Encrypt run aes encryption on file, using the key in keyFile
func Encrypt(keyFile string, file string) error {
	gcm, err := getGCM(keyFile)
	if err != nil {
		return err
	}

	plaintext, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	nonce := make([]byte, 12)
	if _, err := rand.Read(nonce); err != nil {
		return err
	}

	dst := fmt.Sprintf("%s%s", file, encryptExt)
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	err = safeWrite(dst, ciphertext, 0600)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully created %s\n", dst)
	return nil
}

// Decrypt decrypts a file that was encrypted with Encrypt, using the key in keyFile
func Decrypt(keyFile string, src string) error {
	gcm, err := getGCM(keyFile)
	if err != nil {
		return err
	}

	ciphertext, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	if len(ciphertext) < nonceSize {
		return errNotShushEncrypted
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return err
	}

	dst := src[:len(src)-len(encryptExt)]
	err = safeWrite(dst, plaintext, 0600)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully decrypted to %s\n", dst)

	return nil
}

// returns GCM for encrypt/decrypt
func getGCM(keyFile string) (cipher.AEAD, error) {
	b64key, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}

	key := base64decode(b64key)
	if len(key) != keySize {
		return nil, errInvalidKey
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	return cipher.NewGCM(block)
}

// file reading and writing stuff
func writeShards(originalFileName string, shards [][]byte) (shardFiles []string, err error) {
	shardFiles = make([]string, len(shards))
	for i, shard := range shards {
		shardFiles[i] = fmt.Sprintf("%s%s%d", originalFileName, shardExt, i)
		err = safeWrite(shardFiles[i], base64encode(shard), 0600)
		if err != nil {
			return nil, err
		}
	}
	return
}

// readFiles returns a slice of byte slices
func readFiles(files []string) (shards [][]byte, err error) {
	shards = make([][]byte, len(files))
	for i, f := range files {
		shard, err := ioutil.ReadFile(f)
		if err != nil {
			return nil, err
		}

		shards[i] = base64decode(shard)
	}
	return
}

// safeWrite throws errors if the file already exists
func safeWrite(path string, data []byte, perms os.FileMode) (err error) {
	_, e := os.Stat(path)
	if e != nil && os.IsNotExist(e) {
		return ioutil.WriteFile(path, data, perms)
	}

	return errFileExists(path)
}

// encoding and decoding helpers
func base64encode(in []byte) (out []byte) {
	out = make([]byte, base64.StdEncoding.EncodedLen(len(in)))
	base64.StdEncoding.Encode(out, in)
	return out
}

func base64decode(in []byte) (out []byte) {
	out = make([]byte, base64.StdEncoding.DecodedLen(len(in)))
	bytes, _ := base64.StdEncoding.Decode(out, in)
	return out[:bytes]
}
