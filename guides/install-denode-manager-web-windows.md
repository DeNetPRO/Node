# DeNode Manager Web for Windows

## Installation Guide

### Step 0: Download Application
1. Download the archive for your system from https://github.com/DeNetPRO/Node/releases
   ### Windows
    ```
    denode-manager-windows-amd64.zip
    ```

### Step 1: Install And Run
1. Extract the .zip to a folder of your choice (e.g. `C:\denode-manager`)

    The folder contains the Manager server (`server.exe`), the node binary (`denode.exe`) and the web interface (`static\`) — keep them together: the server looks for the node binary and the interface next to itself.
2. Run `server.exe` — it starts the backend locally

### Step 2: Open Application Interface in Browser
1. Open browser and go to http://localhost:1111
   ![Node GUI](../assets/node-gui.png)

## Notes:
1. Launched server should always be running in the background, otherwise the application will not work, check the status using the application's interface
2. You shouldn't use both CLI and GUI at the same time, otherwise you will get an undefined application behaviour
