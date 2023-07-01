const path = require("path");
require("dotenv").config({ path: path.resolve(__dirname, "..", ".env") });
require("@nomiclabs/hardhat-ethers");
require("@nomiclabs/hardhat-etherscan");

const { PROVIDER_KEY, SUPER_USER_PRIVATE_KEY, ETHERSCAN_API_KEY } = process.env;

module.exports = {
  solidity: "0.8.9",
  defaultNetwork: "goerli",
  networks: {
    hardhat: {},
    goerli: {
      url: `https://goerli.infura.io/v3/${PROVIDER_KEY}`,
      accounts: [`0x${SUPER_USER_PRIVATE_KEY}`],
    },
    sepolia: {
      url: `https://sepolia.infura.io/v3/${PROVIDER_KEY}`,
      accounts: [`0x${SUPER_USER_PRIVATE_KEY}`],
    },
    ethereum: {
      chainId: 1,
      url: `https://mainnet.infura.io/v3/${PROVIDER_KEY}`,
      accounts: [`0x${SUPER_USER_PRIVATE_KEY}`],
    },
  },
  etherscan: {
    apiKey: `${ETHERSCAN_API_KEY}`,
  },
};
