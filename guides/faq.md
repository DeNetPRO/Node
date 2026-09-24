## FAQ
**Last update:** 2026-09-24

> In the following answers, using the term "node", we mean the **Datakeeper node**.

## Table of Contents

- [General](#general)
- [Hardware & Setup](#hardware--setup)
- [Node Pools](#node-pools)
- [Running & Operations](#running--operations)
- [Errors & Troubleshooting](#errors--troubleshooting)
- [Rewards & Tokens](#rewards--tokens)
- [Uptime & Penalties](#uptime--penalties)
- [Network & Market](#network--market)

## General

### What is a Datakeeper node and how is it different from a Watcher node?

A **Datakeeper node** is a decentralized storage provider that stores file parts uploaded by users and receives rewards directly from those users in exchange for reliable storage. A **Watcher node** is a free mobile node that monitors the network and earns $WN by verifying file copies.

### Can I run a node on a smartphone?

No. Use a PC, laptop, mini PC, or any other device meeting the requirements above.

### Why do I need PEAQ on my balance?

Tokens are required for your node to send transactions on the network and prove that it actually stores data. This is not a reward for storage, but the funds necessary to keep the node running. Your node address is pre-funded with $PEAQ by the DeNet team to cover fees. We will inform you separately about the termination of the reimbursement.

### What is a wallet?

A wallet is like a digital purse that holds your cryptocurrency. For DeNet you need a wallet supporting Ethereum-based (ERC-20) tokens. Your wallet also serves as your identity on the DeNet network and holds your Datakeeper Node License, which is required to participate as a Datakeeper.

### What is a private key?

A private key is a secret password giving you access to your wallet. Keep it safe and never share it: if someone else gets it, they can steal your money or license.

### What is an RPC and why do we use it?

RPC (Remote Procedure Call) is the protocol that lets your node communicate with the DeNet smart contracts on the blockchain. It acts as a bridge so your node can send transactions, read data, submit storage proofs, and participate in consensus. You configure an RPC endpoint to connect your node to the network.

### How will Datakeepers be paid?

Earnings consist of two types of payments — directly from users who store data in the network and from token distribution for early Datakeepers. The number of tokens depends on how long ago you joined the network. Use the calculator on the [Node Sale](https://nodesale.denet.app/) page to estimate earnings.

### Can running nodes be treated as a business?

Yes, but the exact answer depends on your jurisdiction, the amount of generated income, and local regulations. If it qualifies as a business, you may register it as an LLC, sole proprietorship, or another available form. Some jurisdictions offer exemptions (e.g., income thresholds). Consult a legal professional in your jurisdiction.

## Hardware & Setup

### What type of device can run a Datakeeper node?

You can use a PC, laptop, mini PC, or any other device meeting these minimum requirements:
- Operating system: Windows, Linux, or macOS
- RAM: at least 1 GB
- Internet: minimum 10 mb/s

Storage requirements and disk formatting recommendations are below.

### Recommended technical specifications

| Parameter | Minimum | Recommended | Notes |
|---|---|---|---|
| Storage | 100 GB | 2 TB or more | More stored (and proven) data means more rewards. The cap is determined by the Node Sale rules. [Check yours](https://nodesale.denet.app/profile/) |
| Drive cluster size (NTFS allocation unit, Windows only) | 4 KB (OS default) | 64 KB – 1024 KB (1 MB) | The node writes/reads a huge number of small file parts; the 4 KB default adds filesystem overhead and extra I/O. Larger clusters trade some wasted space for fewer, faster I/O operations. Set this when formatting the drive — it can't be changed afterward without reformatting. Applies to Windows/NTFS; Linux (ext4) and macOS (APFS) don't expose an equivalent setting through standard formatting tools. |

### Can I use an external HDD?

Yes. You can run nodes on a computer or NAS, which typically support external drives. The setup is simple with minimal requirements.

### Can I run the same node on two PCs?

No, running the same node on two PCs is impossible.

### Can I run multiple nodes on a single PC?

Yes, you can run several nodes on one PC. This is a strategy to optimize storage space and earn more.

### Can I transfer my data to a new PC without losing it?

Yes. Connect your HDD/SSD to the new PC and launch your node from the new computer. See the full [migration guide](./migrating-node.md) for step-by-step instructions.

### What happens if my drive/device fails or goes offline?

User data stays safe because of multiple copies across DeNet's decentralized network. For your node, an offline period temporarily stops rewards as proof-sending halts, risking node exclusion and penalties. Prolonged offline status could damage your reputation and rating.

### Does my node have to be online 24/7?

To earn rewards, keep your node online. While offline you temporarily stop earning. Short shutdowns for reboots or outages are not critical; prolonged offline status could damage your reputation and rating.

### Do I need a static IP for my node?

No, a static IP is not necessary. Learn more about this in the special [guide](./public-ip.md)

### Can I expand the storage after the initial setup?

Yes. However, some license have limited capacity enabled. You can check how much space you can share on the [Node Sale](https://nodesale.denet.app/profile/) page. Add more drives or expand existing ones to meet your storage goals.

### Are there any ongoing monthly fees to run nodes?

No, there are no monthly fees for running nodes. There are only network fees to process transactions.

### Do I need disk partitions when running multiple nodes?

The system supports both multiple drives per node and multiple nodes per single drive.

### What causes constant disk activity (read/write) on a node — is it the public-node status?

The continuous disk activity comes mainly from replication and client traffic, not from being publicly reachable. Each 1 MB file part is replicated to at least 3 copies across the pool, and every upload/download of a part hits the disk. The routing/tunnel features used by publicly reachable nodes are memory-and-network work, not disk work. So disk wear scales with the volume of data replicated and served, not with the open port.

### Can an NVMe cache drive be added to reduce HDD wear and speed things up?

An NVMe cache reduces HDD wear only when data gets rewritten (a write-back cache collapses repeated writes). A DeNet data part is written exactly once and never modified, so there is nothing to collapse — every byte still lands on the HDD once, so wear is unchanged and just partly shifts to the cache drive. What a cache does help with is read-heavy moments — serving chunks to clients and to peers catching up — and leveling write bursts. Until a cache tier exists, run the node's storage entirely on an NVMe/SSD drive. A fast/cache node tier is noted as a possible future addition.

## Node Pools

### What is a Node Pool in the DeNet ecosystem?

A Node Pool is a group of 32 Datakeeper Nodes that work together as a team. Instead of each node interacting with the entire network, they only communicate within their pool. This reduces delays, speeds up data handling, and lets the network scale as more pools are added.

### How does pool formation currently work?

We're forming optimized sets by multiple characteristics step by step. The first stage introduced a simple distribution of Datakeepers among pools; the next stage will focus on forming targeted collaborations to make that distribution as effective as possible.

At this stage, we're not focusing on the amount of shared storage — the priority is security and data availability, achieved through the widest possible distribution of nodes across pools. This requires careful coordination, since node operators join independently and at their own discretion.

The next step will integrate an algorithm that evaluates node reputation, publicity, and other parameters, each weighted according to its impact. We'll continue refining this algorithm over time, always optimizing for safety and efficiency.

### How do Node Pools help scale the network?

Pools act like multiple lanes on a highway, moving data faster than a single node could. Data flows to active pools and each pool processes it efficiently.

### Why should I run multiple nodes on one device?

Running 5+ nodes on a single device fills your storage faster by pulling data from different pools. It also increases token rewards (one share per active node) and makes efficient use of spare power and space without extra cost.

### What's the advantage of running 10 nodes specifically?

DeNet caps most users at 10 nodes, balancing participation and decentralization. Running 10 maximizes your storage utilization and rewards while strengthening the network — the "sweet spot" for individual users.

### Who can join a Node Pool and what's required?

Anyone with a Datakeeper Node license can join — one license per node, one node per pool. Even data centers can participate under the same rule.

### Can a single-node operator join multiple pools?

No, one node can join only one pool. [Learn more about Pools](https://medium.com/denetpro/how-to-scale-decentralized-storage-introducing-datakeeper-node-pools-230b7167d22f)

### How does the current user-to-pool distribution system work?

Each storage user is assigned to a single pool, and new users are distributed based on available free space. This approach supports balanced data placement across the network, but does not guarantee the same amount of data in each pool and on each node, because we can not control how much data is uploaded by each user.

In the future, this may change due to the introduction of a node rating system and the most prioritized pools. As a result, more users will be redirected to higher-performing pools, and the load will increase.

### Is manual pool selection available?

Pool assignment is fully automatic — there is no manual pool selection at this time.

## Running & Operations

### How do I get started with the Datakeeper node?

Full [text guide]( https://github.com/DeNetPRO/Node/blob/Dev/readme.md) on all applications.

### What if I don't start my node immediately?

You can launch at any time. However, starting earlier lets your node build a higher rating, tracked from June 16, 2025, maximizing its potential.

### Can I transfer my license to another address?

Licenses are currently tied to the address where issued. Use the [License Management function](./license-management.md) to assign License Manager and Admin roles to other wallet addresses.

### How can I monitor that my node is running correctly?

Check your nodes [onchain activity](./monitoring.md) and an official [Datakeeper Console](https://datakeeper-console.denet.pro/dashboard/) interface.

### How do I get the Datakeeper Role on Discord?

**Option 1 (GUI):** Click "Get Discord Role" → Copy your unique code → Go to the [`🤖・datakeeper-role`](https://discord.com/channels/920205740944273449/1514015705639555233) channel → Click "Paste Code" and paste it.

**Option 2 (Terminal):** Run `./denode config code` → select your address → enter your password and copy the code → Go to the [`🤖・datakeeper-role`](https://discord.com/channels/920205740944273449/1514015705639555233) channel → Click "Paste Code" and paste it.

### How does the node operation cycle work?

Each node operates in a repeating cycle with three stages (stage durations are protocol settings that are adjusted over time, so the total cycle length is not fixed):
- **Filling** — the node creates and sends a snapshot of all stored files. If there are no files, nothing is sent to the smart contract and transaction `0x8929ed2f` won't occur.
- **Challenge** — the node sends no transactions and just stores data.
- **Collecting** — the node sends proofs that the files from the Filling stage are still stored. If no snapshot was sent earlier (e.g., no files), there's nothing to prove.

## Errors & Troubleshooting

### Why are my transactions not going through?

Occasional errors are not a reason for drastic measures. If no transactions are being sent at all, check your internet connection, try changing DNS, or restart the computer. If that doesn't help, contact [support](https://discord.com/channels/920205740944273449/1341396814502559846).

### What if I see "Failed to launch node: couldn't join node pool"?

This usually occurs during node initialization and relates to a failed RPC request.
**Action:** Restart the node (resolves most cases). Use a private RPC for better reliability. Check that your account has sufficient PEAQ tokens — verify at https://peaq.subscan.io/

### What are the common node transaction methods?

Your node often uses specific methods to perform different actions. The most common:
- `0xab8d3936` → Joins a Node Pool
- `0xdcfd5bb0` → Confirms Node reachability
- `0x8929ed2f` → Creates a snapshot of the data currently stored
- `0xbd12599d` → Generates a snapshot for a specific user
- `0xf456307e` → Confirms data storage on the network
- `0xbdd859e9` → Confirms storage for a specific user

You can always check your node's performance on https://peaq.subscan.io/

### How do I run my licenses on different computers without transaction errors?

All nodes (versions above v4.1.3-rc1) normally rely on a shared database that helps avoid conflicts when submitting transactions. If licenses are split across different machines, those nodes don't have access to each other's transaction state, which is what's causing the collisions.

The fix is to delegate the affected licenses to a separate wallet address per server, using the [License Management](./license-management.md) feature. For example, if licenses 1 and 2 run on one server, leave them as-is, and delegate 3 and 4 (running on the other server) to a different wallet address. Without this you will keep getting transaction errors.

### What is the general recommendation to avoid errors?

Always run the latest node version for optimal stability. Download the newest releases [here](https://github.com/DeNetPRO/Node/releases)

### Windows Defender flags denode.exe as a threat — how to fix?

This is a false positive, a known issue. The exclusion is the standard fix; a permanent Microsoft-side resolution is being worked on.

**Step 1. Restore from quarantine:** Windows Security → Virus & threat protection → Protection history → locate `denode.exe` → select **Restore** or **Allow**.

**Step 2. Add an exclusion:** Virus & threat protection → Manage settings → scroll to **Exclusions** → **Add an exclusion** → **Folder** → enter `C:\Users\YourUser\AppData\Local\DeNode Manager` (replace YourUser with your Windows username). This prevents Defender from scanning or removing files in the app folder.

**Step 3. Reinstall (recommended):** Uninstall DeNode Manager, choosing "Delete the application data" for a clean state. Restart the PC. With the folder exclusion already in place, install DeNode Manager cleanly.

After this, your node should run without Windows Defender interference. If issues persist, open a support ticket and attach your logs.

### What can cause previously stored data to become unavailable?

The most common causes are:

- The node was reassigned to a different pool after a restart — data associated with the previous pool is no longer tracked under the node's new assignment.
- The end user who owned the data deleted it themselves — the Datakeeper stops receiving proof challenges and rewards for that data because it no longer exists on the user's side.

The amount of data on a node is a dynamic metric that can fluctuate for many reasons on the user's side. It is not always synchronized with payment calculations at the moment, so you should not perceive a short-term decrease in volume as a direct indicator of income.

## Rewards & Tokens

### How will ingress/egress be compensated when data flows through a node but doesn't land on it?

In the future we plan to introduce traffic-based payments. Users will pay not only for storage but also for data movement. Public nodes handling the highest traffic volumes will receive greater rewards.

### Do withdrawals happen per-node or pooled?

Tokens are credited to the Datakeeper's balance. If you own several licenses, you withdraw the tokens earned by all of them together.

### Where does the liquidity for the TBY pool come from?

At present, all tokens received from sales are automatically transferred into the liquidity pool for Datakeepers (except for a small protocol commission from every sale).

### Can I expand the 10 TB limit to a higher or unlimited license?

The mechanics of increasing the limit are possible but have not yet been highlighted to the community.

### Is there potential for staking in exchange for additional licenses or storage capacity?

We'll announce it if something of this kind appears.

### How is network revenue generated and distributed?

Coming soon: live earnings calculation and forecasting will be built directly into official [Datakeeper Console](https://datakeeper-console.denet.pro/dashboard), so every Datakeeper can track their real-time share and projected growth in one place.

### What determines the difference in earnings between Datakeepers?

Earnings come from real user payments for storage. When a node successfully submits Proof-of-Storage, it receives TBY directly from the users whose data it stores in the pool in which it was defined.

The more data a node stores and successfully proves, the higher its earnings.

## Uptime & Penalties

### How does the penalty system work?

If a node doesn't send any proof during a full cycle, it receives 1 penalty. After 10 penalties (10 missed cycles; how long that takes depends on the current stage durations), the protocol considers the node offline and exits it from the pool. After restarting, the node automatically creates a join transaction to re-enter an available pool. If a node has fewer than 10 penalties (e.g., 5) and then sends a valid proof, the counter resets to 0.

### Do penalties ever reset or are they permanent?

Yes. As soon as the node submits a valid proof, the 10-penalty cycle before ejection restarts for that node. But the conditions can be changed depending on the pool.

## Network & Market

### What will long-term Datakeepers (1+ year, running from the beginning) receive vs. those who joined later?

Early Datakeepers unlock DeNet tokens, forming an integral part of their earnings, withdrawable after TGE. The actual withdrawable amount will be determined by a node's rating, based on key metrics including uptime and successful Proofs-of-Storage.

### What is the current supply/demand ratio — customers to Datakeepers?

There's no priority for one side over the other. We aim to maintain a constant balance between supply and demand so the network can thrive.
