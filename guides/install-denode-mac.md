# DeNet Datakeeper Node Installation Guide for macOS

This guide provides simplified step-by-step instructions for installing and running a DeNet Datakeeper Node on macOS systems.

1. [Quick Start Installation](#quick-start-installation)
2. [Configuring Node](#configuring-node)
3. [Running Your Node](#running-your-node)
4. [Useful Commands](#useful-commands)

## Quick Start Installation

### Method 1: Download via curl (Recommended)
```bash
# Download node executable
curl -LO https://github.com/DeNetPRO/Node/releases/download/v4.0.1-rc10/denode-macos-amd64

# Create directory for the node executable and copy it
mkdir -p ~/denet
cp denode-macos-amd64 ~/denet/denode
cd ~/denet

# Make executable
chmod +x denode

# Remove quarantine attribute (required for macOS)
xattr -d com.apple.quarantine denode

# Run the node (this will prompt for configuration)
./denode
```
![](assets/mac-run.png)

### Method 2: Download from GitHub Website
1. Visit [https://github.com/DeNetPRO/Node/releases](https://github.com/DeNetPRO/Node/releases)
2. Download the appropriate binary for your system:
    - For Intel hardware: `denode-macos-amd64`
    - For Apple Silicon (ARM64): `denode-macos-arm64`
3. Move the downloaded file to your desired location:
   ```bash
   mkdir -p ~/denet
   mv /path/to/downloaded/denode-macos-amd64 ~/denet/denode
   cd ~/denet
   chmod +x denode
   xattr -d com.apple.quarantine denode
   ```

## Configuring Node
- Instructions are [here](./configuring.md)

## Running Your Node

After initial setup of all configuration parameters, simply run:
```bash
cd ~/denet
./denode
```

To run in background:
```bash
DENODE_PASSWORD=your_password nohup ./denode --address you_datakeeper_address --license your_license_number > denode.log 2>&1 &
```

**Output:**
```
$ nohup ./denode > denode.log 2>&1 &
[1] 12345
$ ps aux | grep denode
user     12345  0.1  0.2  123456  7890 pts/0    S    10:30   0:00 ./denode
```

## Useful Commands

- View logs: `tail -f ~/denet/denode.log`
- Check status: `ps aux | grep denode`
- Stop node: `pkill denode`

**Output Example:**
```
$ tail -f ~/denet/denode.log
[2023-01-01 10:30:00] INFO: Node started successfully
[2023-01-01 10:30:05] INFO: Connected to network
[2023-01-01 10:30:10] INFO: Serving data request

$ ps aux | grep denode
user     12345  0.1  0.2  123456  7890 pts/0    S    10:30   0:00 ./denode
user     12347  0.0  0.1  123456  3456 pts/0    S    10:30   0:00 grep denode

$ pkill denode
$ ps aux | grep denode
user     12347  0.0  0.1  123456  3456 pts/0    S    10:30   0:00 grep denode
```
