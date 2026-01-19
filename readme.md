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

## Table of Contents

- [What is DeNet?](#what-is-denet)
- [Installation Process](#installation-process)
- [Supported Platforms](#supported-platforms)
    - [Command Line Interface (CLI) Node](#command-line-interface-cli-node)
    - [Desktop Node Manager](#desktop-node-manager)
- [Node Activity Monitoring](./guides/monitoring.md)
- [License Management](./guides/license-management.md)

## What is DeNet?

**DeNet** is a decentralized data storage protocol that unlocks the global potential of unused storage. It connects storage providers with those in need of sovereign decentralized storage.

### How It Works:

1. Get a [Datakeeper's license](https://nodesale.denet.app/)
2. [Install special software](#installation-process) on your computer
3. Connect to the DeNet network and share your disk space
4. Store user and business data
5. Get rewarded directly from storage users

More details about requirements can be found [here](./guides/requirements.md)

### Who is a Datakeeper?

💽 Datakeepers offer their unused storage space to DeNet Storage users, creating a mutually beneficial conditions for all without any intermediaries.
By utilizing the DeNet Storage Protocol, users pay for the storage they need, while Datakeepers are rewarded for the storage they provide.

# Installation Process
**Step 1**: Make sure that your setup meets the [requirements](./guides/requirements.md)  
**Step 2**: Choose the **installation option** that better suits you

## Desktop Node Manager
- **Best for**: Desktop users, beginners, easy management
- **Features**: Graphical interface for node management
- **Installation**: Available for Windows, macOS, and Linux desktop versions
- **Management**: Intuitive GUI with real-time monitoring
- **Beginner-friendly**: Ideal for users who prefer visual interfaces and point-and-click operations

#### Installation Guide:

1. **Determine your computer's [architecture](#how-to-choose-the-right-architecture) (x86_64 or ARM64)**
2. **Download the appropriate file** from the [Available Desktop Builds](#available-desktop-builds) section
3. **Install the application**:
    - **Windows**: Double-click the downloaded .exe file and follow the prompts
    - **macOS**: Double-click the .dmg file and drag the app to Applications folder
    - **Linux**: Install with appropriate package manager command
4. **Launch the Desktop Node Manager** and [configure](./guides/configuring-manager.md) your node


## Command Line Interface (CLI) Node
- **Best for**: Server environments, headless systems, advanced users
- **Features**: Full functionality via terminal commands
- **Installation**: Available for all supported platforms
- **Management**: Requires manual configuration and monitoring
- **Beginner-friendly**: Not recommended for users unfamiliar with terminal operations

> 💡 **Tip**: For detailed platform-specific installation instructions, see:
> - [Windows](./guides/install-denode-windows.md)
> - [Linux](./guides/install-denode-linux.md)
> - [MacOS](./guides/install-denode-mac.md)

# Supported Platforms

DeNet Datakeeper Nodes can be installed on various operating systems depending on your needs and technical requirements:

## Command Line Interface (CLI) Node

### Available CLI Builds

| Operating System | Architecture | Download Link |
|------------------|--------------|---------------|
| **Windows** | x86_64 | [denode-windows-amd64.exe](https://github.com/DeNetPRO/Node/releases/download/v4.0.1-rc10/denode-windows-amd64.exe) |
| **Linux** | x86_64 | [denode-linux-amd64](https://github.com/DeNetPRO/Node/releases/download/v4.0.1-rc10/denode-linux-amd64) |
| **Linux** | ARM64 | [denode-linux-arm64](https://github.com/DeNetPRO/Node/releases/download/v4.0.1-rc10/denode-linux-arm64) |
| **macOS** | x86_64 | [denode-macos-amd64](https://github.com/DeNetPRO/Node/releases/download/v4.0.1-rc10/denode-macos-amd64) |
| **macOS** | ARM64 | [denode-macos-arm64](https://github.com/DeNetPRO/Node/releases/download/v4.0.1-rc10/denode-macos-arm64) |

> 📝 **Note**: The DeNet Datakeeper Node CLI is distributed as a standalone binary file without any installation package. Just download the appropriate binary for your system

## Desktop Node Manager

### Available Desktop Builds

| Operating System | [Architecture](#how-to-choose-the-right-architecture) | Download Link                                                                                                                          | [Package Format](#package-format-explanation) | Installation Command |
|------------------|----------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------|-----------------------------------------------|---------------------|
| **Windows** | x86_64                                             | [DeNode_Manager_1.0.4_x64-setup.exe](https://github.com/DeNetPRO/Node/releases/download/v4.0.1-rc10/DeNode_Manager_1.0.4_x64-setup.exe) | Installer                                     | Double-click to install |
| **macOS** | x86_64                                             | [DeNode_Manager_1.0.4_x64.dmg](https://github.com/DeNetPRO/Node/releases/download/v4.0.1-rc10/DeNode_Manager_1.0.4_x64.dmg)            | Disk Image                                    | Open .dmg file and drag to Applications |
| **macOS** | ARM64                                              | [DeNode_Manager-1.0.4-arm64.dmg](https://github.com/DeNetPRO/Node/releases/download/v4.0.1-rc10/DeNode_Manager_1.0.4_x64.dmg)         | Disk Image                                    | Open .dmg file and drag to Applications |
| **Linux** | x86_64                                             | [DeNode_Manager_1.0.4_amd64.deb](https://github.com/DeNetPRO/Node/releases/download/v4.0.1-rc10/DeNode_Manager_1.0.4_amd64.deb)        | DEB (Debian/Ubuntu)                           | `sudo dpkg -i DeNode_Manager_1.0.4_amd64.deb` then `sudo apt install -f` |
| **Linux** | ARM64                                              | [DeNode_Manager_1.0.4_arm64.deb](https://github.com/DeNetPRO/Node/releases/download/v4.0.1-rc10/DeNode_Manager_1.0.4_arm64.deb)        | DEB (Debian/Ubuntu)                           | `sudo dpkg -i DeNode_Manager_1.0.4_arm64.deb` then `sudo apt install -f` |
| **Linux** | x86_64                                             | [DeNode_Manager-1.0.4-1.x86_64.rpm](https://github.com/DeNetPRO/Node/releases/download/v4.0.1-rc10/DeNode_Manager-1.0.4-1.x86_64.rpm)  | RPM (Red Hat/Fedora/CentOS)                   | `sudo rpm -ivh DeNode_Manager-1.0.4-1.x86_64.rpm` |
| **Linux** | ARM64                                              | [DeNode_Manager-1.0.4-1.aarch64.rpm](https://github.com/DeNetPRO/Node/releases/download/v4.0.1-rc10/DeNode_Manager-1.0.4-1.aarch64.rpm) | RPM (Red Hat/Fedora/CentOS)                   | `sudo rpm -ivh DeNode_Manager-1.0.4-1.aarch64.rpm` |

> ⚠️ **Important**: Desktop Node Manager is not available for server operating systems (Ubuntu Server, CentOS, etc.) due to lack of GUI support.

> 💡 **Tip**: Use the Desktop Node Manager download page at https://node.denet.app/ to conveniently download the specific application version.

### How to Choose the Right Architecture

When downloading the DeNet Node Binary Build, it's important to select the version that matches your computer's processor architecture:

- **x86_64 (also called AMD64)**: This is the most common architecture for Intel and AMD processors used in most desktops and laptops
- **ARM64 (also called AArch64)**: This is used by Apple Silicon (M1, M2, etc.) Macs and some ARM-based Linux devices

To determine your system architecture:

**On macOS:**
1. Click the Apple menu and select "About This Mac"
2. Click "More info..."
3. Look for "Chip" or "Processor" - if it says "Apple M1", "Apple M2", etc., you have ARM64; otherwise, you have x86_64

**On Linux:**
1. Open a terminal and run: `uname -m`
2. If it shows `x86_64`, you have x86_64 architecture
3. If it shows `aarch64`, you have ARM64 architecture

### Package Format Explanation

- **DEB packages** (ending in `.deb`) are used for Debian, Ubuntu, and other Debian-based Linux distributions
    - Install with: `sudo dpkg -i package.deb` then `sudo apt install -f`

- **RPM packages** (ending in `.rpm`) are used for Red Hat, Fedora, CentOS, and other RPM-based Linux distributions
    - Install with: `sudo rpm -ivh package.rpm`

> 💡 **Tip**: If you're unsure what Linux distribution you're using, run `cat /etc/os-release` in your terminal to find out.

