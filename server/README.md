# ERC-721-Checks. Golang backend service

Follow the steps below to set up your admin environment.

## Requirements

### 1. Golang

Golang need to be installed on your machine. How to install it can be found [here](https://go.dev/doc/install).

### 2. PostgreSQL

How to install it can be found [here](https://www.cherryservers.com/blog/how-to-install-and-setup-postgresql-server-on-ubuntu-20-04).

### 3. Abigen, Solc

Tools for compiling your smart contract solidity code to use build files to interact with the smart contract.
How to install abigen, follow this [link](https://geth.ethereum.org/docs/tools/abigen).
How to install solc, follow this [link](https://goethereumbook.org/smart-contract-compile/).

### 4. Blockchain Node Service Provider

You can use any provider that supports the Ethereum blockchain, such as Infura, Alchemy, etc.

### 5. Deployed smart contract

Follow this [link](https://github.com/ArtemBurakov/ERC-721-Checks/tree/main/server/contract) to deploy and integrate smart contract in Golang backend service.

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

- Navigate to `ERC-721-Checks/server/cmd/admin` and run this command for starting admin cli. (Before running follow `Compiling smart contract` part).

```bash
  go run main.go
```

- Navigate to `ERC-721-Checks/server/cmd/eventlistener` and run this command for start listening to the smart contract transfer events. (Before running follow `Compiling smart contract` part).

```bash
  go run main.go
```

## Compiling smart contract

To compile your smart contract and get abi follow these steps. Run these commands inside `ERC-721-Checks/server/contract` folder.

- Create abi

```bash
  solc --abi ./contracts/Checks.sol -o build --overwrite
```

- Create bin

```bash
  solc --bin ./contracts/Checks.sol -o build --overwrite
```

- Create smart contract go file for interaction with you contract

```bash
  abigen --abi=./build/Checks.abi --bin=./build/Checks.bin --pkg=checks --out=./../internal/checks/checks.go
```

## Environment Variables

You will need to add the following environment variables to your `.env` file inside `ERC-721-Checks/server` folder.

`DATABASE_HOST` - db host

`DATABASE_PORT` - db port

`DATABASE_NAME` - db name

`DATABASE_USER` - db user

`DATABASE_USER_PASSWORD` - db user password

`TESTNET_PROVIDER` - url from your provider. Can be testnet or mainnet

`SUPER_USER_PRIVATE_KEY` - your metamask crypto wallet private key
