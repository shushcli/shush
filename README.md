# Shush 🤫
This simple program will help you run [Shamir's Secret Sharing](https://en.wikipedia.org/wiki/Shamir%27s_Secret_Sharing) algorithm on _any_ file using the `split` and `merge` commands. It also contains  tools to easily `generate` an [AES](https://en.wikipedia.org/wiki/Advanced_Encryption_Standard) key and `encrypt` and `decrypt` files using said key.

## Why is that useful?
I've you've distributed the shards of an AES key to your family or team, they will be able to recover any encrypted data you left behind if you lose it, become incapacitated, or worse.

## Why not just split a key into chunks?
With Sharmir's algorithm, you can specify a threshold for recovery that is lower than the total number of shards. This approach protects you against some members of your family or team losing their shards.

## Safely generating and distributing keys & encrypted payloads
Run this program in a secure operating system like [Tails](https://en.wikipedia.org/wiki/Tails_%28operating_system%29) with no internet connection. Be extremely careful about how you store your key! Distribute shards to your team on physical media like flash drives. You may also want to relay information about who holds other key shards to the folks on your team, though ideally this isn't in writing or included on the flash drives.

## What to include when distributing keys shards
You may want to consider including any of the following things when distributing key shards:
- instructions on how to merge keys and decrypt files
- information about the location of other potential payloads
- a copy of the encrypted payload
- a copy of shush
- a copy of the shush source code

## Safely merging key shards and decrypting payloads
Since the payload likely has sensitive contents, you should take similar precautions (tails, offline, etc.) when re-assembling key shards and decrypting payloads.

## Encrypting additional secrets at a later date
If you hold onto your original AES key, you can create new encrypted payloads whenever you want, and redistribute them or put them online without having to generate new keys or distribute new key shards.

## Usage
### 1) Generate a new AES Key:
```bash
# shush generate <filename>
shush generate my.key
```

### 2) Encrypt a secret file or archive with your AES Key:
```bash
shush encrypt -key=my.key secrets.tar
```

Which will write an encrypted payload file to disk:
```
secrets.tar.shush
```

### 3) Split the key into 5 shards, requiring a threshold of at least 3 shards for recovery:
```bash
shush split -t=3 -s=5 my.key
```

Which will write the following shards to disk:
```
my.key.shard0
my.key.shard1
my.key.shard2
my.key.shard3
my.key.shard4
```

### 4) Distribute the keys (and secrets) to a team of trusted parties (do this offline), and wait for something bad to happen.

### 5) Someone from your team will recollect some of the original shards, and merge them back into an AES key:
```bash
# Note that the shard format allows you to use a wildcard: my.key.shard*
shush merge my.key.shard0 my.key.shard2 my.key.shard4
```

### 6) Recover the payload using the recombined key
```bash
shush decrypt -key=my.key secrets.tar.shush
```

## Building & Installing
```bash
# On a unix-based system with go installed...
go build -o shush main.go
# install on your system
mv shush /usr/local/bin
```
