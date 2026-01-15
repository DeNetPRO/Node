# DeNet Datakeeper Node Configuration 

### Follow the instructions in your terminal

1. **Choose the account or import a new one**:
   - To use already imported account: Choose it in the displayed list
   - To import: Choose the **import** option and paste the copied [private key](./requirements.md#step-1-copy-your-private-key), then press Enter

2. **Set Password**: Enter a strong password
    - The private key will be stored securely on your device, encrypted with this password

3. **Choose License ID**: enter license bumber you'd like to start
    - Enter any license number from the **Available licenses** list
    > ⚠️ **Important**: Do not choose the same license number when running multiple node instances. This is not prohibited, but it can lead to undefined behavior

4. **Enter IP Address**: Press enter to use the default one (0.0.0.0)
    - If you'd like to set up the [public IP](./public-ip.md), enter it in the following format: xxx.xxx.xxx.xxx
   
5. **Choose Port**: Press Enter to use default one (55050)
    - Or specify another (value from 10000 to 65535)
    - If you use the public IP address, make sure that the port is [forwarded](./public-ip.md#port-forwarding-requirements)

6. **Specify Storage Directory**: Enter path to the user files storage
    - **e.g.**, /home/user/denet_storage (Linux/macOS) or C:\denet_storage (Windows)
    - Ensure the directory exists and has sufficient space

7. **Set Storage Space**:
    - Specify the amount of disk space to allocate for DeNet Storage (e.g., 10). Enter the value (only number, without GiB) when prompted.
   
8. **Optional Second Drive**: Enter 'N' to skip.
    - Or if you want to use another drive, provide its path when prompted.

9. **Select [RPC](./faq.md#what-is-an-rpc-and-why-do-we-use-it) for peaq Blockchain**: Press Enter to use default one.
    - Or choose the RPC endpoint (Select RPC for peaq (ChainID: 3338)).  

10. **Verify Operation**:
    - Watch the terminal output. If no errors appear, your DeNet Node is running correctly.
      ![](assets/successful-launch.png)


**NOTE:**
1. Configuration parameters could be viewed or changed using [CLI config commands](./denode-command.md#config-management-commands)
2. Shared space parameters could be changed using [CLI disks commands](./disks-management.md)
