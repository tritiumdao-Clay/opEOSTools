require("@nomicfoundation/hardhat-toolbox");

/** @type import('hardhat/config').HardhatUserConfig */
require("hardhat-contract-sizer");
require("@openzeppelin/hardhat-upgrades");
require("hardhat-gas-reporter");
require("@nomiclabs/hardhat-web3-legacy");
const dotenv = require("dotenv");
dotenv.config();

module.exports = {
  solidity: {
    version: "0.8.20",
    settings: {
      optimizer: {
        enabled: true,
        runs: 200,
        details: {
          peephole: true,
          jumpdestRemover: true,
          orderLiterals: false,
          deduplicate: true,
          constantOptimizer: true,
          yul: true,
        },
      },
    },
  },
  networks: {
    hardhat: {
      //forking: {
      //  url: `http://127.0.0.1:8545`,
      //  blockNumber: 19316651,
      //},
      //forking: {
      //  url: `http://127.0.0.1:8545`,
      //  blockNumber: 100,
      //},
      accounts: [{
        "privateKey": process.env.DEPLOYER_PRIVATE_KEY,
        "balance": "10000000000000000000000000",
      }],
    },
    eosevmtest: {
      url: `https://api.testnet.evm.eosnetwork.com`,
      accounts: [process.env.DEPLOYER_PRIVATE_KEY],
      gasPrice: 150000000000,
      gas: 1000000,
    },
    eosevm: {
      url: `https://api.evm.eosnetwork.com`,
      accounts: [process.env.DEPLOYER_PRIVATE_KEY],
      gasPrice: 150000000000,
      gas: 1000000,
    },
    local: {
      url: `http://127.0.0.1:8545`,
      accounts: [process.env.DEPLOYER_PRIVATE_KEY],
      gasPrice: 150000000000,
      gas: 1000000,
    },
  },
  contractSizer: {
    alphaSort: true,
    disambiguatePaths: false,
    runOnCompile: false,
    strict: true,
  },
  gasReporter: {
    enabled: false,
    currency: "CHF",
    gasPrice: 21,
  },
  mocha: {
    timeout: 100000,
  },
};
