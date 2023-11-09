const hre = require("hardhat");
const fs = require("fs");
const optimismSDK = require("@eth-optimism/sdk");

async function main() {
    //reference: https://github.com/ethereum-optimism/optimism-tutorial/tree/main/standard-bridge-standard-token
    const [owner, userAddr] = await hre.ethers.getSigners();
    owner.provider.estimateGas = async(transaction) => {
        return hre.config.networks.hardhat.gas;
    }
    console.log("owner addr:", await owner.getAddress())

    const l1Addr = process.env.L1_TOKEN_ADDRESS,
        l2Addr = process.env.L2_TOKEN_ADDRESS;

    let l1Url = process.env.L1_RPC
    let l1RpcProvider = new hre.ethers.JsonRpcProvider(l1Url)
    let privateKey = process.env.DEPLOYER_PRIVATE_KEY
    let l1Wallet = new hre.ethers.Wallet(privateKey, l1RpcProvider)

    let fname2 = "../../optimism/packages/contracts-bedrock/artifacts/src/L1/L1StandardBridge.sol/L1StandardBridge.json"
    let ftext2 = fs.readFileSync(fname2).toString().replace(/\n/g, "")
    let L1StandardBridgeData = JSON.parse(ftext2)
    let l1StandardBridge = new hre.ethers.Contract(process.env.L1StandardBridgeProxy, L1StandardBridgeData.abi, l1Wallet)

    console.log("--------------------L2 to L1-----------------------")
    let l2Url = process.env.L2_RPC
    let l2RpcProvider = new hre.ethers.JsonRpcProvider(l2Url)
    let l2Wallet = new hre.ethers.Wallet(privateKey, l2RpcProvider)

    let fname3 = "../../optimism/packages/contracts-bedrock/artifacts/src/L2/L2StandardBridge.sol/L2StandardBridge.json"
    let ftext3 = fs.readFileSync(fname3).toString().replace(/\n/g, "")
    let L2StandardBridgeData = JSON.parse(ftext3)
    let l2StandardBridge = new hre.ethers.Contract(process.env.L2StandardBridgeProxy, L2StandardBridgeData.abi, l2Wallet)

    //eth in OptimismPortalProxy, erc20 in L2StandardBridge
    //eth: 0xDeadDeAddeAddEAddeadDEaDDEAdDeaDDeAD0000
    //erc20: erc20Addr
    let txL2ETH = await l2StandardBridge.withdraw("0xDeadDeAddeAddEAddeadDEaDDEAdDeaDDeAD0000",1, 1000000, 0x0, {"value":0x1})
    let txL2ERC20 = await l2StandardBridge.withdraw(l2Addr, 1, 1000000, 0x0)
    console.log("txL2ETH:", txL2ETH)
    console.log("txL2ERC20:", txL2ERC20)

    console.log("---------------------------finalize-withdraw------------")
    let l1ChainId = 15557
    let l2ChainId = 42096
    let crossChainMessenger = new optimismSDK.CrossChainMessenger({
        l1ChainId: l1ChainId,
        l2ChainId: l2ChainId,
        l1SignerOrProvider: l1Wallet,
        l2SignerOrProvider: l2Wallet
    })
    //Initiate the withdrawal on L2
    let withdrawalTx1 = await crossChainMessenger.withdrawERC20(l1Addr, l2Addr, 1)
    await withdrawalTx1.wait()

    console.log("----------------------------")
    //Wait until the root state is published on L1, and then prove the withdrawal. This is likely to take less than 240 seconds.
    await crossChainMessenger.waitForMessageStatus(withdrawalTx1.hash, optimismSDK.MessageStatus.READY_TO_PROVE)
    let withdrawalTx2 = await crossChainMessenger.proveMessage(withdrawalTx1.hash)
    await withdrawalTx2.wait()
    console.log("-------------")
    await l1Contract.balanceOf(l1Wallet.address)
    await l2Contract.balanceOf(l1Wallet.address)

}

main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
});

