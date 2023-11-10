#source ./scripts/2MintL1ERC20.sh
source .env
cast send --rpc-url $L1_RPC --legacy --cast-async --gas-price 150000000000 --private-key $DEPLOYER_PRIVATE_KEY $L1_TOKEN_ADDRESS "mint(address,uint256)" $DEPLOYER_ADDR 10000
#cast send --rpc-url $L1_RPC --legacy --cast-async --gas-price 150000000000 --private-key $DEPLOYER_PRIVATE_KEY $L1_TOKEN_ADDRESS "mint(address,uint256)" 0xC9DB2f861dec8d63bC2203E6CC727b6cff233d74 10000000000000000000000

