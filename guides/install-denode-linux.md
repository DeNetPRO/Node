# DeNet Datakeeper Node Installation Guide for Linux

This guide provides simplified step-by-step instructions for installing and running a DeNet Datakeeper Node on Linux systems.

## Table of Contents
1. [Quick Start Installation](#quick-start-installation)
    - [Method 1: Download via curl (Recommended)](#method-1-download-via-curl-recommended)
    - [Method 2: Download from GitHub Website](#method-2-download-from-github-website)
2. [Configuring Node](#configuring-node)
3. [Running Your Node](#running-your-node)
4. [Advanced: Systemd Service Setup (Optional)](#advanced-systemd-service-setup-optional)
5. [Useful Commands](#useful-commands)

## Quick Start Installation

### Method 1: Download via curl (Recommended)
```bash
# Download node executable
curl -LO https://github.com/DeNetPRO/Node/releases/download/v4.0.1-rc10/denode-linux-amd64

# Create directory for the node executable and copy it
mkdir -p ~/denet
cp denode-linux-amd64 ~/denet/denode
cd ~/denet

# Make executable
chmod +x denode

# Run the node (this will prompt for configuration)
./denode
```
![](assets/linux-run.png)

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
   ```
   
## Configuring Node
- Instructions [here](./configuring.md)

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

## Advanced: Systemd Service Setup (Optional)

For automatic startup on boot:

1. Create a dedicated user (recommended):
   ```bash
   sudo adduser denet
   ```

2. Create service file:
   ```bash
   sudo nano /etc/systemd/system/denode.service
   ```

3. Add this content (replace `username` with your actual username):
   ```ini
   [Unit]
   Description=DeNode Service
   After=network.target

   [Service]
   User=username
   Group=username
   Type=simple
   ExecStart=/home/username/denet/denode
   Restart=always
   RestartSec=5

   [Install]
   WantedBy=multi-user.target
   ```

4. Enable and start:
   ```bash
   sudo systemctl daemon-reload
   sudo systemctl enable denode.service
   sudo systemctl start denode.service
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
