# ERC-721-Checks. Smart contract

Follow the steps below to deploy and manage your smart contract.

## Requirements

### 1. NPM, Node

To successfully build project, generate node_modules and use hardhat, you need to have node and npm installed on your machine.

### 2. Hardhat

Hardhat is a development environment for Ethereum software. It consists of different components for editing, compiling, debugging and deploying your smart contracts and dApps, all of which work together to create a complete development environment. Additional information can be found [here](https://hardhat.org/).

### 3. Crypto Wallet

To create your own smart contract and manage it, you need to have a personal crypto wallet. You can learn how to create one [here](https://metamask.io/).

### 4. Blockchain Node Service Provider

You can use any provider that supports the Ethereum blockchain, such as Infura, Alchemy, etc.

### 5. Etherscan

This is a tool for deploying, verifying, accessing your smart contract functions, and more. You can learn how to create an account [here](https://etherscan.io/).

## Getting started

- Clone the project

```bash
  git clone https://github.com/ArtemBurakov/ERC-721-Checks.git
```

- Go to the project directory

```bash
  cd ERC-721-Checks/server/contract
```

- Install dependencies

```bash
  npm install
```

## Deployment

To deploy your smart contract follow these steps. Run these commands inside `ERC-721-Checks/server/contract` folder.

- Compile contract

```bash
  npx hardhat compile
```

> run `npx hardhat clean` before compile if you need to clear cached files

- Deploy contract

```bash
  npx hardhat run scripts/deploy.js --network goerli
```

> by `--network` argument you can specify which network to use. You can manage networks by modifying `hardhat.config.js` file.

- Verify contract

```bash
  npx hardhat verify --network goerli <DEPLOYED_CONTRACT_ADDRESS>
```

> After deploy, your contract need to be verified in case of future interactions

## Running tests

To run tests run this command inside `ERC-721-Checks/server/contract` folder.

```bash
  npx hardhat test
```

> feel free to rewrite tests inside `/test/checks.test.js` file or use your existing smart contract for testing

## Environment Variables

You will need to add the following environment variables to your `.env` file inside `ERC-721-Checks/server` folder.

`PROVIDER_KEY` - your provider key

`ETHERSCAN_API_KEY` - api key that you will get from Etherscan

`ACCOUNT_PRIVATE_KEY` - your metamask crypto wallet private key

`MINTER_TEST_PRIVATE_KEY` - metamask crypto wallet private key for minter (for tests only)
