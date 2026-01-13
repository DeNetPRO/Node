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
- [Datakeeper Node Installation](#datakeeper-node-installation)
- [Node Manager Installation](#node-manager-installation)
- [License Management](#license-management)

# Becoming Datakeeper: Node Setup Instructions

This guide walks you through setting up and running a Datakeeper Node, enabling your device to join the DeNet decentralized storage network.

## Overview

There are two ways to manage your DeNet Node:

1. **CLI Node**: A command-line interface application that runs in your terminal
2. **Desktop Node Manager**: A graphical desktop application for easier management

Download the Desktop Node Manager: [Node Manager Desktop Application](https://node.denet.app/)

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
