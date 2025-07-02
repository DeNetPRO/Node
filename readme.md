<p align="center">
    <img src="assets/LOGO.png">
</p>

<p align="center">
    Monetize your computer's storage now!
    <br/>
    <br/>
    <a href="https://denet.pro">
        <img alt="website.png" src="assets/denet.pro.svg" height="31" width="120"/>
    </a>
    <a href="https://t.me/+Yu5KnSruttc5ZGRi">
        <img alt="tg.png" src="https://img.shields.io/badge/Telegram-2CA5E0?style=for-the-badge&logo=telegram&logoColor=white" height="30" width="120"/>
    </a>
    <a href="https://discord.gg/cPz9m4cSWv">
        <img alt="discord.png" src="https://img.shields.io/badge/Discord-5865F2?style=for-the-badge&logo=discord&logoColor=white" height="30" width="120" />
    </a>
    <a href="https://www.youtube.com/channel/UCeCxt3tYbtSkJvaznNjQimQ">
        <img alt="youtube.png" src="https://img.shields.io/badge/YouTube-FF0000?style=for-the-badge&logo=youtube&logoColor=white" height="30" width="120" />
    </a>
</p>

💽 **Datakeepers** offer their unused storage space to DeNet Storage users, creating a mutually beneficial conditions for all without any intermediaries. \
By utilizing the DeNet Storage Protocol, users pay for the storage they need, while Datakeepers are rewarded for the storage they provide.

# Becoming Datakeeper: node set up instruction

This guide walks you through setting up and running a DeNet Node, enabling your device to join the DeNet decentralized storage network

### Table of contents:

- [Denode](#denode)
  - [Installation](#installation)
    - [Requirements](#requirements)
    - [Step 0: License Verification](#step-0-verify-your-account-has-license)
    - [Step 1: Copy Private Key](#step-1-copy-your-private-key)
    - [Step 2: Download Datakeeper Node](#step-2-download-datakeeper-node)
    - [Step 3: Start Node](#step-3-start-denet-node)
      - [Windows](#windows)
      - [Linux](#linux)
      - [macOS](#macos)
    - [Step 4: Run DeNet Node](#step-4-run-denet-node)
    - [Step 5: Monitor Transactions](#step-5-monitor-transactions)
- [Node Manager (GUI)](#denode-manager-gui)
  - [Installation](#installation-1) 
    - [Step 0: Prepare Environment](#step-0-prepare-environment)
    - [Step 1: Download Node Manager](#step-1-download-application)
    - [Step 2: Install And Run](#step-2-install-and-run)
      - [Windows](#windows-1)
      - [Linux](#linux-1)
      - [macOS](#macos-1)

# Denode
## Installation
## Requirements

  - A wallet address (DeNet app/Metamask/any other wallet) holding a Datakeeper Node License.
  - Device with free disk space for storing DeNet user data.
  - Terminal access (Command Prompt/PowerShell on Windows, Terminal on macOS/Linux).
  - DeNet Node application downloaded.
  - PEAQ balance: tokens will be distributed to Datakeepers automatically and will be regularly credited for successfully completed transactions, if the node is running and does not disconnect from the network, no deposits will be required.

## Step 0: Verify your account has license

  - Open https://peaq.subscan.io/account/YOUR_ADDRESS
  - Replace YOUR_ADDRESS with your wallet address.
  - The license(s) should be seen as a sNL ERC-721 token.

![](assets/license.png)

## Step 1: Copy Your Private Key
You need the private key from a wallet with a Datakeeper Node License.
### From DeNet App
1. Open the DeNet app.
2. Go to "Profile" -> "Settings"-> “Security”.
3. Copy the 64-character HEX private key (e.g., a1b2c3d4...).
4. Save it securely. Never share your private key!
### From other wallet (we take Metamask as an example)
1. Open Metamask in your browser or app.
2. Select the account with the Datakeeper Node License.
3. Go to "Account Details" > "Export Private Key."
4. Enter your Metamask password and copy the private key.
5. **Store it securely! Do not share it anywhere!**

_If you use any other wallet, the steps may differ but should be similar to the list above._

## Step 2: Download Datakeeper Node

1. Visit [https://github.com/DeNetPRO/Node/releases](https://github.com/DeNetPRO/Node/releases)
2. Download the latest application executable for your OS (Linux, macOS, or Windows)
   ![](assets/executables.png)
   **macOS**: use amd64 for Intel hardware, arm64 for Apple Silicon.

3. Copy application to another directory.
   **Example:**
   - macOS/Linux: copy denode executable to `~/denet/` directory
   - Windows: copy to `C:\denet\` directory

## Step 3: Start DeNet Node
Launch the node via a terminal.

### Windows
1. Open Terminal: Press **Win + R**, type cmd or powershell, and press Enter.
   ![](assets/win-cmd.png)
2. Start the Application:
- Navigate to the Application Folder
   ![](assets/win-folder.png)
- Run application executable
   ![](assets/win-run.png)
### Linux
1. Open Terminal: Use Ctrl + Alt + T or your terminal shortcut.
   Or SSH to your remote host.
2. Run the following commands to create folder, copy and run denode
```bash
curl -LO https://github.com/DeNetPRO/Node/releases/download/v4.0.0-rc1/denode-linux-amd64
mkdir ~/denet
cp denode-linux-amd64 ~/denet/denode
cd ~/denet
chmod +x denode
```
   ![](assets/linux-run.png)

### macOS
1. Open "Terminal" via Spotlight or Applications
   ![](assets/mac-terminal.png)
2. Run the following commands to create folder, copy and run denode
```bash
curl -LO https://github.com/DeNetPRO/Node/releases/download/v4.0.0-rc1/denode-macos-amd64
mkdir ~/denet
cp ~/Downloads/denode-macos-amd64 ~/denet/denode
cd ~/denet
chmod +x denode
xattr -d com.apple.quarantine denode
```
   ![](assets/mac-run.png)

## Step 4: Run DeNet Node
1. **Enter private key**: Paste the copied private key and press Enter.
- The key is stored securely on your device, encrypted with this password.

2. **Set Password**: Enter a strong password
- The private key is encrypted with the password.

3. **Choose Port**: Press Enter for the default one
- Or specify another (value from 10000 to 65535)
4. **Specify Storage Directory**: Enter path to the user files storage
- **e.g.**, /home/user/denet_storage (Linux/macOS) or C:\denet_storage (Windows).
- Ensure the directory exists and has sufficient space.
5. **Set Storage Space**:
- Specify the amount of disk space to allocate for DeNet Storage (e.g., 10). Enter the value (only number, without GiB) when prompted.
7. **Optional Second Drive**: Enter 'N' to skip.
- Or if you want to use another drive, provide its path when prompted.
8. **Select RPC for peaq Blockchain**: Press Enter to use default one.
- Or choose the RPC endpoint (Select RPC for peaq (ChainID: 3338)).
9. **Verify Operation**:
- Watch the terminal output. If no errors appear, your DeNet Node is running correctly.
   ![](assets/successful-launch.png)

## Step 5: Monitor Transactions

Track your node’s activity using the peaq Subscan web interface.
1. Visit the peaq Subscan website (e.g., https://peaq.subscan.io/account/YOUR_ADDRESS).
  - Search for your node’s transactions using your Datakeeper address.
2. Check transaction statuses. Green check marks indicate successful transactions, confirming your node is working correctly.
   ![](assets/successful-trxs.png)
## Troubleshooting

- **Errors in Terminal**: Carefully check error message. Most of the errors are related to the lack of Internet, insufficient balance of gas tokens, or the result of manually changing the data generated by the node. Ask for help from community members or contact support in Discord.
- **Failed Transactions in Subscan**: There may be some unusual situations where transactions fail. If you encounter such cases, please open a support ticket.
  ![](assets/failed-trx.png)
- **Port Conflicts**: If port 55050 is in use, try another port (e.g., 55051)
  ![](assets/port-error.png)
- **Subscan Issues:** If transactions don’t appear, confirm your node is running and has enough gas tokens (> 0.03 $PEAQ).
- **Not Opened macOS:**
   - Run the following command to allow `denode` executable:
   - `xattr -d com.apple.quarantine denode`
   ![](assets/mac-error-01.png)
- **Permission denied**
   ```bash
   user@desktop:~/denet$ ./denode
   -bash: ./denode: Permission denied
   ```
   - Allow execution by running `chmod +x denode`
## Notes
- Keep your terminal open to maintain the node’s operation. Closing it stops the node. Otherwise, set up the node as a background service (see [Advanced Settings: Systemd Service](#advanced-settings)).
- Additional steps (e.g., advanced settings) will be added as needed — check for updates from DeNet.

Congratulations, Datakeeper! Your DeNet Node is now contributing to the decentralized storage network.

A graphical user interface (GUI) for seamless node operation coming soon. Stay tuned!

#### Systemd service (Linux)
- It is recommended to run as non-root user.
- Follow [Step 4](#step-4-run-denet-node) first. 
  You should have config files in the path `/home/denet/.denode/`
- We used **denet** user as an example, replace it with your username.

**/etc/systemd/system/denode.service**
```ini
[Unit]
Description=DeNode Service
After=network.target

[Service]
User=denet
Group=denet
Type=simple
ExecStart=/usr/local/bin/denode
EnvironmentFile=/home/denet/denode.env
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```
**/home/denet/denode.env**
```text
DENODE_PASSWORD='your password'
```

Run the following commands to enable startup and run node
    ```shell
    sudo systemctl daemon-reload
    sudo systemctl enable denode.service
    sudo systemctl start denode.service
    ```
Now your node will be running and start at boot.

**View latest logs**
`journalctl -u denode -r`

# Denode Manager GUI
## Installation

## Step 0: Prepare environment
1. Remove old credentials and configurations, they should be imported from scratch:
    ### Linux/macOS
    ```shell
    rm -rf ~/.denode
    ```
    ### Windows
    Use cmd or graphical user interface to remove the folder.
   ![](assets/win-rm-denode.png)

2. For **Linux/macOS**: Download installation and management scripts from the [scripts](https://github.com/DeNetPRO/Node) directory
    ```shell
    install.sh
    denode-manager.sh
    ```
## Step 1: Download Application
1. Download archive for your system from https://github.com/DeNetPRO/Node/releases as well as for [denode](#step-2-download-datakeeper-node)
    ### Windows
    ```
    denode-manager-win-amd64.msi
    ```
    ### Linux
    ```
    denode-manager-linux-amd64.zip
    denode-manager-linux-arm64.zip
    ```
    ### macOS
    ```
    denode-manager-darwin-amd64.zip
    denode-manager-darwin-arm64.zip
    ```
## Step 2: Install And Run
### macOS/Linux
1. Open terminal as for [denode installation](#macos)
2. Allow scripts execution on this device 
    ```shell
    cd ~/Downloads
    chmod +x install.sh denode-manager.sh
    ```
3. Run installation script that will install the application in ~/.denode-manager by default
    ```shell
   sudo bash install.sh
   ```
   ![](assets/mac-install-gui.png)
4. Then start the application and check its state using ***denode-manager.sh*** script
    ```shell
   sudo bash denode-manager.sh
   ```
   ![](assets/mac-start-server.png)
### Windows
Use .msi installer
## Step 3: Open Application Interface in Browser
1. Open browser and go to http://localhost:1111
   ![](assets/node-gui.png)

## NOTES:
1. Application should always be running in the background, otherwise the application will not work, check the status using ***denode-manager.sh*** script (Mac/Linux)
2. You shouldn't use both CLI and GUI at the same time, otherwise you will get an undefined applications behaviour.
3. We recommend to change rpc to the private one exactly after the node launch. It will allow to avoid problems with the default version limitations.

#### Ask your questions here and get help:

<a href="https://discord.gg/cPz9m4cSWv">
    <img alt="discord.png" src="https://img.shields.io/badge/Discord-5865F2?style=for-the-badge&logo=discord&logoColor=white" height="30" width="120" />
</a>
