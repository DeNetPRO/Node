# Moving a Datakeeper Node to a New Computer

This guide covers moving an existing Datakeeper node — its license, account, and (optionally) already-stored data — from one computer to another, for both the CLI node and the DeNode Manager Desktop.

## Table of Contents

1. [Before You Start](#before-you-start)
2. [Step 1: Stop the Node on the Old Computer](#step-1-stop-the-node-on-the-old-computer)
3. [Step 2: Back Up Your Private Key](#step-2-back-up-your-private-key)
4. [Step 3: Move Your Storage Drive (Optional)](#step-3-move-your-storage-drive-optional)
5. [Step 4: Install the Node Software on the New Computer](#step-4-install-the-node-software-on-the-new-computer)
6. [Step 5: Import Your Account](#step-5-import-your-account)
7. [Step 6: Configure and Start the Node](#step-6-configure-and-start-the-node)

---

## Before You Start

> ⚠️ **Important:** A license can only run actively on one computer at a time. Running the same license on two machines simultaneously causes transaction conflicts, not faster syncing.

Decide what you're moving:
- **Just the setup** (account + config), and you'll let the node re-sync data from scratch — simplest option.
- **The setup and the already-stored data**, by physically moving the drive — keeps your existing stored files, avoids re-download.

> 💡 **Moving only some of your licenses?** If you're keeping the rest running on the old computer, delegate the moved licenses to a separate wallet address via [License Management](./license-management.md) instead of splitting one account across two machines — otherwise you'll hit transaction conflicts.

---

## Step 1: Stop the Node on the Old Computer

- **CLI:** stop the running process (`Ctrl+C` in the terminal, or close the terminal window).
- **DeNode Manager Desktop:** quit the application, then confirm the background `denode` process has actually terminated — closing the window alone may not stop it (see the [reinstallation troubleshooting note](./install-denode-manager-windows.md#reinstallation-issues)).

Do this before activating the license on the new computer, to avoid the transaction conflicts described above.

---

## Step 2: Back Up Your Private Key

You'll need to re-import your account on the new computer.

- **CLI:** `./denode account export` — copy and store the output securely.
- **Other wallet:** use your wallet's export function (e.g., Metamask → Account Details → Export Private Key), as described in the [requirements guide](./requirements.md#step-1-copy-your-private-key).

Never share your private key with anyone — see [What is a private key?](./faq.md#what-is-a-private-key)

---

## Step 3: Move Your Storage Drive (Optional)

If you want to keep your already-stored data, connect the same HDD/SSD to the new computer.

**Can't move the drive physically?** (e.g., it's inside a laptop, or the machines are in different locations) — copy the data over the network instead:
- Same LAN: share the storage folder and copy it with `robocopy` (Windows) or `rsync -a` (Linux/macOS), or mount it over SMB/NFS.
- Remote machines: `rsync -avz` or `scp` over SSH.

The destination folder name doesn't need to match the old one — you set the exact path explicitly in Step 6. Note that for large amounts of stored data, a network copy is usually much slower than physically moving the drive.

---

## Step 4: [Install the Node Software](../readme.md#installation-process) on the New Computer

You're not required to keep using the same node type — e.g., you can switch from CLI to DeNode Manager Desktop on the new machine if that suits you better.

---

## Step 5: Import Your Account

During first-run configuration, choose the **import** option instead of selecting an existing account, and paste the private key from Step 2:

- CLI: see [Configuring DeNode CLI](./configuring-cli.md#follow-the-instructions-in-your-terminal), step 1.
- Manager Desktop: see [DeNode Manager Configuration](./configuring-manager.md#account-setup).

You'll be asked to set a new password — the private key is stored locally, encrypted with it.

---

## Step 6: Configure and Start the Node

Reconfigure the license from scratch on the new computer, rather than trying to carry over the old configurations files — it's simpler and less error-prone. Point the storage directory to your moved drive (Step 3) or a fresh empty folder. Point it to the folder that **contains** the `denode.storage-<license-id>` subdirectory of the old setup — the node keeps stored data inside that subdirectory.

