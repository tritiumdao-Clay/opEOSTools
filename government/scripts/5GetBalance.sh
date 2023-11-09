#Fifth,
#source ./scripts/getBalance.sh
source .env
echo -n $DEPLOYER_ADDR" eth balance(L1):"
cast rpc eth_getBalance $DEPLOYER_ADDR "latest"  --rpc-url $L1_RPC
echo -n $DEPLOYER_ADDR" eth balance(L2):"
cast rpc eth_getBalance $DEPLOYER_ADDR2 "latest"  --rpc-url $L2_RPC
echo -n $DEPLOYER_ADDR" erc20 balance(L1):"
curl --location $L1_RPC --header 'Content-Type: application/json' --data '{"jsonrpc":"2.0","method":"eth_call","params": [{"data":"0x70a08231000000000000000000000000'$DEPLOYER_ADDR'", "to":"'$L1_TOKEN_ADDRESS'"}, "latest"],"id":67}'
echo -n $DEPLOYER_ADDR" erc20 balance(L2):"
curl --location $L2_RPC --header 'Content-Type: application/json' --data '{"jsonrpc":"2.0","method":"eth_call","params": [{"data":"0x70a08231000000000000000000000000'$DEPLOYER_ADDR'", "to":"'$L2_TOKEN_ADDRESS'"}, "latest"],"id":67}'
