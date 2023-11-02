package main

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
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

	L1StandardBridgeData, err := os.ReadFile("./abi/L1StandardBridge")
	if err != nil {
		fmt.Println("load L1StandardBridgeData abi fail")
		panic(err)
	}

	contractAbi, err := abi.JSON(bytes.NewReader(L1StandardBridgeData))
	amount, _ := new(big.Int).SetString("100", 10)
	data, err := contractAbi.Pack("depositETH", common.HexToAddress(toAddress), amount)
	if err != nil {
		fmt.Println("contractAbi.Pack error ,", err)
		return "", err
	}

}
func SignTxn(from string, _to string, data []byte, nonce uint64, value int64, gas *big.Int, gasPrice *big.Int, privkey *ecdsa.PrivateKey) (*GethTxn, error) {

	var parsed_tx = new(GethTxn)
	var amount = big.NewInt(value)
	var bytesto [20]byte
	_bytesto, _ := hex.DecodeString(_to[2:])
	copy(bytesto[:], _bytesto)
	to := common.Address([20]byte(bytesto))

	signer := types.NewEIP155Signer(nil)
	tx := types.NewTransaction(nonce, to, amount, gas, gasPrice, data)
	signature, _ := crypto.Sign(tx.SigHash(signer).Bytes(), privkey)
	signed_tx, _ := tx.WithSignature(signer, signature)

	json_tx, _ := signed_tx.MarshalJSON()
	_ = json.Unmarshal(json_tx, parsed_tx)
	parsed_tx.From = from
	fmt.Println("data", parsed_tx.Data)
	return parsed_tx, nil
}
