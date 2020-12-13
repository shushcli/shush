# Shush 🤫
This simple program will help you run [Shamir's Secret Sharing](https://en.wikipedia.org/wiki/Shamir%27s_Secret_Sharing) algorithm on _any_ file using the `split` and `merge` commands. It also contains  tools to easily `generate` an [AES](https://en.wikipedia.org/wiki/Advanced_Encryption_Standard) key and `encrypt` and `decrypt` files using said key.

## Why?
If you've distributed the shards of an AES key to your team (read: family, friends, coworkers), they will be able to recover any encrypted data in case you lose it, become incapacitated, or worse.

## Can't I just split a key into chunks, and distribute the chunks?
With Shamir's algorithm, you can specify a `threshold` for recovery that is lower than the total number of `shards`. This approach protects you against some members of your team losing their shards.

## Safely generating and distributing keys & encrypted payloads
Run this program in [Tails](https://en.wikipedia.org/wiki/Tails_%28operating_system%29) with no internet connection. Be extremely careful about how you store your key! Distribute shards to your team on physical media (like flash drives). You may also want to notify your team members who else is on their team, but ideally that information will live in their heads, not in their emails.

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

### Encrypt and Decrypt Files
```bash
# Generate a new AES Key
shush generate my.key

# Encrypt a secret file or archive with your AES Key
shush encrypt -key=my.key secrets.tar

# Decrypt a payload using an AES key
shush decrypt -key=my.key secrets.tar.shush
```

### Split and Merge Files
```bash
# Split a file into 5 shards, requiring a threshold of at least 3 shards for recovery
shush split -t=3 -s=5 my.key

# Merge shards back into the original file
shush merge my.key.shard0 my.key.shard2 my.key.shard4

# You can also use a wildcard if the names are preserved.
shush merge my.key.shard*
```

## Building & Installing
```bash
# On a unix-based system with go installed...
go build -o shush main.go
# install on your system
mv shush /usr/local/bin
```
