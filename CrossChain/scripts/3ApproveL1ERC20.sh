# source ./scripts/3ApproveL1ERC20.sh
source .env
cast send --rpc-url $L1_RPC --legacy --cast-async --gas-price 150000000000 --private-key $DEPLOYER_PRIVATE_KEY $L1_TOKEN_ADDRESS "approve(address,uint256)" $L1StandardBridgeProxy 1000000
echo "allowance:"
cast call --rpc-url $L1_RPC --legacy --private-key $DEPLOYER_PRIVATE_KEY $L1_TOKEN_ADDRESS "allowance(address,address)" $DEPLOYER_ADDR  $L1StandardBridgeProxy

