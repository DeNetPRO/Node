# DeNet Datakeeper Node Installation Guide for Windows

This guide provides simplified step-by-step instructions for installing and running a DeNet Datakeeper Node on Windows systems.

## Quick Start Installation

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

## Configuring Node
- Instructions [here](./configuring.md)


After the initial installation, you can run the node in background:
```powershell
$env:DENODE_PASSWORD="your_password"
Start-Process -FilePath ".\denode.exe" -ArgumentList "--address", "your_datakeeper_address", "--license", "your_license_number" -PassThru
```

## Useful Commands

- Check if node is running: `Get-Process denode`
- Stop node: `Stop-Process -Name denode`

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
