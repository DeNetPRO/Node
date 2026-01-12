# DeNet Datakeeper Node Configuration 

### Follow the instructions in your terminal

1. **Enter private key**: Paste the copied private key and press Enter.
    - The key is stored securely on your device, encrypted with the password.


2. **Set Password**: Enter a strong password
    - The private key is encrypted with the password.    


3. **Choose Port**: Press Enter for the default one
    - Or specify another (value from 10000 to 65535) . 
   

4. **Specify Storage Directory**: Enter path to the user files storage
    - **e.g.**, /home/user/denet_storage (Linux/macOS) or C:\denet_storage (Windows).
    - Ensure the directory exists and has sufficient space.
    - NOTE: use different storage paths for managing different licenses
    - **e.g.**, `/home/user/denet_storage_1` for license with id 1, `/home/user/denet_storage_2` for license with id 2, etc.


5. **Set Storage Space**:
    - Specify the amount of disk space to allocate for DeNet Storage (e.g., 10). Enter the value (only number, without GiB) when prompted.


6. **Optional Second Drive**: Enter 'N' to skip.
    - Or if you want to use another drive, provide its path when prompted.


8. **Select RPC for peaq Blockchain**: Press Enter to use default one.
    - Or choose the RPC endpoint (Select RPC for peaq (ChainID: 3338)).  
   

9. **Verify Operation**:
    - Watch the terminal output. If no errors appear, your DeNet Node is running correctly.
      ![](assets/successful-launch.png)


**NOTE:** Storage path could be changed or added using [CLI commands](./disks-management.md)
