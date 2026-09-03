# DeNet Datakeeper Node Installation Guide for macOS

This guide provides simplified step-by-step instructions for installing and running a DeNet Datakeeper Node on macOS systems.

1. [Quick Start Installation](#quick-start-installation)
2. [Configuring Node](#configuring-node)
3. [Run in background](#run-in-background)
4. [Running Multiple Nodes](#running-multiple-nodes)
5. [Useful Commands](#useful-commands)

## Quick Start Installation

### Method 1: Download via curl (Recommended)
```bash
# Download node executable (replace with darwin-arm64 or darwin-amd64 for your architecture)
curl -LO https://github.com/DeNetPRO/Node/releases/download/v4.1.3/denode-darwin-arm64

# Create directory for the node executable and copy it
mkdir -p ~/denet
cp denode-darwin-arm64 ~/denet/denode
cd ~/denet

# Make executable and remove quarantine attribute (required on macOS)
chmod +x denode
xattr -d com.apple.quarantine ./denode 2>/dev/null

# Run the node (this will prompt for configuration)
./denode
```

### Method 2: Download from GitHub Website
1. Visit [https://github.com/DeNetPRO/Node/releases](https://github.com/DeNetPRO/Node/releases)
2. Download the appropriate binary for your system:
    - For Intel hardware: `denode-darwin-amd64`
    - For Apple Silicon (ARM64): `denode-darwin-arm64`
3. Move the downloaded file to your desired location:
   ```bash
   mkdir -p ~/denet
   mv /path/to/downloaded/denode-darwin-arm64 ~/denet/denode
   cd ~/denet
   chmod +x denode
   xattr -d com.apple.quarantine ./denode 2>/dev/null
   ./denode
   ```
   
## Configuring Node

After the node is installed, you can proceed with the configuration

- Instructions are [here](configuring-cli.md)

## Run In Background

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

## Running Multiple Nodes

When running multiple DeNet Datakeeper Nodes on the same machine, each node needs its own configuration directory and unique parameters to avoid conflicts:

1. **Create separate directories for each node:**
   ```bash
   # Create directories for different nodes
   mkdir -p ~/denet-node1 ~/denet-node2
   cp ~/denet/denode ~/denet-node1/
   cp ~/denet/denode ~/denet-node2/
   ```

2. **Configure each node separately:**
   ```bash
   # Configure first node
   cd ~/denet-node1
   ./denode  # This will prompt for configuration
   
   # Configure second node
   cd ~/denet-node2
   ./denode  # This will prompt for configuration
   ```

3. **Run each node with unique parameters:**
   ```bash
   # Run first node in background
   cd ~/denet-node1
   DENODE_PASSWORD=your_password nohup ./denode --address your_address_1 --license license_1 > denode.log 2>&1 &

   # Run second node in background
   cd ~/denet-node2
   DENODE_PASSWORD=your_password nohup ./denode --address your_address_2 --license license_2 > denode.log 2>&1 &
   ```

4. **Monitor processes:**
   ```bash
   # List all denode processes
   ps aux | grep denode
   
   # View logs for each node
   tail -f ~/denet-node1/denode.log
   tail -f ~/denet-node2/denode.log
   
   # Kill specific node by PID
   kill -9 <PID>
   
   # Kill all denode processes
   pkill denode
   ```

**Important:** Each node requires its own unique license id. When configuring each node, make sure to use a different license number for each instance.

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
