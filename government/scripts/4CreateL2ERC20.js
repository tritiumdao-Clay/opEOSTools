//npx hardhat --network eosevmtest run ./scripts/1CreateL1ERC20.js
const hre = require("hardhat");

async function main() {
    const [owner] = await hre.ethers.getSigners();
    console.log("owner adress:", owner.address)

    let abi = [
        " function createStandardL2Token(address _remoteToken,string memory _name,string memory _symbol) external returns (address)"
    ]
    const ADDR = "0x4200000000000000000000000000000000000012"
    const testERC20 = await hre.ethers.getContractAt(abi, ADDR, owner)

    const L1ADDR = process.env.L1_TOKEN_ADDRESS
    tx = await testERC20.createStandardL2Token(L1ADDR, "L2Name", "L2Symbol")
    console.log("L2 token trans:", tx)
}

main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
});