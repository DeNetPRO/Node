# DeNet Datakeeper
<a href="https://denet.pro">
    <img alt="tg.png" src="https://img.shields.io/badge/website-000000?style=for-the-badge&logo=About.me&logoColor=white" height="31" width="120" href="https://discord.gg/cPz9m4cSWv"/>
</a> 
<a href="https://t.me/+Yu5KnSruttc5ZGRi">
    <img alt="tg.png" src="https://img.shields.io/badge/Telegram-2CA5E0?style=for-the-badge&logo=telegram&logoColor=white" height="30" width="120" href="https://discord.gg/cPz9m4cSWv"/>
</a> 
<a href="https://discord.gg/cPz9m4cSWv">
    <img alt="discord.png" src="https://img.shields.io/badge/Discord-5865F2?style=for-the-badge&logo=discord&logoColor=white" height="30" width="120" />
</a>
<a href="https://www.youtube.com/channel/UCeCxt3tYbtSkJvaznNjQimQ">
    <img alt="youtube.png" src="https://img.shields.io/badge/YouTube-FF0000?style=for-the-badge&logo=youtube&logoColor=white" height="30" width="120" href="https://discord.gg/cPz9m4cSWv"/>
</a>

[Download latest release](https://github.com/DeNetPRO/Node/releases/latest)

To use the app you need an account. The account is just an Ethereum wallet that you can import or create.
When you run the app it prompts you for password if you already have an account otherwise it runs the account import command.
If you run the app for the first time you need to setup the configuration. It's simple, all you need to do is to answer the question prompts.

Information about the IP address and port that you specified for remote connections will be added to a smart contract.


- [Requirements](#requiremnts)
- [Installation](#installation)
- [Usage](#usage)
- [Video tutorials](#video-tutorials)
- [FAQ](#faq)


## Requiremnts:
- Public & static IP address
- At least 0.1 Polygon MATIC (in order to send proofs)
- 1GiB of RAM
- Datakeeper ID ([how to get it?](./doc/get_id.md))
- Stable internet connection
- Provided free disk space > 512 GiB
- Operating system: Linux, MacOS, Windows (all x64)


## Installation
💡 Make sure you meet all the requirements, and you've already received datakeeper ID [here](./doc/get_id.md).

- [Linux/MacOS installation](./doc/installation_unix.md)
- [Windows installation](./doc/installation_windows.md)

## Usage
| Account management| Description | Demo video |
|---|---|----|
| account create | create new Ethereum wallet (if didn't exist)| [video](https://www.youtube.com/watch?v=So8VAjv9o1Y) |
| account import | imports Ethereum wallet by its private key | [video](https://www.youtube.com/watch?v=vVRMHlqLA0w) |
| account export | show your private key | [video](https://www.youtube.com/watch?v=bnstbPGdjKY)|
| account delete | delete your current account | - |

## Video tutorials
- [Launching datakeeper on MacOS](https://www.youtube.com/watch?v=vw7yyDjyhS8)

## FAQ 

### Is it ok to run DeNet Node on VPS ❓

▪️ Yes, but ....

If you already use a VPS it is ok to launch DeNet Node on it, but if launching the node is your only purpose, we recommend you try to launch it on your own machine instead of subscribing on a VPS. So if you have a PC or laptop with Linux or MacOS on board, you'll also need a public IP. It can be obtained via your ISP and is going to be cheaper than VPS.

### How much space should I share ❓

▪️ You can share as much as you can but sharing out more space will allow you to store more files and increases reward totals and your chances of being the first to receive it.

### Does the disk have to be ssd ❓

▪️ There are no special requirements, but the faster your I/O speed, the more files you are able to store simultaneously and send proofs. However, your hard drive must be in working condition without any damage.

### How fast should the internet speed be ❓

▪ Recommended speed is 100Mb/sec, and the connection should be stable.

### How can I calculate my earnings ❓

▪ Nodes are rewarded in TBY tokens. One TBY can be received for storing 1 TB of data for a year. If you stored 500 GB for a year you will receive 0.5 TBY and so on. The inflation rate will be 2% out of TBY Total Supply. Daily system reward is shared among all Datakeepers. 
You can also calculate your expected income in more detail depending on your Node ID in [Dune Analytics](https://dune.com/djdeniro/denet-v3). 

### What Node ID will I get❓

▪️You can take a look at current node supply in [contract](https://polygonscan.com/token/0xcb19bede3e4f64b6b0085d99127f6d0a25b7180d).
That means if there are already 5 nodes in the contract, you will be the sixth to register there and your ID is going to be number 6 in this contract. That also means that you must have 6 TBY on your balance in order to register your node.

### When should I expect to get an income ❓

▪️ You will start earning TBY_Mined when files are uploaded to your node and you find proof for these files (based on the DeNet Proof-of-Storage algorithm).

### How much can I earn for providing 1TiB ❓

▪️ The approximate reward for providing 1TB is $100. But we're improving our rewarding system and creating conditions for more profitable storage usage. 

