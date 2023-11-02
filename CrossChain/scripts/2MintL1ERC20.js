// npx hardhat --network eosevmtest run ./scripts/2MintL1ERC20.js
const hre = require("hardhat");

async function main() {
    const [owner] = await hre.ethers.getSigners();
    owner.provider.estimateGas = async(transaction) => {
    return hre.config.networks.hardhat.gas;
    }
    console.log("owner address:", owner.address)

    let l1Addr = process.env.L1_TOKEN_ADDRESS
    const testERC20 = await hre.ethers.getContractAt("TestERC20", l1Addr, owner)
    await testERC20.mint(owner.address, 10000)
    let balance = await testERC20.balanceOf(owner.address)
    console.log("owner balance:", balance)
}

main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
});
