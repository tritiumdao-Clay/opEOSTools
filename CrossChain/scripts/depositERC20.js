const hre = require("hardhat");
const fs = require("fs");

async function main() {
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

    console.log("--------------------L1 to L2-----------------------")
    let fname2 = "../../optimism/packages/contracts-bedrock/artifacts/src/L1/L1StandardBridge.sol/L1StandardBridge.json"
    let ftext2 = fs.readFileSync(fname2).toString().replace(/\n/g, "")
    let L1StandardBridgeData = JSON.parse(ftext2)
    let l1StandardBridge = new hre.ethers.Contract(process.env.L1StandardBridgeProxy, L1StandardBridgeData.abi, l1Wallet)

    let txL1ETH = await l1StandardBridge.depositETH(1000000, 0x0, {'value':0x1})
    let txL1ERC20 = await l1StandardBridge.depositERC20(l1Addr, l2Addr, 100, 1000000, 0x0)
    console.log("txL1ETH:", txL1ETH)
    console.log("txL1ERC20:", txL1ERC20)

    console.log("--------------------L2 to L1-----------------------")
    let l2Url = process.env.L2_RPC
    let l2RpcProvider = new hre.ethers.JsonRpcProvider(l2Url)
    let l2Wallet = new hre.ethers.Wallet(privateKey, l2RpcProvider)

    let fname3 = "../../optimism/packages/contracts-bedrock/artifacts/src/L2/L2StandardBridge.sol/L2StandardBridge.json"
    let ftext3 = fs.readFileSync(fname3).toString().replace(/\n/g, "")
    let L2StandardBridgeData = JSON.parse(ftext3)
    let l2StandardBridge = new hre.ethers.Contract(process.env.L2StandardBridgeProxy, L2StandardBridgeData.abi, l2Wallet)

    //eth in OptimismPortalProxy, erc20 in L2StandardBridge
    let txL2ETH = await l2StandardBridge.withdraw("0xDeadDeAddeAddEAddeadDEaDDEAdDeaDDeAD0000",1, 1000000, 0x0, {"value":0x1})
    let txL2ERC20 = await l2StandardBridge.withdraw(l2Addr, 1, 1000000, 0x0)
    console.log("txL2ETH:", txL2ETH)
    console.log("txL2ERC20:", txL2ERC20)






    /*
    let fname = "../../optimism/packages/contracts-bedrock/artifacts/src/universal/OptimismMintableERC20Factory.sol/OptimismMintableERC20Factory.json"
    let ftext = fs.readFileSync(fname).toString().replace(/\n/g, "")
    let optimismMintableERC20FactoryData = JSON.parse(ftext)
    console.log(optimismMintableERC20FactoryData)

    let optimismMintableERC20Factory = new ethers.Contract(
        "0x4200000000000000000000000000000000000012",
        optimismMintableERC20FactoryData.abi,
        owner)
    let deployTx = await optimismMintableERC20Factory.createOptimismMintableERC20(
        process.env.L1_TOKEN_ADDRESS,
        "Token Name on L2",
        "L2-SYMBOL"
    )
    console.log("----------------")
    console.log(deployTx)
    let deployRcpt = await deployTx.wait()
    console.log("----------------")
    console.log(deployRcpt)
    console.log("----------------")
    let l1Addr = process.env.L1_TOKEN_ADDRESS
    let event = deployRcpt.events.filter(x => x.event == "OptimismMintableERC20Created")[0]
    let l2Addr = event.args.localToken
    console.log("l1 token:", l1Addr)
    console.log("l2 token:", l2Addr)

     */


}

main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
});
