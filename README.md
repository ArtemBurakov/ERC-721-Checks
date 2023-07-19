# ERC-721-Checks

This is a marketplace where users can create NFT tokens with a unique check that defines the purchase. The project provides a tool for creating a secure and transparent token economy or e-commerce platform using NFTs.

## Tech Stack

- Backend: Golang
- Frontend: React
- Database: PostgreSQL
- Decentralized Storage: IPFS node
- Blockchain: Ethereum
- Smart contract: Solidity
- Deployment: Docker and AWS

## Getting Started

- [How to set up backend part](https://github.com/ArtemBurakov/ERC-721-Checks/tree/main/server)
- [How to set up frontend part](https://github.com/ArtemBurakov/ERC-721-Checks/tree/main/client)
- [How to set up your own smart contract and manage it](https://github.com/ArtemBurakov/ERC-721-Checks/tree/main/server/contract)

## Requirements

### 1. Golang

To make this application fully functional, you need to have the Golang backend service. Additional information can be found [here](https://github.com/ArtemBurakov/ERC-721-Checks/tree/main/server).

### 2. React App

To make this application fully functional, you need to run React App. Additional information can be found [here](https://github.com/ArtemBurakov/ERC-721-Checks/tree/main/client).

### 3. Public IPFS node

To make this application fully functional, you need to have your own public IPFS node for saving NFTs metadata. Additional information can be found [here](https://github.com/ArtemBurakov/ERC-721-Checks/tree/main/server). Also feel free to use someone else public node.

### 4. PostgreSQL database

How to install it can be found [here](https://www.cherryservers.com/blog/how-to-install-and-setup-postgresql-server-on-ubuntu-20-04).

### 5. Crypto Wallet

To create your own smart contract and manage it, you need to have a personal crypto wallet. You can learn how to create one [here](https://metamask.io/).

### 6. Blockchain Node Service Provider

You can use any provider that supports the Ethereum blockchain, such as Infura, Alchemy, etc.

### 7. Etherscan

This is a tool for deploying, verifying, accessing your smart contract functions, and more. You can learn how to create an account [here](https://etherscan.io/).

## Deploying to AWS EC2 with Docker

### 1. Docker and Docker compose

You to install docker and docker compose to build images and run containers. Docker Compose using `ERC-721-Checks/docker-compose.yml` file for running all 4 app containers with one command.

- To build app via docker compose run this command inside root folder. Ensure that you are using your env variables in `ERC-721-Checks/docker-compose.yml` file.

```bash
  docker compose up -d
```

> check docker compose documentation for running this command with different params like `-d` for starting app in background. Then you can check running containers with `docker ps` command.

- To use admin and eventlistener cli run these commands

```bash
  docker exec -it <YOUR_GO_CLI_CONTAINER_ID> /bin/bash
  ls
  ./admin or ./eventlistener
```

### 2. AWS EC2

To deploy this application on Amazon AWS follow this [guide](https://everythingdevops.dev/how-to-deploy-a-multi-container-docker-compose-application-on-amazon-ec2/).

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.
