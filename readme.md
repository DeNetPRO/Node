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

## Table of Contents

- [Overview](#overview)
- [Requirements](#requirements)
    - [Step 0: Verify Your Account Has License](#step-0-verify-your-account-has-license)
    - [Step 1: Copy Your Private Key](#step-1-copy-your-private-key)
- [Supported Platforms](#supported-platforms)
- [Installation Options](#installation-options)
- [Datakeeper Node](#datakeeper-node-installation)
- [Node Manager](#node-manager-installation)
- [License Management](#license-management)

TODO: add transactions monitoring info

# Becoming Datakeeper: Node Setup Instructions

This guide walks you through setting up and running a Datakeeper Node, enabling your device to join the DeNet decentralized storage network.

## Overview

There are two primary ways to manage your DeNet Node:

1. [**CLI Node**](#datakeeper-node-installation): A command-line interface application that runs in your terminal
2. [**Desktop Node Manager**](#node-manager-installation): A graphical desktop application for easier management

Download the Desktop Node Manager: [Node Manager Desktop Application](https://node.denet.app/)

# Supported Platforms

DeNet Datakeeper Nodes can be installed on various operating systems depending on your needs and technical requirements:

## Available Builds

### DeNet Node Binary Builds

| Operating System | Architecture | Download Link |
|------------------|--------------|---------------|
| **Windows** | x86_64 | [denode-windows-amd64.exe](https://github.com/DeNetPRO/Node/releases/download/v4.0.1-rc10/denode-windows-amd64.exe) |
| **Linux** | x86_64 | [denode-linux-amd64](https://github.com/DeNetPRO/Node/releases/download/v4.0.1-rc10/denode-linux-amd64) |
| **Linux** | ARM64 | [denode-linux-arm64](https://github.com/DeNetPRO/Node/releases/download/v4.0.1-rc10/denode-linux-arm64) |
| **macOS** | x86_64 | [denode-macos-amd64](https://github.com/DeNetPRO/Node/releases/download/v4.0.1-rc10/denode-macos-amd64) |
| **macOS** | ARM64 | [denode-macos-arm64](https://github.com/DeNetPRO/Node/releases/download/v4.0.1-rc10/denode-macos-arm64) |

### How to Choose the Right Architecture

When downloading the DeNet Node Binary Build, it's important to select the version that matches your computer's processor architecture:

- **x86_64 (also called AMD64)**: This is the most common architecture for Intel and AMD processors used in most desktops and laptops
- **ARM64 (also called AArch64)**: This is used by Apple Silicon (M1, M2, etc.) Macs and some ARM-based Linux devices

To determine your system architecture:

**On macOS:**
1. Click the Apple menu and select "About This Mac"
2. Click "System Report..."
3. Under "Hardware", look for "Chip" or "Processor" - if it says "Apple M1", "Apple M2", etc., you have ARM64; otherwise, you have x86_64

**On Linux:**
1. Open a terminal and run: `uname -m`
2. If it shows `x86_64`, you have x86_64 architecture
3. If it shows `aarch64`, you have ARM64 architecture

### Desktop Node Manager Builds

| Operating System | Architecture | Download Link | Package Format | Installation Command |
|------------------|--------------|---------------|----------------|---------------------|
| **Windows** | x86_64 | [DeNode_Manager_1.0.4_x64-setup.exe](https://github.com/DeNetPRO/Node/releases/download/v4.0.1-rc10/DeNode_Manager_1.0.4_x64-setup.exe) | Installer | Double-click to install |
| **macOS** | x86_64 | [DeNode_Manager-1.0.4.dmg](https://github.com/DeNetPRO/Node/releases/download/v4.0.1-rc10/DeNode_Manager-1.0.4.dmg) | Disk Image | Open .dmg file and drag to Applications |
| **macOS** | ARM64 | [DeNode_Manager-1.0.4-arm64.dmg](https://github.com/DeNetPRO/Node/releases/download/v4.0.1-rc10/DeNode_Manager-1.0.4-arm64.dmg) | Disk Image | Open .dmg file and drag to Applications |
| **Linux** | x86_64 | [DeNode_Manager_1.0.4_amd64.deb](https://github.com/DeNetPRO/Node/releases/download/v4.0.1-rc10/DeNode_Manager_1.0.4_amd64.deb) | DEB (Debian/Ubuntu) | `sudo dpkg -i DeNode_Manager_1.0.4_amd64.deb` then `sudo apt install -f` |
| **Linux** | ARM64 | [DeNode_Manager_1.0.4_arm64.deb](https://github.com/DeNetPRO/Node/releases/download/v4.0.1-rc10/DeNode_Manager_1.0.4_arm64.deb) | DEB (Debian/Ubuntu) | `sudo dpkg -i DeNode_Manager_1.0.4_arm64.deb` then `sudo apt install -f` |
| **Linux** | x86_64 | [DeNode_Manager-1.0.4-1.x86_64.rpm](https://github.com/DeNetPRO/Node/releases/download/v4.0.1-rc10/DeNode_Manager-1.0.4-1.x86_64.rpm) | RPM (Red Hat/Fedora/CentOS) | `sudo rpm -ivh DeNode_Manager-1.0.4-1.x86_64.rpm` |
| **Linux** | ARM64 | [DeNode_Manager-1.0.4-1.aarch64.rpm](https://github.com/DeNetPRO/Node/releases/download/v4.0.1-rc10/DeNode_Manager-1.0.4-1.aarch64.rpm) | RPM (Red Hat/Fedora/CentOS) | `sudo rpm -ivh DeNode_Manager-1.0.4-1.aarch64.rpm` |

> ⚠️ **Important**: Desktop Node Manager is not available for server operating systems (Ubuntu Server, CentOS, etc.) due to lack of GUI support.

## Package Format Explanation

- **DEB packages** (ending in `.deb`) are used for Debian, Ubuntu, and other Debian-based Linux distributions
    - Install with: `sudo dpkg -i package.deb` then `sudo apt install -f`

- **RPM packages** (ending in `.rpm`) are used for Red Hat, Fedora, CentOS, and other RPM-based Linux distributions
    - Install with: `sudo rpm -ivh package.rpm`

> 💡 **Tip**: If you're unsure what Linux distribution you're using, run `cat /etc/os-release` in your terminal to find out.

## Installation Options

The DeNet Datakeeper Node offers flexible installation approaches to accommodate different deployment scenarios:

### 1. Command Line Interface (CLI) Node
- **Best for**: Server environments, headless systems, advanced users
- **Features**: Full functionality via terminal commands
- **Installation**: Available for all supported platforms
- **Management**: Requires manual configuration and monitoring

### 2. Desktop Node Manager
- **Best for**: Desktop users, beginners, easy management
- **Features**: Graphical interface for node management
- **Installation**: Available for Windows, macOS, and Linux desktop versions
- **Management**: Intuitive GUI with real-time monitoring

## Requirements

- A wallet address (DeNet app/Metamask/any other wallet) holding a Datakeeper Node License
- Device with free disk space for storing DeNet user data
- Terminal access (Command Prompt/PowerShell on Windows, Terminal on macOS/Linux)
- DeNet Node application downloaded
- PEAQ balance: tokens will be distributed to Datakeepers automatically and will be regularly credited for successfully completed transactions, if the node is running and does not disconnect from the network, no deposits will be required

## Step 0: Verify Your Account Has License

- Open https://peaq.subscan.io/account/YOUR_ADDRESS
- Replace YOUR_ADDRESS with your wallet address
- The license(s) should be seen as a sNL ERC-721 token

![](assets/license.png)

## Step 1: Copy Your Private Key

You need the private key from a wallet with a Datakeeper Node License.

### From DeNet App
1. Open the DeNet app
2. Go to "Profile" -> "Settings" -> "Security"
3. Copy the 64-character HEX private key (e.g., a1b2c3d4...)
4. Save it securely. Never share your private key!

### From Other Wallet (Example: Metamask)
1. Open Metamask in your browser or app
2. Select the account with the Datakeeper Node License
3. Go to "Account Details" > "Export Private Key"
4. Enter your Metamask password and copy the private key
5. **Store it securely! Do not share it anywhere!**

If you use any other wallet, the steps may differ but should be similar to the list above._

## Datakeeper Node Installation

Choose your operating system for detailed installation instructions:

- [Windows](./guides/install-denode-windows.md)
- [Linux](./guides/install-denode-linux.md)
- [MacOS](./guides/install-denode-mac.md)

## Node Manager Installation

For a graphical interface to manage your node, install the Desktop Node Manager application.
Detailed installation instructions can be found in: [Node Manager Installation Guide](./guides/denode-manager.md)

## License Management

🔐 **Managing License Addresses in DeNet**

📌 **Roles for the License**
- **License Owner** - The license holder with full management rights
- **Admin** - Can change the Manager, but cannot transfer the license
- **Manager** - Has a limited set of actions, with no rights in the smart contract

🧭 **Steps for Managing Addresses**

🔗 **Connecting a Wallet**
- Open the [web page manager](https://nodemanager.denet.app/) and connect:
    - Click the **Connect Wallet** button
    - Choose a connection method:
        - Metamask
        - Wallet Connect
- After connecting, your address will appear in the top right corner

📋 **Viewing the List of Licenses**
- After connecting the wallet:
    - Each license is displayed as a card with an ID (e.g., ID #2276) and a **Manage** button

⚙️ **Managing a License**
- Click **Manage** on the desired license. A window will open:
    - "You can change address here"
    - Two fields are available:
        - **New Admin Address** - For changing the administrator
        - **New Manager Address** - For changing the manager

🛠 **Changing Addresses**

▶️ **Changing the Admin**
- Enter the new Ethereum address in the **New Admin Address** field
- Click the **Change Admin** button
- Confirm the transaction in your wallet

▶️ **Changing the Manager**
- Enter the new manager's address in the **New Manager Address** field
- Click **Change Manager**
- Confirm the action in your wallet

📌 **Notes**
- All changes are processed through a smart contract, requiring gas fees
- Ensure the provided addresses are valid (Ethereum format)

#### Ask your questions here and get help:

<a href="https://discord.gg/cPz9m4cSWv">
    <img alt="discord.png" src="https://img.shields.io/badge/Discord-5865F2?style=for-the-badge&logo=discord&logoColor=white" height="30" width="120" />
</a>
