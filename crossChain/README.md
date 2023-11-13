

```javascript
//ETH_DEAD_ADDR=0xDeadDeAddeAddEAddeadDEaDDEAdDeaDDeAD0000
//L2StandardBridgeProxy=0x4200000000000000000000000000000000000010
cast send --rpc-url $L2_RPC --legacy --cast-async --value 1 --private-key $DEPLOYER_PRIVATE_KEY $L2StandardBridgeProxy "withdraw(address, uint256, uint32, bytes)" $ETH_DEAD_ADDR 1 1000000 ""
```

```javascript
// prove tx:
go run main.go --rpc $L1_RPC --network opeostest --start-http false --private-key $KEY --withdrawal $TXHASH
```

```javascript
// finalize tx:
go run main.go --rpc $L1_RPC --network opeostest --start-http false --private-key $KEY --withdrawal $TXHASH
```
