//Fifth,
const hre = require("hardhat");

async function main() {
    const [owner] = await hre.ethers.getSigners();
    owner.provider.estimateGas = async(transaction) => {
        return hre.config.networks.hardhat.gas;
    }
    console.log("owner adress:", owner.address)

    let l1Url = process.env.L1_RPC
    let l1RpcProvider = new hre.ethers.JsonRpcProvider(l1Url)
    let privateKey = process.env.DEPLOYER_PRIVATE_KEY
    let l1Wallet = new hre.ethers.Wallet(privateKey, l1RpcProvider)

    let l2Url = process.env.L2_RPC
    let l2RpcProvider = new hre.ethers.JsonRpcProvider(l2Url)
    let l2Wallet = new hre.ethers.Wallet(privateKey, l2RpcProvider)

    let l1Factory = await hre.ethers.getContractFactory("TestERC20")
    let l1Contract = new hre.ethers.Contract(process.env.L1_TOKEN_ADDRESS, l1Factory.interface, l1Wallet)
    let l2Contract = new hre.ethers.Contract(process.env.L2_TOKEN_ADDRESS, l1Factory.interface, l2Wallet)

    let beforeL1ETHBalance = await l1Wallet.getBalance()
    console.log("L1ETHBalance:", beforeL1ETHBalance)
    let beforeL1ERC20Balance = await l1Contract.balanceOf(owner.address)
    console.log("L1ERC20Balance:", beforeL1ERC20Balance)

    let beforeL2ERC20Balance = await l2Contract.balanceOf(owner.address)
    console.log("L2ERC20Balance:", beforeL2ERC20Balance)
    let beforeL2ETHBalance = await l2Wallet.balanceOf(owner.address)
    console.log("L2ETHBalance", beforeL2ETHBalance)
}

main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
});
{

}
