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
