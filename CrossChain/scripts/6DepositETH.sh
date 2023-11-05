#Fourth,
# source depositAndWithdraw.sh
source .env
cast send --rpc-url $L1_RPC --legacy --cast-async --gas-price 150000000000 --value 1 --private-key $DEPLOYER_PRIVATE_KEY $L1StandardBridgeProxy "depositETH(uint32,bytes)" 1 ""

