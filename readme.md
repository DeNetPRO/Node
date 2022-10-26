# DeNet Node - a CLI application for sharing unused storage capacity for reward.

## Distributed binary usage

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

[account import](https://www.youtube.com/watch?v=vVRMHlqLA0w)

If you want to create a new account instead of importing an existing one just run the following command:

```bash
./denet-node account create
```

[account create](https://www.youtube.com/watch?v=So8VAjv9o1Y)



[export key](https://www.youtube.com/watch?v=bnstbPGdjKY)



## Commands
Open a terminal / console in the folder where you downloaded DeNet-Node and run the command below. You will see a list of available commands that you can use:

```bash
./DeNet-Node --help
``` 

There is an example how to use the commands
```bash
./DeNet-Node [command]
```

| Available Commands | Description |
|---|---|
| account | Account is a command that lets you to manage accounts in the different blockchain network |
| config | Config is a command that lets you change your account configuration |
| help | Help about any command |

```bash
./DeNet-Node account [command]
```

| Account Command | Description |
|---|---|
| create | create a new blockchain account |
| import | imports your account by private key |
| key | discloses your private key |
| login | log in a blockchain accounts |

```bash
./DeNet-Node config [command]
```

| Config Command | Description |
|---|---|
| update | updates your account configuration |


## API Documentation
[Documentation](https://app.gitbook.com/o/-MhCmHmTRDb0MF2vIQKk/s/-MhI3_4Kt2DnLxDFkDH8)
