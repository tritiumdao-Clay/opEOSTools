//First, L1 env
//npx hardhat --network eosevmtest run ./scripts/deployERC20.js
//0x82c4F896FfdfFcbcABC8aF9E4966437B91B9B470
const hre = require("hardhat");

async function main() {
    const [owner] = await hre.ethers.getSigners();
    owner.provider.estimateGas = async(transaction) => {
        return hre.config.networks.hardhat.gas;
    }
    console.log("owner adress:", owner.address)
    const TestERC20 = await hre.ethers.getContractFactory("TestDeposit")
    const testERC20 = await TestERC20.deploy()
    await testERC20.waitForDeployment()
    console.log("debug0", await testERC20.getAddress())
}

main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
});
