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

[The latest release](https://github.com/DeNetPRO/Node/releases/latest)

## Provide your storage space & earn $80 for 1TB!

💽  Datakeepers offer their unused storage space to DeNet Storage users, creating a mutually beneficial conditions for all without any intermediaries. By utilizing the DeNet Storage Protocol, users pay for the storage they need, while Datakeepers are rewarded for the storage they provide. Rent out your computer, HDD or old unused devices and benefit!

💰 [Calculate your potential reward](https://p.denet.app/datakeeper)

Follow the navigation to become a Datakeeper:

- [Requirements](#requirements)
- [Step 1: Download](#download)
- [Step 2: Installation process](#installation)
- [Step 3: Get your Datakeeper ID](#id)
- [Step 4: Send proofs and get rewards](#rewards)
- [FAQ](#faq)


## Requirements:
- Public & static IP address
- At least 0.1 Polygon MATIC (in order to send proofs)
- 1GiB of RAM
- Stable internet connection, minimum speed - 20Mb/sec
- Provided free disk space > 512 GiB
- Operating system: Linux, MacOS, Windows (all x64)

## Download

💡 **Make sure you meet all the requirements first.**

As the step 1, you need to download DeNode for your OC. You can always find the latest release [here](https://github.com/DeNetPRO/Node/releases/latest)

## Installation

To use the app you need an account. The account is just an Ethereum wallet that you can import or create. When you run the app it prompts you for password if you already have an account otherwise it runs the account import command. If you run the app for the first time you need to setup the configuration. It's simple, all you need to do is to answer the question prompts.

Information about the IP address and port that you specified for remote connections will be added to a smart contract.

| Account management| Description | Demo video |
|---|---|----|
| account create | create new Ethereum wallet (if didn't exist)| [video](https://www.youtube.com/watch?v=So8VAjv9o1Y) |
| account import | imports Ethereum wallet by its private key | [video](https://www.youtube.com/watch?v=vVRMHlqLA0w) |
| account export | show your private key | [video](https://www.youtube.com/watch?v=bnstbPGdjKY)|
| account delete | delete your current account | - |

**Installation process:**

1. Run next command to install node locally:

Using curl:
```console
$ sh -c "$(curl -fsSL https://raw.githubusercontent.com/denetpro/node/master/scripts/install.sh)"
```
Using wget:
```console
$ sh -c "$(wget  -O - https://raw.githubusercontent.com/denetpro/node/master/scripts/install.sh)"
```

2. Launch denode binary
```console
denode
```

3. Follow  the instructions and input next fields:

- Export Private key

- Set the password to protect your private key

- Your public IP address

- Port or use default 55050 (press enter)

- Share existing free partition (disk) dirpath

    - Path to existing directory

    - Size of this partition (in GiB)
- Choose blockchain (currently only polygon is available)

    - Press enter on polygon

    - Choose custom RPC, or use default one (press enter)


4. Complete


#### If you see similliar logs, congratulations! 🎉 You will soon get your first files and rewards.
![image](https://github.com/Arthurbrain/Node/assets/143201292/abdd3eec-a7d4-45d4-8cc3-04e2b98504d7)

## Get your Datakeeper ID

To be part of the network, you should get a Datakeeper ID. To get an ID you need to make a TBY deposit. The amount of TBY necessary to launch each new node is determined by the number of existing Datakeepers. 

For example, if there are already 10 Datakeepers in the network, your ID will be assigned as number 11, along with the necessary TBY deposit. It is important to note that the deposit remains your own property and it's not spent!

## Rewards

Send your first proofs and earn $TBY and $TBYmined as a reward.

**Earn $TBYmined for sending proofs. Withdraw $TBYmined to $DE token**
- Use to influence protocol parameters at Consensus
- Use it for your own benefit

**Earn $TBY as a system reward from the protocol (2% APY of $TBY total supply)**
- Soon: convert it to $TBYmined with charity proofs.
- Use it for storing data


## FAQ 

### Can I run DeNet Node on a VPS  ❓

▪️ Yes, it is possible to run DeNode on a VPS. However, we recommend running the node on your own machine if it is your only purpose. If you have a PC or laptop with Linux or MacOS, obtaining a public IP from your ISP will be cheaper than using a VPS.

### Do I need an SSD for my hard drive  ❓

▪️ There are no specific requirements for your hard drive, but faster I/O speed will allow you to store more files simultaneously and send proofs. Your hard drive must be in working condition without any damage.

### What Internet speed do I need ❓

▪️ The minimum speed required is 20Mb/sec. However, higher and more stable internet speeds will provide better income opportunities.

### How do I calculate my potential earnings ❓

▪ DeNet Datakeepers have a stable income for sending proofs, and there are additional incentive programs. Visit [DeNet Payments](https://p.denet.app/datakeeper) and use our calculator to estimate your potential income.

### When will I start earning income ❓

▪ You will start earning TBY_Mined when files are uploaded to your node, and you send proof of storing these files based on the DeNet Proof-of-Storage algorithm.

### How much can I earn for providing 1TiB ❓

▪️ The approximate reward for providing 1TB is $80, but our rewarding system is constantly improving to create more profitable storage usage conditions.

### What is the purpose of getting a Datakeeper ID ❓

▪️ To become a Datakeeper and run a DeNode, you must make a TBY deposit equal to the ID you'll take up. The deposit improves network stability, enhances upload and download speeds, and provides a secure environment for all members. For more information, please refer to our [article](https://medium.com/denetpro/denet-storage-protocol-v3-to-address-key-challenge-of-decentralization-f19b9041b0fa).

### Why do I need TBY on my balance ❓
▪️ You can only send proofs and earn rewards if your TBY balance ≥ your Datakeeper ID.

### Will my deposit be spent when running a node ❓

▪️ No, the deposit amount remains the same and is only used to assert your ID.

### What can I do with my deposit if I want to leave ❓

▪️ Your deposit remains your property and can be withdrawn at any time. However, note that your ID may be taken by another member after you leave.
