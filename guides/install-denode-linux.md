# DeNet Datakeeper Node Installation Guide for Linux

This guide provides simplified step-by-step instructions for installing and running a DeNet Datakeeper Node on Linux systems.

## Table of Contents
1. [Quick Start Installation](#quick-start-installation)
    - [Method 1: Download via curl (Recommended)](#method-1-download-via-curl-recommended)
    - [Method 2: Download from GitHub Website](#method-2-download-from-github-website)
2. [Configuring Node](#configuring-node)
3. [Run In Background](#run-in-background)
4. [Running Multiple Nodes](#running-multiple-nodes)
5. [Advanced: Systemd Service Setup (Optional)](#advanced-systemd-service-setup-optional)
6. [Reinstalling Your Node](#reinstalling-your-node)
7. [Useful Commands](#useful-commands)

## Quick Start Installation

### Method 1: Download via curl (Recommended)
```bash
# Download node executable
curl -LO https://github.com/DeNetPRO/Node/releases/download/v4.1.4/denode-linux-amd64

# Create directory for the node executable and copy it
mkdir -p ~/denet
cp denode-linux-amd64 ~/denet/denode
cd ~/denet

# Make executable
chmod +x denode

# Run the node (this will prompt for configuration)
./denode
```

### Method 2: Download from GitHub Website
1. Visit [https://github.com/DeNetPRO/Node/releases](https://github.com/DeNetPRO/Node/releases)
2. Download the appropriate binary for your system:
    - For AMD64 architecture: `denode-linux-amd64`
    - For ARM64 architecture: `denode-linux-arm64`
3. Move the downloaded file to your desired location:
   ```bash
   mkdir -p ~/denet
   mv /path/to/downloaded/denode-linux-amd64 ~/denet/denode
   cd ~/denet
   chmod +x denode
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

> 💡 **Note**: typed this way, the password stays in your shell history and in the environment of everything the node starts. On a machine with more than one account, run the node as a [systemd service](./systemd-service.md) instead — it reads the password from a root-only file.

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

## Advanced: Systemd Service Setup (Optional)

To start the node on boot and have it restart itself after a crash, run it as a systemd service. Full instructions — the unit file, where the wallet password goes, multiple nodes, upgrades — are in the **[systemd service guide](./systemd-service.md)**.

In short, with the node already installed and configured under `YOUR_USER`:

1. Keep the password out of the unit file (unit files are world-readable):
   ```bash
   sudo mkdir -p /etc/denode && sudo chmod 700 /etc/denode
   sudo install -o root -g root -m 600 /dev/null /etc/denode/denode.env
   sudo nano /etc/denode/denode.env      # one line: DENODE_PASSWORD=your_password
   ```

2. Create `/etc/systemd/system/denode.service` with the content from the [systemd service guide](./systemd-service.md#step-2-write-the-unit-file), replacing `YOUR_USER`, the wallet address and the license ID.

3. Enable and start it:
   ```bash
   sudo systemctl daemon-reload
   sudo systemctl enable --now denode.service
   sudo systemctl status denode.service     # active (running)
   ```

> ⚠️ Two things the node will not work without: `User=` must be the account you configured the node under (its config and encrypted key live in `~/.denode`), and `ExecStart` must pass `--address` and `--license` — a service has no keyboard, so a node that still has questions to ask exits immediately.

## Reinstalling Your Node

If you need to reinstall your DeNet Datakeeper Node, follow these steps:

1. **Stop the current node process**:
   ```bash
   pkill denode
   ```

2. **Remove the existing installation directory**:
   ```bash
   rm -rf ~/denet
   ```

3. **(Optional) Remove systemd service** (if previously configured):
   ```bash
   sudo systemctl stop denode.service
   sudo systemctl disable denode.service
   sudo rm /etc/systemd/system/denode.service
   sudo systemctl daemon-reload
   ```

4. **Follow the Quick Start Installation steps** from the beginning to install the node again

## Useful Commands

Started by hand (`nohup`):

- View logs: `tail -f ~/denet/denode.log`
- Check status: `ps aux | grep denode`
- Stop node: `pkill denode`

Started as a [systemd service](./systemd-service.md):

- View logs: `journalctl -u denode -f`
- Check status: `sudo systemctl status denode`
- Stop node: `sudo systemctl stop denode`

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
