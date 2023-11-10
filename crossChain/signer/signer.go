package signer

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type Signer interface {
	Address() common.Address
	SignerFn(chainID *big.Int) bind.SignerFn
	SignData([]byte) ([]byte, error)
}

func CreateSigner(privateKey string) (Signer, error) {
	if privateKey != "" {
		key, err := crypto.HexToECDSA(privateKey)
		if err != nil {
			return nil, fmt.Errorf("error parsing private key: %w", err)
		}
		return &ecdsaSigner{key}, nil
	}
	return nil, errors.New("private key not config")
}
