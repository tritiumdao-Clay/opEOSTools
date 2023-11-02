//npx hardhat --network eosevmtest run ./scripts/3ApproveL1ERC20.js
const hre = require("hardhat");

async function main() {
    const [owner] = await hre.ethers.getSigners();
    owner.provider.estimateGas = async(transaction) => {
    return hre.config.networks.hardhat.gas;
    }
    console.log("owner adress:", owner.address)

    let l1Addr = process.env.L1StandardBridgeProxy
    let erc20Addr = process.env.L1_TOKEN_ADDRESS
    const testERC20 = await hre.ethers.getContractAt("TestERC20", erc20Addr, owner)
    await testERC20.approve(l1Addr, 1000000)
    let balance = await testERC20.allowance(owner.address, l1Addr)
    console.log("allowance:", balance)
}

main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
});
