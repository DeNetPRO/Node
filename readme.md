# DeNet Node - CLI app for receiving rewards by sharing unused storage

For using the app you need an account. Account is just an Ethereum wallet that you can import or create.

When you run the app it prompts for password if you already have an account or runs account import command. 

If you run the app for the first time you need to create configuration. It's simple, all you need to do is to answer prompted questions.

Info about the IP address and port that you specified for remote connections will be added to a smart contract. 
#### So there are two things you need:
#### - public IP address
#### - at least 0.1 MATIC on your Ethereum wallet that is going to be used as an account 

Having MATICs in your account wallet is also needed for paying transaction fees when sending file storage proofs, please top it up on time. DeNet smart contracts are deployed in Polygon. More networks will be added in the future.

You can run the app in terminal emulator by navigating to the directory that contains the binary file and typing the following command: 

```bash
./denet-node
```

[Account import command demo](https://www.youtube.com/watch?v=vVRMHlqLA0w)

If you want to create a new account instead of importing an existing one just run the following command:

```bash
./denet-node account create
```

[Account create command demo](https://www.youtube.com/watch?v=So8VAjv9o1Y)

If you generated a new wallet when creating an account and need to receive the private key run the following command: 

```bash
./denet-node account export
```

[Account export command demo](https://www.youtube.com/watch?v=bnstbPGdjKY)

## Currently available commands cheat sheet

| Account Command | Description |
|---|---|
| account create | generates new Ethereum wallet that is used as an account |
| account import | imports Ethereum wallet by its private key |
| account export | discloses your private key |

## Minimal system requirements
- 1 GB of RAM 
- Public IP address and the ability to provide access for remote requests to the device
- 50GB of free storage space.
- 100% Uptime 
- OS: Linux or MacOS (Windows version will be released later)

## FAQ

### Is it ok to run DeNet Node on VPS?

Yes, but ....

If you already use a VPS it is ok to launch DeNet Node on it, but if launching the node is your only purpose, we recommend you trying to launch it on your own machine instead of subscribing on a VPS. So if you have a PC or laptop with Linux or MacOS on board, you'll also need a public IP. It can be obtained via your ISP and is going to cost you cheaper than VPS. 

### How much space should I share?

You can share as much as you can, but sharing more space is making you able to store more files and increases total reward and your chances of receiving it first. 

##  Links

[Discord Miners Group](https://discord.com/channels/920205740944273449/1033033015489728552)


