# DeNode CLI Account and Config Management Guide

## Overview
The `account` and `config` commands allow users to manage Ethereum wallets and node configuration files for the Denode (peaq DePIN node).

Global flags available for all commands:
- `-d, --debug` Enable debug output
- `-h, --help`  Show help

## Account Management Commands

### Command: account delete

**Purpose:** Delete your ethereum account

**Usage:**
```
./denode account delete
```

### Command: account export

**Purpose:** Export your ethereum account private key

**Usage:**
```
./denode account export
```

### Command: account import

**Purpose:** Imports your wallet by private key

**Usage:**
```
./denode account import
```

## Config Management Commands

### Command: config clear

**Purpose:** Deletes all user's configuration files when it's needed to start from scratch

**Usage:**
```
./denode config clear
```

### Command: config code

**Purpose:** Generates a code for obtaining a role on the Discord server

**Usage:**
```
./denode config code
```

### Command: config generate

**Purpose:** Generates a new config file in the working directory (all command line parameters are required for non-interactive mode)

**Flags:**
* `--address` string  (required) Account address to create config for
* `--license` string  (required for non-interactive) License ID
* `--storage` string  (required for non-interactive) Storage directory path
* `--share` int    Storage share size in GiB (minimum 1)
* `--rpc` string    peaq blockchain RPC URL
* `--ip` string    Node server IP address (default "0.0.0.0")
* `--port` int    Port number (10000-65535) (default 55050)

**Usage example:**
```
./denode config generate \
  --address 0x123...abc \
  --license 12345 \
  --storage /home/user/denode_storage \
  --share 200 \
  --rpc https://rpc-peaq.peaq.network \
  --ip 0.0.0.0 \
  --port 55050
```

### Command: config get

**Purpose:** Display the decrypted configuration

**Flags:**
* `--address` string  Account address to get config for
* `--license` string  (required for non-interactive mode) License ID

**Usage (interactive mode):**
```
./denode config get --address 0x123...abc
```

**Usage (non-interactive mode):**
```
./denode config get --address 0x123...abc --license 12345
```

### Command: config set

**Purpose:** Modifies specified fields in a configuration file (only in interactive mode)

**Flags:**
* `--address` string  (required) Account address to modify config for

**Usage:**
```
./denode config set --address 0x123...abc
```

The command will launch an interactive menu where you can change individual fields (license, storage path, share size, RPC, port, etc.).
