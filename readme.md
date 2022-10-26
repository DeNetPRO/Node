# DeNet Node - a CLI application for sharing unused storage capacity for reward.

For using the app you need an account. Account is just an Ethereum wallet that you can import or create.

When you run the app it prompts for password if you already have an account or runs account import command. 

If you run the app for the first time you need to create configuration. It's simple, all you need to do is to answer prompted questions.

Info about the IP address and port that you specified for remote connections will be added to a smart contract. 
#### So there are two things you need:
#### - public IP address
#### - at least 0.1 MATIC on your Ethereum wallet that is going to be used as an account

Having MATICs in your account wallet is also needed for paying transaction fees when sending file storage proofs, please top it up on time.

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

[Account export command demo](https://www.youtube.com/watch?v=bnstbPGdjKY)

## Currently available commands cheat sheet

| Account Command | Description |
|---|---|
| account create | generates new Ethereum wallet that is used as an account |
| account import | imports Ethereum wallet by its private key |
| account export | discloses your private key |

