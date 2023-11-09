package withdraw

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum-optimism/optimism/op-bindings/bindings"
	"github.com/ethereum-optimism/optimism/op-node/withdrawals"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

func TxBlock(ctx context.Context, l2c *rpc.Client, l2TxHash common.Hash) (*big.Int, error) {
	l2 := ethclient.NewClient(l2c)
	receipt, err := l2.TransactionReceipt(ctx, l2TxHash)
	if err != nil {
		return nil, err
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		return nil, errors.New("unsuccessful withdrawal receipt status")
	}
	return receipt.BlockNumber, nil
}

func ProofFinalized(ctx context.Context, portal *bindings.OptimismPortal, l2TxHash common.Hash) (bool, error) {
	return portal.FinalizedWithdrawals(&bind.CallOpts{}, l2TxHash)
}

func ProvenWithdrawal(ctx context.Context, l2c *rpc.Client, portal *bindings.OptimismPortal, l2TxHash common.Hash) (struct {
	OutputRoot    [32]byte
	Timestamp     *big.Int
	L2OutputIndex *big.Int
}, error) {
	empty := *new(struct {
		OutputRoot    [32]byte
		Timestamp     *big.Int
		L2OutputIndex *big.Int
	})

	l2 := ethclient.NewClient(l2c)
	receipt, err := l2.TransactionReceipt(ctx, l2TxHash)
	if err != nil {
		return empty, err
	}

	ev, err := withdrawals.ParseMessagePassed(receipt)
	if err != nil {
		return empty, err
	}

	hash, err := withdrawals.WithdrawalHash(ev)
	if err != nil {
		return empty, err
	}

	return portal.ProvenWithdrawals(&bind.CallOpts{}, hash)
}

func ProveWithdrawal(ctx context.Context, l1 *ethclient.Client, l2c *rpc.Client, l2oo *bindings.L2OutputOracle, portal *bindings.OptimismPortal, l2TxHash common.Hash, opts *bind.TransactOpts) error {
	l2 := ethclient.NewClient(l2c)
	l2g := gethclient.New(l2c)

	l2OutputBlock, err := l2oo.LatestBlockNumber(&bind.CallOpts{})
	if err != nil {
		return err
	}

	l2OutputIndex, err := l2oo.GetL2OutputIndexAfter(&bind.CallOpts{}, l2OutputBlock)
	if err != nil {
		return err
	}

	header, err := l2.HeaderByNumber(ctx, l2OutputBlock)
	if err != nil {
		return err
	}
	params, err := withdrawals.ProveWithdrawalParameters(ctx, l2g, l2, l2TxHash, header, &l2oo.L2OutputOracleCaller)
	if err != nil {
		return err
	}

	tx, err := portal.ProveWithdrawalTransaction(
		opts,
		bindings.TypesWithdrawalTransaction{
			Nonce:    params.Nonce,
			Sender:   params.Sender,
			Target:   params.Target,
			Value:    params.Value,
			GasLimit: params.GasLimit,
			Data:     params.Data,
		},
		l2OutputIndex,
		params.OutputRootProof,
		params.WithdrawalProof,
	)
	if err != nil {
		return err
	}

	fmt.Printf("Proved withdrawal for %s: %s\n", l2TxHash.String(), tx.Hash().String())

	// Wait 5 mins max for confirmation
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	return waitForConfirmation(ctxWithTimeout, l1, tx.Hash())
}

func ProveWithdrawal2(ctx context.Context, l1 *ethclient.Client, l2c *rpc.Client, l2oo *bindings.L2OutputOracle, portal *bindings.OptimismPortal, l2TxHash common.Hash) (string, error) {
	l2 := ethclient.NewClient(l2c)
	l2g := gethclient.New(l2c)

	fmt.Println("debug20")
	l2OutputBlock, err := l2oo.LatestBlockNumber(&bind.CallOpts{})
	if err != nil {
		return "", err
	}

	fmt.Println("debug21")
	l2OutputIndex, err := l2oo.GetL2OutputIndexAfter(&bind.CallOpts{}, l2OutputBlock)
	if err != nil {
		return "", err
	}

	fmt.Println("debug22")
	header, err := l2.HeaderByNumber(ctx, l2OutputBlock)
	if err != nil {
		return "", err
	}
	fmt.Println("debug23")
	params, err := withdrawals.ProveWithdrawalParameters(ctx, l2g, l2, l2TxHash, header, &l2oo.L2OutputOracleCaller)
	if err != nil {
		return "", err
	}

	fmt.Println("debug24")
	type typesWithdrawalTransaction struct {
		Nonce    *big.Int       `json:"nonce"`
		Sender   common.Address `json:"sender"`
		Target   common.Address `json:"target"`
		Value    *big.Int       `json:"value"`
		GasLimit *big.Int       `json:"gasLimit"`
		Data     []byte         `json:"data"`
	}
	type outputRootProof struct {
		Version                  [32]byte `json:"version"`
		StateRoot                [32]byte `json:"stateRoot"`
		MessagePasserStorageRoot [32]byte `json:"MessagePasserStorageRoot"`
		LatestBlockhash          [32]byte `json:"latestBlockhash"`
	}

	var res = struct {
		TypesWithdrawalTransaction typesWithdrawalTransaction `json:"typesWithdrawalTransaction"`
		L2OutputIndex              *big.Int                   `json:"l2OutputIndex"`
		OutputRootProof            outputRootProof            `json:"outputRootProof"`
		WithdrawalProof            [][]byte                   `json:"withdrawalProof"`
	}{
		TypesWithdrawalTransaction: typesWithdrawalTransaction{
			Nonce:    params.Nonce,
			Sender:   params.Sender,
			Target:   params.Target,
			Value:    params.Value,
			GasLimit: params.GasLimit,
			Data:     params.Data,
		},
		L2OutputIndex: l2OutputIndex,
		OutputRootProof: outputRootProof{
			Version:                  params.OutputRootProof.Version,
			StateRoot:                params.OutputRootProof.StateRoot,
			MessagePasserStorageRoot: params.OutputRootProof.MessagePasserStorageRoot,
			LatestBlockhash:          params.OutputRootProof.LatestBlockhash,
		},
		WithdrawalProof: params.WithdrawalProof,
	}

	resBytes, err := json.Marshal(res)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(resBytes), nil
}

func CompleteWithdrawal2(ctx context.Context, l1 *ethclient.Client, l2c *rpc.Client, l2oo *bindings.L2OutputOracle, portal *bindings.OptimismPortal, l2TxHash common.Hash, finalizationPeriod *big.Int) (string, error) {
	l2 := ethclient.NewClient(l2c)
	l2g := gethclient.New(l2c)

	fmt.Println("debug30")
	receipt, err := l2.TransactionReceipt(ctx, l2TxHash)
	if err != nil {
		return "", err
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		return "", errors.New("unsuccessful withdrawal receipt status")
	}

	fmt.Println("debug31")
	l2WithdrawalBlock, err := l2.BlockByNumber(ctx, receipt.BlockNumber)
	if err != nil {
		return "", err
	}

	l2OutputBlockNr, err := l2oo.LatestBlockNumber(&bind.CallOpts{})
	if err != nil {
		return "", err
	}

	l2OutputBlock, err := l2.BlockByNumber(ctx, l2OutputBlockNr)
	if err != nil {
		return "", err
	}

	if l2OutputBlock.NumberU64() < l2WithdrawalBlock.NumberU64() {
		return "", errors.New(fmt.Sprintf("the latest L2 output is %d and is not past L2 block %d that includes the withdrawal yet, no withdrawal can be completed yet", l2OutputBlock.NumberU64(), l2WithdrawalBlock.NumberU64()))
	}

	l1Head, err := l1.HeaderByNumber(ctx, nil)
	if err != nil {
		return "", err
	}

	if l2WithdrawalBlock.Time()+finalizationPeriod.Uint64() >= l1Head.Time {
		return "", errors.New(fmt.Sprintf("withdrawal tx %s was included in L2 block %d (time %d) but L1 only knows of L2 proposal %d (time %d) at head %d (time %d) which has not reached output confirmation yet (period is %d)",
			l2TxHash, l2WithdrawalBlock.NumberU64(), l2WithdrawalBlock.Time(), l2OutputBlock.NumberU64(), l2OutputBlock.Time(), l1Head.Number.Uint64(), l1Head.Time, finalizationPeriod.Uint64()))
	}

	header, err := l2.HeaderByNumber(ctx, l2OutputBlockNr)
	if err != nil {
		return "", err
	}

	params, err := withdrawals.ProveWithdrawalParameters(ctx, l2g, l2, l2TxHash, header, &l2oo.L2OutputOracleCaller)
	if err != nil {
		return "", err
	}

	type typesWithdrawalTransaction struct {
		Nonce    *big.Int       `json:"nonce"`
		Sender   common.Address `json:"sender"`
		Target   common.Address `json:"target"`
		Value    *big.Int       `json:"value"`
		GasLimit *big.Int       `json:"gasLimit"`
		Data     []byte         `json:"data"`
	}

	var res = typesWithdrawalTransaction{
		Nonce:    params.Nonce,
		Sender:   params.Sender,
		Target:   params.Target,
		Value:    params.Value,
		GasLimit: params.GasLimit,
		Data:     params.Data,
	}

	resBytes, err := json.Marshal(res)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(resBytes), nil
}

func CompleteWithdrawal(ctx context.Context, l1 *ethclient.Client, l2c *rpc.Client, l2oo *bindings.L2OutputOracle, portal *bindings.OptimismPortal, l2TxHash common.Hash, finalizationPeriod *big.Int, opts *bind.TransactOpts) error {
	l2 := ethclient.NewClient(l2c)
	l2g := gethclient.New(l2c)

	receipt, err := l2.TransactionReceipt(ctx, l2TxHash)
	if err != nil {
		return err
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		return errors.New("unsuccessful withdrawal receipt status")
	}

	l2WithdrawalBlock, err := l2.BlockByNumber(ctx, receipt.BlockNumber)
	if err != nil {
		return err
	}

	l2OutputBlockNr, err := l2oo.LatestBlockNumber(&bind.CallOpts{})
	if err != nil {
		return err
	}

	l2OutputBlock, err := l2.BlockByNumber(ctx, l2OutputBlockNr)
	if err != nil {
		return err
	}

	if l2OutputBlock.NumberU64() < l2WithdrawalBlock.NumberU64() {
		fmt.Printf("the latest L2 output is %d and is not past L2 block %d that includes the withdrawal yet, no withdrawal can be completed yet", l2OutputBlock.NumberU64(), l2WithdrawalBlock.NumberU64())
		return nil
	}

	l1Head, err := l1.HeaderByNumber(ctx, nil)
	if err != nil {
		return err
	}

	if l2WithdrawalBlock.Time()+finalizationPeriod.Uint64() >= l1Head.Time {
		fmt.Printf("withdrawal tx %s was included in L2 block %d (time %d) but L1 only knows of L2 proposal %d (time %d) at head %d (time %d) which has not reached output confirmation yet (period is %d)",
			l2TxHash, l2WithdrawalBlock.NumberU64(), l2WithdrawalBlock.Time(), l2OutputBlock.NumberU64(), l2OutputBlock.Time(), l1Head.Number.Uint64(), l1Head.Time, finalizationPeriod.Uint64())
		return nil
	}

	header, err := l2.HeaderByNumber(ctx, l2OutputBlockNr)
	if err != nil {
		return err
	}

	params, err := withdrawals.ProveWithdrawalParameters(ctx, l2g, l2, l2TxHash, header, &l2oo.L2OutputOracleCaller)
	if err != nil {
		return err
	}

	tx, err := portal.FinalizeWithdrawalTransaction(
		opts,
		bindings.TypesWithdrawalTransaction{
			Nonce:    params.Nonce,
			Sender:   params.Sender,
			Target:   params.Target,
			Value:    params.Value,
			GasLimit: params.GasLimit,
			Data:     params.Data,
		},
	)
	if err != nil {
		return err
	}

	fmt.Printf("Completed withdrawal for %s: %s\n", l2TxHash.String(), tx.Hash().String())

	// Wait 5 mins max for confirmation
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	return waitForConfirmation(ctxWithTimeout, l1, tx.Hash())
}

func waitForConfirmation(ctx context.Context, client *ethclient.Client, tx common.Hash) error {
	for {
		receipt, err := client.TransactionReceipt(ctx, tx)
		if err == ethereum.NotFound {
			fmt.Printf("waiting for tx confirmation\n")
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(5 * time.Second):
			}
		} else if err != nil {
			return err
		} else if receipt.Status != types.ReceiptStatusSuccessful {
			return errors.New("unsuccessful withdrawal receipt status")
		} else {
			break
		}
	}
	fmt.Printf("%s confirmed\n", tx.String())
	return nil
}
