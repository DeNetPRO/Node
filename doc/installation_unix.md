Linux & MacOS installation 
------------------

This guide shows you how to launc datakeeper on your Linux / MacOS machine.\
If you still have questions, you can ask it in discord or telegram support channels\
</br><img src="https://img.shields.io/badge/Telegram-2CA5E0?style=for-the-badge&logo=telegram&logoColor=white" height="30" width="120"/> 
<img src="https://img.shields.io/badge/Discord-5865F2?style=for-the-badge&logo=discord&logoColor=white" height="30" width="120"/> 


## Requiremnts:
- Public & static IP address
- At least 0.1 Polygon MATIC (in order to send proofs)
- 1GiB of RAM
- Datakeeper ID rented ([how to?](./doc/get_id.md))
- Stable internet connection
- Provided free disk space > 512 GiB
- Operating system: Linux, MacOS, Windows (all x64)

<!--
    Hello, this is step by step instruction on how to start your first DeNet Datakeeper node with Linux, and start earning money, providing your free disk space.
    First of all you should have Datakeeper ID created, watch how to do it in our last video, link will be down below.

    So, let's get started !
-->
### Installation proccess
------------------
<!--
    First, we should have node installed, so run this command. All commands will be also attached below.
--> 

1. Run next command to install node locally:

```console
sh -c "$(curl -fsSL https://raw.githubusercontent.com/denetpro/node/master/scripts/install.sh)"
```
    or using wget

```console
sh -c "$(wget  -O - https://raw.githubusercontent.com/denetpro/node/master/scripts/install.sh)"
```

2. Launch denode binary
```console
denode
```

3. Next follow the instructions and input next fields:

- Export Private key, with datakeeper ID rented from it

- Set the passowrd to protect your private key

- Your public IP address

- Port or use default 55050 (press enter)

- Share existing free partition (disk) dirpath

    - Path to existing directory

    - Size of this partition (in GiB)
- Choose blockchain (currently only polygon is available)

    - Press enter on polygon

    - Choose custom RPC, or use default one (press enter)


4. Complete


#### If you see similliar logs, you can be happy and wait for first files and rewards
</br><img src="imgs/start_logs.png" width="550"/>

#### 🪙 Withdraw $TBYmined to $DE token
- use to influence protocol parameters in [Consensus](https://consensus.denet.app/#welcome_to_consensus)
- Use it for your own fate

#### 🪙 Earn $TBY as a system reward from protocol  (2% APY of $TBY total supply)
- use it for storing data
- soon: convert it to $TBYmined using charity proofs. Read how it will work [here](https://medium.cojm/denetpro/denet-storage-protocol-v3-to-address-key-challenge-of-decentralization-f19b9041b0fa#:~:text=close%20the%20deposit.-,Charity%20Proof,-%3A%20The%20DeNet)

### Also...
#### ❗️ Make sure you always have enough $MATIC on yur balance, in order to send proofs.
#### ❗️ Don't let your $TBY balance become less than your Datakeeper ID, unless you want to stop earning rewards.
