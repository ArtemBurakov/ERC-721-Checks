# ERC-721-Checks. React App

## Requirements

### 1. NPM, Node

To successfully build project, generate node_modules and run react app, you need to have node and npm installed on your machine.

### 2. Public IPFS node

To make this application fully functional, you need to have your own public IPFS node for saving NFTs metadata Also feel free to use someone else public node.

## Getting started

- Clone the project

```bash
  git clone https://github.com/ArtemBurakov/ERC-721-Checks.git
```

- Go to the project directory

```bash
  cd ERC-721-Checks/client
```

- Install dependencies

```bash
  npm install
```

- Run server locally

```bash
  npm start
```

> before start, create `.env` file with public IPFS node url env variable

## Available Scripts

### `npm start`

Runs the app in the development mode.\
Open [http://localhost:3000](http://localhost:3000) to view it in your browser.

The page will reload when you make changes.\
You may also see any lint errors in the console.

### `npm run build`

Builds the app for production to the `build` folder.\
It correctly bundles React in production mode and optimizes the build for the best performance.

The build is minified and the filenames include the hashes.\
Your app is ready to be deployed!

See the section about [deployment](https://facebook.github.io/create-react-app/docs/deployment) for more information.

## Environment Variables

You will need to add the following environment variables to your `.env` file inside `ERC-721-Checks/client` folder. For starting on localhost only. For production use env variables in `ERC-721-Checks/docker-compose.yml` file.

`REACT_APP_IPFS_URL` - your/someone else ipfs node url
