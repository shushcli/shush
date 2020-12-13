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
	"path"
	"path/filepath"
	"strings"

	"github.com/hashicorp/vault/shamir"
)

const nonceSize = 12
const suffix = ".shush"

var errInvalidKey = errors.New("invalid key file provided")

// Gen creates a new aes key and writes it to disk
func Gen(keyName string) error {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return err
	}

	err = safeWrite(keyName, base64encode(key), 0600)
	if err != nil {
		return err
	}

	return nil
}

// Split reads the fileName, and writes the shards to disk
func Split(fileName string, parts int, threshold int) error {
	secret, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	shards, err := shamir.Split(secret, parts, threshold)
	if err != nil {
		return err
	}

	shardNames, err := writeFiles(fileName, shards)
	if err != nil {
		return err
	}

	fmt.Println("Successfully wrote shards:")
	for _, f := range shardNames {
		fmt.Println(" ", f)
	}
	fmt.Print("\n\n")

	return nil
}

// Merge reads the files, and writes the recovered secret
func Merge(files []string) error {
	if len(files) < 2 {
		return fmt.Errorf("You must supply at least 2 shards to attempt to combine them into a secret")
	}

	shards, err := readFiles(files)
	if err != nil {
		return err
	}

	fmt.Println("Merging shards:")
	for _, f := range files {
		fmt.Println(" ", filepath.Base(f))
	}
	fmt.Print("\n\n")

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
	b64key, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return err
	}

	key := base64decode(b64key)
	if len(key) != 32 {
		return errInvalidKey
	}

	contents, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	nonce := make([]byte, 12)
	if _, err := rand.Read(nonce); err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	ciphertext := gcm.Seal(nonce, nonce, contents, nil)

	dst := fmt.Sprintf("%s%s", file, suffix)

	err = safeWrite(dst, ciphertext, 0600)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully created %s\n", dst)
	return nil
}

// Decrypt run aes encryption on file, using the key in keyFile
func Decrypt(keyFile string, file string) error {
	b64key, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return err
	}

	key := base64decode(b64key)
	if len(key) != 32 {
		return errInvalidKey
	}

	ciphertext, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	if len(ciphertext) < nonceSize {
		return errors.New("provided file isn't shush encrypted")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return err
	}

	fileWithoutSuffix := file[:len(file)-len(suffix)]
	err = safeWrite(fileWithoutSuffix, plaintext, 0600)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully decrypted to %s\n", fileWithoutSuffix)

	return nil
}

// file reading and writing stuff
func writeFiles(secret string, shards [][]byte) (shardFiles []string, err error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	for i, shard := range shards {
		fileName := fmt.Sprintf("%s.shard%d", secret, i)
		shardFiles = append(shardFiles, fileName)
		err = safeWrite(path.Join(cwd, shardFiles[i]), base64encode(shard), 0600)
		if err != nil {
			return nil, err
		}
	}
	return
}

func readFiles(files []string) (shards [][]byte, err error) {
	for _, f := range files {
		shard, err := ioutil.ReadFile(f)
		if err != nil {
			return nil, err
		}

		shards = append(shards, base64decode(shard))
	}
	return
}

// safeWrite throws errors if the file already exists
func safeWrite(path string, data []byte, perms os.FileMode) (err error) {
	_, e := os.Stat(path)
	if e != nil && os.IsNotExist(e) {
		return ioutil.WriteFile(path, data, perms)
	}

	return fmt.Errorf("attempted to write to a file that already exists: %s", path)
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
