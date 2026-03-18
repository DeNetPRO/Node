# DeNet Desktop Node Manager Installation Guide for Windows

This guide provides step-by-step instructions for installing and running the DeNet Desktop Node Manager on Windows systems. The Desktop Node Manager offers a user-friendly graphical interface for managing your DeNet Datakeeper nodes.

## Table of Contents

1. [Quick Start](#quick-start)
2. [Download the Application](#download-the-application)
3. [Installation](#installation)
4. [Launch DeNode Manager](#launch-denode-manager)
5. [Troubleshooting](#troubleshooting)
6. [Next Steps](#next-steps)

---

## Quick Start

> ⚠️ **Important:** Before starting, please review the [System and Account Requirements](./requirements.md) to ensure you have everything needed.

Get started with DeNet Desktop Node Manager in just a few steps:

1. **Download:** Get the `.exe` file from the [GitHub Releases page](https://github.com/DeNetPRO/Node/releases)
2. **Install:** Double-click the installer and follow the wizard
3. **Launch:** Open DeNode Manager from the Start menu or Desktop
4. **Configure:** Follow the setup wizard to import your wallet account and activate your node

That's it! Your node will be ready to run after configuration.

For detailed instructions, see:
- [Configuring Node Manager Guide](./configuring-manager.md) - Complete configuration steps
- [Node Activity Monitoring](./monitoring.md) - Monitor your node performance
- [License Management](./license-management.md) - Manage your Datakeeper licenses

---

## Download the Application

### Method 1: Direct Download from GitHub (Recommended)

1. Visit the [DeNet Node Releases page](https://github.com/DeNetPRO/Node/releases)
2. Download the Windows version:
   - **Windows (x86_64):** `DeNode_Manager_1.0.5_x64-setup.exe`

### Method 2: Using PowerShell

Open PowerShell and run:

```powershell
# Download the installer
curl -LO https://github.com/DeNetPRO/Node/releases/download/v4.0.1-rc11/DeNode_Manager_1.0.5_x64-setup.exe
```

### Check Your System Architecture

```powershell
# Check your system architecture
wmic os get OSArchitecture

# Should show "64-bit" or "x64"
```

---

## Installation

### Step 1: Run the Installer

1. Navigate to your Downloads folder
2. Double-click `DeNode_Manager_1.0.5_x64-setup.exe`
3. If prompted by User Account Control, click **Yes** to allow the installation

[image: win-run.png]

### Step 2: Follow the Installation Wizard

1. Choose your installation language and click **OK**
2. Click **Next** on the welcome screen
3. Accept the license agreement
4. Choose the installation location (default is recommended)
5. Click **Install** to begin the installation

[image: win-folder.png]

### Step 3: Complete Installation

1. Wait for the installation to complete
2. Click **Finish** to exit the wizard
3. Optionally, check "Launch DeNode Manager" to start the application immediately

> 💡 **Tips:** 
> - The installer will create a Start Menu shortcut and optionally a Desktop shortcut
> - You can change the installation location if needed, but the default `C:\Program Files\DeNode Manager` is recommended

---

## Launch DeNode Manager

### Opening the Application

**From Start Menu:**
1. Click the **Start** button
2. Find **DeNode Manager** in the application list
3. Click to launch

**From Desktop:**
1. Double-click the **DeNode Manager** shortcut on your Desktop (if created during installation)

**From Installation Folder:**
1. Navigate to `C:\Program Files\DeNode Manager`
2. Double-click `DeNode Manager.exe`

On first launch, you will see the initial configuration screen:

[image: windows_screen_1.png]

The application will display a list of licenses associated with your wallet. Initially, all licenses will be disabled until configured.

> 💡 **Tip:** The first launch may take some time as the application initializes its components.

---

## Troubleshooting

Contact support if you encounter issues:
- [Discord Support](https://discord.gg/cPz9m4cSWv)

---

## Next Steps

Congratulations! You have successfully installed and launched the DeNet Desktop Node Manager. Here's what to do next:

### 1. Configure Your Node

Learn how to configure and activate your node:
- [Configuring Node Manager Guide](./configuring-manager.md)

### 2. Monitor Your Node

Learn how to monitor node activity and performance:
- [Node Activity Monitoring](./monitoring.md)

### 3. Manage Your License

Understand license management and transactions:
- [License Management](./license-management.md)

### 4. Set Up Public IP (Optional)

Improve node performance with a public IP address:
- [Public IP Setup Guide](./public-ip.md)

---

## Additional Resources

### Documentation

- [FAQ - Frequently Asked Questions](./faq.md)
- [Command Line Interface Guide](./denode-command.md)
- [Disk Management Guide](./disks-management.md)

### Installation Guides for Other Platforms

- [macOS Installation](./install-denode-manager-macos.md)
- [Linux Installation](./install-denode-manager-linux.md)

### Community & Support

- **Official Website:** [denet.pro](https://denet.pro)
- **Discord Server:** [Join our community](https://discord.gg/cPz9m4cSWv)

---

**Need Help?**  
If you encounter any issues not covered in this guide, please reach out through our community channels:
- [Discord Support](https://discord.gg/cPz9m4cSWv)
