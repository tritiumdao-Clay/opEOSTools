//npx hardhat --network eosevmtest run ./scripts/1CreateL1ERC20.js
const hre = require("hardhat");

async function main() {
    const [owner] = await hre.ethers.getSigners();
    owner.provider.estimateGas = async(transaction) => {
        return hre.config.networks.hardhat.gas;
    }
    console.log("owner adress:", owner.address)
    const TestERC20 = await hre.ethers.getContractFactory("TestERC20")
    const testERC20 = await TestERC20.deploy("TestERC20", "T20")
    await testERC20.deployed()
    console.log("L1 token address:", testERC20.address)
}

main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
});