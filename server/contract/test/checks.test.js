require("dotenv").config();
require("@nomicfoundation/hardhat-chai-matchers");

const { expect } = require("chai");
const { ethers } = require("hardhat");
const { MINTER_TEST_PRIVATE_KEY } = process.env;

const gasLimit = 3000000;
const tokenURI = "ipfs://QmRAQB6YaC91dP37UdDnjFY5vQuiBrcqdyoW1Cu321wxkD4";

describe("Checks", function () {
  this.timeout(100000);

  let user;
  let admin;
  let checksContract;
  let minterWithProvider;

  before(async function () {
    const provider = ethers.provider;

    admin = await ethers.getSigner();
    user = await ethers.Wallet.createRandom();
    minterWithProvider = new ethers.Wallet(
      MINTER_TEST_PRIVATE_KEY,
      provider
    ).connect(provider);

    const Checks = await ethers.getContractFactory("Checks");
    checksContract = await Checks.deploy();
    await checksContract.deployed();

    // Use deployed contract
    // const contractAddress = checksContract.address;
    // console.log("Contract address:", contractAddress);
    // const Checks = await ethers.getContractFactory("Checks");
    // checksContract = await Checks.attach(contractAddress);
  });

  describe("Deployment", function () {
    it("Should set the admin as the default admin role", async function () {
      expect(
        await checksContract.hasRole(
          await checksContract.DEFAULT_ADMIN_ROLE(),
          admin.address
        )
      ).to.be.true;
    });

    it("Should set the admin as a minter", async function () {
      expect(
        await checksContract.hasRole(
          await checksContract.MINTER_ROLE(),
          admin.address
        )
      ).to.be.true;
    });
  });

  describe("Minting", function () {
    it("Should allow the admin to set a minter", async function () {
      await checksContract
        .connect(admin)
        .setMinter(minterWithProvider.address)
        .then((tx) => tx.wait());

      expect(
        await checksContract.hasRole(
          await checksContract.MINTER_ROLE(),
          minterWithProvider.address
        )
      ).to.be.true;
    });

    it("Should not allow the default user to set a minter", async function () {
      expect(
        await checksContract
          .connect(minterWithProvider)
          .setMinter(user.address, { gasLimit })
      ).to.be.revertedWith("Caller is not an admin");
    });

    it("Should not allow setting an invalid minter address", async function () {
      await expect(
        checksContract.connect(admin).setMinter(ethers.constants.AddressZero)
      ).to.be.revertedWith("Invalid minter address");
    });

    it("Should allow a minter to mint a token", async function () {
      await checksContract
        .connect(minterWithProvider)
        ._mint(user.address, tokenURI)
        .then((tx) => tx.wait());
      expect(await checksContract.ownerOf(0)).to.equal(user.address);
      expect(await checksContract.tokenURI(0)).to.equal(tokenURI);
    });

    it("Should not allow minting to an invalid recipient address", async function () {
      expect(
        await checksContract
          .connect(minterWithProvider)
          ._mint(ethers.constants.AddressZero, tokenURI, { gasLimit })
      ).to.be.revertedWith("Invalid recipient address");
    });

    it("Should not allow minting to the contract itself", async function () {
      expect(
        await checksContract
          .connect(minterWithProvider)
          ._mint(checksContract.address, tokenURI, { gasLimit })
      ).to.be.revertedWith("Cannot mint to the contract itself");
    });

    it("Should not allow minting with an empty URI", async function () {
      expect(
        await checksContract
          .connect(minterWithProvider)
          ._mint(user.address, "", { gasLimit })
      ).to.be.revertedWith("URI cannot be empty");
    });
  });
});
