package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	lib "github.com/shushcli/shush/lib"
)

var (
	errMissingSubCommand = errors.New("missing a valid sub-command")
	errMissingShards     = errors.New("missing list of files to merge")

	// gen errors
	errMissingKeyFile = errors.New("missing name for key file")

	// split errors
	errInvalidShardCount = errors.New("invalid number of shards")
	errInvalidThreshold  = errors.New("invalid threshhold provided")
	errMissingPath       = errors.New("missing file path of the secret")

	// encrypt/decrypt errors
	errEncryptMissingFileArg = errors.New("missing the filename to encrypt")
	errDecryptMissingFileArg = errors.New("missing the filename to decrypt")
)

var splitCmd = flag.NewFlagSet("split", flag.ExitOnError)
var encryptCmd = flag.NewFlagSet("encrypt", flag.ExitOnError)
var decryptCmd = flag.NewFlagSet("decrypt", flag.ExitOnError)

func main() {
	err := route()

	if err != nil {
		fmt.Printf("Error: %s\n\n", err)

		// print flags for relevant sub-command
		if splitCmd.Parsed() {
			fmt.Println("key split flags:")
			splitCmd.PrintDefaults()
		} else if encryptCmd.Parsed() {
			fmt.Println("encrypt flags:")
			encryptCmd.PrintDefaults()
		} else if decryptCmd.Parsed() {
			fmt.Println("decrypt flags:")
			decryptCmd.PrintDefaults()
		}

		usage()

		os.Exit(1)
	}

	os.Exit(0)
}

// handle sub-commands
func route() error {
	if len(os.Args) < 2 {
		return errMissingSubCommand
	}

	switch os.Args[1] {
	case "generate":
		return handleGen()
	case "split":
		return handleSplit()
	case "merge":
		return handleMerge()
	case "encrypt":
		return handleEncrypt()
	case "decrypt":
		return handleDecrypt()
	default:
		return errMissingSubCommand
	}
}

func usage() {
	fmt.Print(`
USAGE:

Generate a new AES Key:
	shush key gen my.key

Encrypt a secret with your key:
	shush encrypt -key=my.key secrets.tar

Split your key into 5 shards, requiring a threshold of at least 3 shards for recovery:
	shush key split -t=3 -s=5 my.key
	
Merge shards back into an AES key:
	shush key merge my.key.shard0 my.key.shard1 my.key.shard4

Merge shards with a wildcard:
	shush key merge my.key.shard*
	
Decrypt a secret with your key:
	shush decrypt -key=my.key secrets.tar.shush
`)
}

func handleGen() error {
	if len(os.Args) < 3 {
		return errMissingKeyFile
	}

	return lib.Gen(os.Args[2])
}

func handleSplit() error {
	threshold := splitCmd.Int("t", 0, "Threshold: How many shards are needed to reconstruct the messsage?")
	shardCount := splitCmd.Int("s", 0, "Shards: How many total shards will we generate")
	splitCmd.Parse(os.Args[2:])

	if *shardCount < 2 {
		return errInvalidShardCount
	} else if *threshold > *shardCount || *threshold < 2 {
		return errInvalidThreshold
	}

	args := splitCmd.Args()
	if len(args) < 1 {
		return errMissingPath
	}

	err := lib.Split(args[0], *shardCount, *threshold)
	if err != nil {
		return err
	}

	return nil
}

func handleMerge() error {
	if len(os.Args) < 3 {
		return errMissingShards
	}

	err := lib.Merge(os.Args[2:])
	if err != nil {
		return err
	}
	return nil
}

func handleEncrypt() error {
	keyFile := encryptCmd.String("key", "", "Key: Path to your key file")
	encryptCmd.Parse(os.Args[2:])

	if *keyFile == "" {
		return errMissingKeyFile
	}

	if len(os.Args) < 4 {
		return errEncryptMissingFileArg
	}

	return lib.Encrypt(*keyFile, os.Args[3])
}

func handleDecrypt() error {
	keyFile := decryptCmd.String("key", "", "Key: Path to your key file")
	decryptCmd.Parse(os.Args[2:])

	if *keyFile == "" {
		return errMissingKeyFile
	}

	if len(os.Args) < 4 {
		return errDecryptMissingFileArg
	}

	return lib.Decrypt(*keyFile, os.Args[3])
}
