source .env
cast send --rpc-url $L2_RPC --legacy --cast-async --value 1 --private-key $DEPLOYER_PRIVATE_KEY $L2StandardBridgeProxy "withdraw(address, uint256, uint32, bytes)" $ETH_DEAD_ADDR 1 1000000 ""
cast send --rpc-url $L2_RPC --legacy --cast-async --value 1 --private-key $DEPLOYER_PRIVATE_KEY $L2ToL1MessagePasser "initiateWithdrawal(address, uint256, bytes)" $DEPLOYER_ADDR2 1000000 ""


# 向 $L2StandardBridgeProxy 调用函数 withdraw
# proposer向L1链发布proposer(stateRoot, struct OutputProposal{xxx})
# 监听 L2ToL1MessagePasser 地址 MessagePassed(nonce, sender, target, value, gasLimit, data, withdrawHash)
# 调用L1的 OptimismPortal 合约 的 proveWithdrawalTransaction(xx), finalizeWithdrawalTransaction(xx)


