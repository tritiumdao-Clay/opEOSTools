const hre = require("hardhat");
const fs = require("fs");

async function main() {
    const [owner, userAddr] = await hre.ethers.getSigners();
    console.log("owner addr:", await owner.getAddress())
    owner.provider.estimateGas = async(transaction) => {
        return hre.config.networks.hardhat.gas;
    }

    const l1Addr = process.env.L1_TOKEN_ADDRESS,
        l2Addr = process.env.L2_TOKEN_ADDRESS;

    let fname = "../../optimism/packages/contracts-bedrock/artifacts/src/universal/OptimismMintableERC20.sol/OptimismMintableERC20.json"
    let ftext = fs.readFileSync(fname).toString().replace(/\n/g, "")
    let optimismMintableERC20Data = JSON.parse(ftext)
    let l2Contract = new hre.ethers.Contract(l2Addr, optimismMintableERC20Data.abi, owner)

    let l1Url = process.env.L1_RPC
    let l1RpcProvider = new hre.ethers.JsonRpcProvider(l1Url)
    let privateKey = process.env.DEPLOYER_PRIVATE_KEY
    let l1Wallet = new hre.ethers.Wallet(privateKey, l1RpcProvider)

    let l1Factory = await hre.ethers.getContractFactory("TestERC20")
    let l1Contract = new hre.ethers.Contract(process.env.L1_TOKEN_ADDRESS, l1Factory.interface, l1Wallet)
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
