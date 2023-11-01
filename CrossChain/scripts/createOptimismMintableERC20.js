const hre = require("hardhat");
const fs = require("fs");

async function main() {
    const [owner, userAddr] = await hre.ethers.getSigners();
    console.log("owner addr:", await owner.getAddress())

    owner.provider.estimateGas = async(transaction) => {
        return hre.config.networks.hardhat.gas;
    }

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
    console.log("l2 token:", l2Addr)


}

main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
});
