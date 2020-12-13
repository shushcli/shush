package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/shushcli/shush/shards"
)

func main() {
	encryptCommand := flag.NewFlagSet("encrypt", flag.ExitOnError)
	threshold := encryptCommand.Int("t", 0, "Threshold: How many shards are needed to reconstruct the messsage?")
	totalShards := encryptCommand.Int("s", 0, "Shards: How many total shards will we generate")

	decryptCommand := flag.NewFlagSet("decrypt", flag.ExitOnError)

	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Println("encrypt or decrypt sub-command is required")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "encrypt":
		handleEncrypt(encryptCommand, totalShards, threshold)
	case "decrypt":
		handleDecrypt(decryptCommand)
	default:
		fmt.Println("encrypt or decrypt sub-command is required")
		os.Exit(1)
	}

	os.Exit(0)
}

func handleEncrypt(encryptCommand *flag.FlagSet, totalShards *int, threshold *int) {
	encryptCommand.Parse(os.Args[2:])
	if *totalShards < 2 {
		fmt.Println("invalid number of totalShards")
		encryptCommand.PrintDefaults()
		os.Exit(1)
	} else if *threshold > *totalShards || *threshold < 2 {
		fmt.Println("invalid threshhold")
		encryptCommand.PrintDefaults()
		os.Exit(1)
	}

	args := encryptCommand.Args()
	if len(args) < 1 {
		fmt.Print("You must include a file path to the secret as a trailing argument.\n\n")
		fmt.Print("Example: shush encrypt -t=3 -s=6 \"secret.txt\"\n\n")
		encryptCommand.PrintDefaults()
		os.Exit(1)
	}

	err := shards.Create(args[0], *totalShards, *threshold)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func handleDecrypt(decryptCommand *flag.FlagSet) {
	decryptCommand.Parse(os.Args[2:])
	args := decryptCommand.Args()
	if len(args) < 1 {
		fmt.Println("You must include a glob pattern to match the shards")
		fmt.Print("Example: shush decrypt \"secret.txt.shard*\"")
		os.Exit(1)
	}

	err := shards.Recover(args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
