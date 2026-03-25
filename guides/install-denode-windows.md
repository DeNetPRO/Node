# DeNet Datakeeper Node Installation Guide for Windows

This guide provides simplified step-by-step instructions for installing and running a DeNet Datakeeper Node on Windows systems.

## Table of Contents
1. [Quick Start Installation](#quick-start-installation)
    - [Method 1: Download via curl (Recommended)](#method-1-download-via-curl-recommended)
    - [Method 2: Download from GitHub Website](#method-2-download-from-github-website)
2. [Configuring Node](#configuring-node)
3. [Running Your Node](#running-your-node)
4. [Running Multiple Nodes](#running-multiple-nodes)
5. [Advanced: Windows Service Setup (Optional)](#advanced-windows-service-setup-optional)
6. [Reinstalling Your Node](#reinstalling-your-node)
7. [Useful Commands](#useful-commands)

## Quick Start Installation

### Method 1: Download via curl (Recommended)
```powershell
# Download node executable
curl -LO https://github.com/DeNetPRO/Node/releases/download/v4.0.1-rc12/denode-windows-amd64.exe

# Create directory for the node executable and copy it
mkdir -Force C:\denet
Move-Item -Path "denode-windows-amd64.exe" -Destination "C:\denet\denode.exe"
Set-Location "C:\denet"

# Run the node (this will prompt for configuration)
.\denode.exe
```

### Download from GitHub Website
1. Visit [https://github.com/DeNetPRO/Node/releases](https://github.com/DeNetPRO/Node/releases)
2. Download the appropriate binary for your system:
    - For AMD64 architecture: `denode-windows-amd64.exe`
3. Move the downloaded file to your desired location
   - Using [powershell](#opening-powershell):
      ```powershell
      mkdir -Force C:\denet
      Move-Item -Path "C:\path\to\downloaded\denode-windows-amd64.exe" -Destination "C:\denet\denode.exe"
      Set-Location "C:\denet"
      ```
   
   - Or using File Explorer

## Opening PowerShell

Before proceeding with the installation, you'll need to open PowerShell to execute the commands. Here are several methods to open PowerShell:

### Method 1: Using Windows Search (Recommended for beginners)
1. Press the **Windows key** on your keyboard
2. Type "PowerShell" in the search box
3. Click on **Windows PowerShell** (not "PowerShell ISE")
4. If prompted by User Account Control, click **Yes** to allow the app to make changes

### Method 2: Using Run Dialog (Quick access)
1. Press **Windows key + R** to open the Run dialog
2. Type "powershell" and press **Enter**
3. Alternatively, you can type "powershell.exe" and press Enter

### Method 3: From File Explorer (For executing commands in specific directories)
1. Open File Explorer
2. Navigate to the folder where you downloaded the executable (e.g., C:\denet)
3. Hold **Shift** and right-click in the empty space
4. Select **"Open PowerShell window here"**
5. This opens PowerShell directly in the folder where your executable is located

### Method 4: Using Command Prompt (Alternative approach)
1. Press **Windows key + R**
2. Type "cmd" and press Enter
3. In the Command Prompt window, type "powershell" and press Enter to switch to PowerShell

### Method 5: Task Manager (Advanced users)
1. Right-click on the taskbar and select **Task Manager**
2. Go to the **File** tab and select **Run new task**
3. Type "powershell" and check **Create this task with administrative privileges** if needed
4. Click **OK**

**Important Notes:**
- Always run PowerShell as Administrator if you encounter permission errors during installation
- To run as Administrator, right-click on the PowerShell icon and select **"Run as administrator"**
- The PowerShell window should display a prompt like `PS C:\Users\YourUsername>` indicating it's ready for commands

## Running Your Node

After initial setup of all configuration parameters, simply run:
```powershell
cd C:\denet
.\denode.exe
```

To run in background:
```powershell
$env:DENODE_PASSWORD="your_password"
Start-Process -FilePath ".\denode.exe" -ArgumentList "--address", "your_datakeeper_address", "--license", "your_license_number" -PassThru
```

## Configuring Node
After the node is installed, you can proceed with the configuration

- Instructions are [here](configuring-cli.md)


## Running Multiple Nodes

When running multiple DeNet Datakeeper Nodes on the same machine, each node needs its own configuration directory and unique parameters to avoid conflicts:

1. **Create separate directories for each node:**
   ```powershell
   # Create directories for different nodes
   mkdir -Force C:\denet-node1
   mkdir -Force C:\denet-node2
   Copy-Item -Path "C:\denet\denode.exe" -Destination "C:\denet-node1\denode.exe"
   Copy-Item -Path "C:\denet\denode.exe" -Destination "C:\denet-node2\denode.exe"
   ```

2. **Configure each node separately:**
   ```powershell
   # Configure first node
   Set-Location "C:\denet-node1"
   .\denode.exe  # This will prompt for configuration
   
   # Configure second node
   Set-Location "C:\denet-node2"
   .\denode.exe  # This will prompt for configuration
   ```

3. **Run each node with unique parameters:**
   ```powershell
   # Run first node in background
   Set-Location "C:\denet-node1"
   $env:DENODE_PASSWORD="your_password"
   Start-Process -FilePath ".\denode.exe" -ArgumentList "--address", "your_address_1", "--license", "license_1" -PassThru

   # Run second node in background
   Set-Location "C:\denet-node2"
   $env:DENODE_PASSWORD="your_password"
   Start-Process -FilePath ".\denode.exe" -ArgumentList "--address", "your_address_2", "--license", "license_2" -PassThru
   ```

4. **Monitor processes:**
   ```powershell
   # List all denode processes
   Get-Process denode
   
   # View logs for each node
   Get-Content -Path "C:\denet-node1\denode.log" -Tail 10
   Get-Content -Path "C:\denet-node2\denode.log" -Tail 10
   
   # Kill specific node by PID
   Stop-Process -Id <PID>
   
   # Kill all denode processes
   Get-Process denode | Stop-Process
   ```

**Important:** Each node requires its own unique license id. When configuring each node, make sure to use a different license number for each instance.

## Advanced: Windows Service Setup (Optional)

For automatic startup on boot:

1. Create a dedicated user (recommended):
   ```powershell
   # Create a new user account for the service
   New-LocalUser -Name "denet" -Password (ConvertTo-SecureString "password" -AsPlainText -Force) -Description "DeNet Node Service User"
   ```

2. Create service file:
   ```powershell
   # Create a PowerShell script to run the node
   $serviceScript = @"
   #!/usr/bin/env pwsh
   cd C:\denet
   .\denode.exe
   "@
   $serviceScript | Out-File -FilePath "C:\denet\run-node.ps1" -Encoding UTF8
   ```

3. Install the service:
   ```powershell
   # Install the service using NSSM (Non-Sucking Service Manager)
   # Download NSSM from https://nssm.cc/
   nssm install denode "C:\denet\run-node.ps1"
   ```

4. Start the service:
   ```powershell
   Start-Service denode
   ```

## Reinstalling Your Node

If you need to reinstall your DeNet Datakeeper Node, follow these steps:

1. **Stop the current node process**:
   ```powershell
   Get-Process denode | Stop-Process
   ```

2. **Remove the existing installation directory**:
   ```powershell
   Remove-Item -Recurse -Force C:\denet
   ```

3. **(Optional) Remove Windows service** (if previously configured):
   ```powershell
   # Uninstall the service if it was created with NSSM
   nssm remove denode
   ```

4. **Follow the Quick Start Installation steps** from the beginning to install the node again

## Useful Commands

- Check if node is running: `Get-Process denode`
- Stop node: `Stop-Process -Name denode`
- View logs: `Get-Content -Path "C:\denet\denode.log" -Tail 10`
- Check status: `Get-Process denode`

**Output Example:**
```
PS C:\denet> Get-Process denode
Handles  NPM(K)    PM(K)      WS(K)     CPU(s)     Id  SI ProcessName
-------  ------    -----      -----     ------     --  -- -----------
    123      12     1234       5678      10.5    1234   1 denode

PS C:\denet> Stop-Process -Name denode
PS C:\denet> Get-Process denode
Get-Process : Cannot find a process with the name "denode". Verify the process name and call the cmdlet again.
```
