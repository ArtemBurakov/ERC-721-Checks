# ERC-721-Checks. Golang backend service

Follow the steps below to set up your admin environment.

## Requirements

### 1. Golang

Golang need to be installed on your machine. How to install it can be found [here](https://go.dev/doc/install).

### 2. Abigen, Solc

Tools for compiling your smart contract solidity code to use build files to interact with the smart contract.
How to install abigen, follow this [link](https://geth.ethereum.org/docs/tools/abigen).
How to install solc, follow this [link](https://goethereumbook.org/smart-contract-compile/).

### 3. Blockchain Node Service Provider

You can use any provider that supports the Ethereum blockchain, such as Infura, Alchemy, etc.

## Getting started

- Clone the project

```bash
  git clone https://github.com/ArtemBurakov/ERC-721-Checks.git
```

- Go to the project directory

```bash
  cd ERC-721-Checks/server
```

- Install dependencies

```bash
  go mod tidy
```

- Navigate to `ERC-721-Checks/server/cmd` and run this command for starting admin cli. (Before running follow `Compiling smart contract` part).

```bash
  go run server.go
```

- Navigate to `ERC-721-Checks/server/listener` and run this command for start listening to the smart contract transfer events. (Before running follow `Compiling smart contract` part).

```bash
  go run listener.go
```

## Compiling smart contract

To compile your smart contract and get abi follow these steps. Run these commands inside `ERC-721-Checks` folder.

- Create abi

```bash
  solc --abi ./contracts/<CONTRACT_FILE_NAME>.sol -o build --overwrite
```

- Create bin

```bash
  solc --bin ./contracts/<CONTRACT_FILE_NAME>.sol -o build --overwrite
```

- Create smart contract go file for interaction with you contract

```bash
  abigen --abi=./build/<CONTRACT_FILE_NAME>.abi --bin=./build/<CONTRACT_FILE_NAME>.bin --pkg=contract --out=./../server/contract/<CONTRACT_FILE_NAME>.go
```

## Environment Variables

You will need to add the following environment variables to your `.env` file inside `ERC-721-Checks/server` folder.

`DEPLOYED_CONTRACT_ADDRESS` - your deployed contract address

`TESTNET_PROVIDER` - url from your provider. Can be testnet or mainnet

`SUPER_USER_PRIVATE_KEY` - your metamask crypto wallet private key
