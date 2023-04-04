/**
 * @type import('hardhat/config').HardhatUserConfig
 */

require("dotenv").config();
require("@nomiclabs/hardhat-ethers");

const { PROVIDER_KEY, ACCOUNT_PRIVATE_KEY } = process.env;

module.exports = {
  solidity: "0.8.1",
  //defaultNetwork: "sepolia",
  networks: {
    hardhat: {},
    // sepolia: {
    //   url: `https://sepolia.infura.io/v3/${PROVIDER_KEY}`,
    //   accounts: [`0x${ACCOUNT_PRIVATE_KEY}`],
    // },
    ethereum: {
      chainId: 1,
      url: `https://mainnet.infura.io/v3/${PROVIDER_KEY}`,
      accounts: [`0x${ACCOUNT_PRIVATE_KEY}`],
    },
  },
};
