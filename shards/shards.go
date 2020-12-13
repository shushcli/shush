package shards

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/hashicorp/vault/shamir"
)

// Create reads the secret file, and writes the shards to disk
func Create(secretFileName string, parts int, threshold int) error {
	secret, err := ioutil.ReadFile(secretFileName)
	if err != nil {
		return err
	}

	shards, err := shamir.Split(secret, parts, threshold)
	if err != nil {
		return err
	}

	shardNames, err := writeFiles(secretFileName, shards)
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

// Recover reads the file file, and writes the recovered secret to out file
func Recover(file string) error {
	files, err := buildFileList(file)
	if err != nil {
		return err
	}

	shards, err := readFiles(files)
	if err != nil {
		return err
	}

	fmt.Println("Decrypting shards:")
	for _, f := range files {
		fmt.Println(" ", filepath.Base(f))
	}
	fmt.Print("\n\n")

	result, err := shamir.Combine(shards)
	if err != nil {
		return err
	}

	parts := strings.Split(file, ".")
	base := parts[0]
	after := strings.Join(parts[1:len(parts)-1], ".")
	out := fmt.Sprintf("%s-recovered.%s", base, after)

	err = ioutil.WriteFile(out, result, os.ModePerm)
	if err != nil {
		return err
	}

	fmt.Printf("Wrote result to %s\n\n", out)

	return nil
}

func buildFileList(glob string) (files []string, err error) {
	cwd, err := os.Getwd()
	if err != nil {
		return
	}

	files, err = filepath.Glob(path.Join(cwd, glob))
	if err != nil {
		return
	}

	if len(files) < 2 {
		err = fmt.Errorf("You must supply at least 2 shards to attempt to combine them into a secret")
	}

	return
}

// file reading and writing  stuff
func writeFiles(secret string, shards [][]byte) (shardFiles []string, err error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	for i, shard := range shards {
		fileName := fmt.Sprintf("%s.shard%d", secret, i)
		shardFiles = append(shardFiles, fileName)
		err = ioutil.WriteFile(path.Join(cwd, shardFiles[i]), base64encode(shard), os.ModePerm)
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
