package main

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"os"
)

type GethTxn struct {
	To       string `json:"to"`
	From     string `json:"from"`
	Gas      string `json:"gas"`
	GasPrice string `json:"gasPrice"`
	Value    string `json:"value"`
	Data     string `json:"input"`
}

func main() {
	privateKeyECDSA, err := crypto.HexToECDSA("f9b70dd856352559c44ba6c17b91b502197b78390e7253d4aef2e76032b4a683")
	if err != nil {
		fmt.Println("load private key", err.Error())
		panic(err)
	}
	publicKey := privateKeyECDSA.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		fmt.Println("load public key")
		panic(err)
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Println("address:", fromAddress)

	fmt.Println("debug0")
	pwd, _ := os.Getwd()
	fmt.Println("debug0:", pwd)
	L1StandardBridgeData, err := os.ReadFile(pwd + "/abi/L1StandardBridge.json")
	if err != nil {
		fmt.Println("load L1StandardBridgeData abi fail")
		panic(err)
	}

	contractAbi, err := abi.JSON(bytes.NewReader(L1StandardBridgeData))
	fmt.Println(contractAbi.Methods["depositERC20"])
	//data, err := contractAbi.Pack("depositETH", uint32(100), []byte(""))
	//function depositERC20(address _l1Token, address _l2Token, uint256 _amount, uint32 _minGasLimit, bytes _extraData) returns()
	l1Addr := common.HexToAddress("0xc7ffb803bfC73E59a9C8a201CAB28C5a0Cb2Da96")
	l2Addr := common.HexToAddress("0x0c2ffeba77ab37eec68b09fd2ae1acbd205cc7b7")
	amount, _ := new(big.Int).SetString("100", 10)
	gasLimit := uint32(1000000)
	_ = amount
	data, err := contractAbi.Pack("depositERC20", l1Addr, l2Addr, amount, gasLimit, []byte(""))
	if err != nil {
		fmt.Println("package fail:", err.Error())
		panic(err)
	}
	fmt.Println("data:", hex.EncodeToString(data))

	//os.Exit(1)
	client, err := ethclient.Dial("https://api.testnet.evm.eosnetwork.com")
	if err != nil {
		fmt.Println("connect server fail:", err.Error())
		panic(err)
	}
	nonce, err := client.NonceAt(context.Background(), fromAddress, nil)
	if err != nil {
		fmt.Println("get nonce fail:", err.Error())
		panic(err)
	}
	gasPrice, _ := new(big.Int).SetString("150000000000", 10)
	value, _ := new(big.Int).SetString("0", 10)
	gas := uint64(1000000)

	fmt.Println("debug4")
	contractAddress := common.HexToAddress("0x433c8294AAB8027e1c20b1389C55283B67a640F5")

	fmt.Println("debug5")
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		Gas:      gas,
		GasPrice: gasPrice,
		To:       &contractAddress,
		Value:    value,
		Data:     data,
	})
	fmt.Println("debug6")

	chainId, _ := new(big.Int).SetString("15557", 10)
	signer := types.NewEIP155Signer(chainId)

	fmt.Println("debug7")
	sigTx, err := types.SignTx(tx, signer, privateKeyECDSA)
	if err != nil {
		fmt.Println("sign fail:", err.Error())
		panic(err)
	}
	fmt.Println("debug8")
	fmt.Println("hash:", sigTx.Hash())
	rawTrans, err := sigTx.MarshalBinary()
	if err != nil {
		fmt.Println("marsh fail:", err.Error())
		panic(err)
	}
	fmt.Println("marsh fail:", hex.EncodeToString(rawTrans))

	fmt.Println("debug9")
	err = client.SendTransaction(context.Background(), sigTx)
	if err != nil {
		fmt.Println("broadcast fail:", err.Error())
		panic(err)
	}

}
