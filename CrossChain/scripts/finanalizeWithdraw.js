//Sixth, L1 env
const hre = require("hardhat");

async function main() {
    const [owner] = await hre.ethers.getSigners();
    owner.provider.estimateGas = async(transaction) => {
        return hre.config.networks.hardhat.gas;
    }
    console.log("owner adress:", owner.address)

}

main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
});
