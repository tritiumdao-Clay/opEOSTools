package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/gin-contrib/cors"
	"log"
	"math/big"
	"os"
	"time"

	"crossChain/prove/signer"
	"crossChain/prove/withdraw"
	"github.com/ethereum-optimism/optimism/op-bindings/bindings"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/gin-gonic/gin"
)

type network struct {
	l2RPC         string
	portalAddress string
	l2OOAddress   string
}

var networks = map[string]network{
	"opeostest": {
		l2RPC:         "https://testnet-rpc.opeos.io",
		portalAddress: "0x70d544de5f1c7C4a9f09a82a07eB8F360B040169",
		l2OOAddress:   "0xB28aF3ac0c2847DE28345bb3b821dd8744f6fC6F",
	},
	"opeos": {
		l2RPC:         "",
		portalAddress: "",
		l2OOAddress:   "",
	},
}

var rpcFlag string
var networkFlag string
var withdrawalFlag string
var privateKey string
var startHTTP bool

type L2ToL1 struct {
	WithdrawHash []string `json:"l2ToL1WithdrawHash"`
	ProveHash    []string `json:"l2ToL1ProveHash"`
	FinalizeHash []string `json:"l2ToL1FinalizeHash"`
}

type WithdrawHashDatabaseItem struct {
	UserAddr   string   `json:"userAddr"`
	L1ToL2Hash []string `json:"l1ToL2Hash"`
	L2ToL1Hash L2ToL1   `json:"l2ToL1Hash"`
}

var database map[string]WithdrawHashDatabaseItem

func writeFile(path string) error {
	a, err := json.Marshal(database)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, a, 0644)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	var networkKeys []string
	for n := range networks {
		networkKeys = append(networkKeys, n)
	}

	flag.StringVar(&rpcFlag, "rpc", "", "L1 RPC")
	flag.StringVar(&networkFlag, "network", "opeostest", "network name")
	flag.StringVar(&withdrawalFlag, "withdrawal", "", "L2 withdraw txHash")
	flag.StringVar(&privateKey, "private-key", "", "private key")
	flag.BoolVar(&startHTTP, "start-http", true, "whether to start http server")
	flag.Parse()

	log.Default().SetFlags(0)

	n, ok := networks[networkFlag]
	if !ok {
		log.Fatalf("unknown network: %s", networkFlag)
	}
	if rpcFlag == "" {
		log.Fatalf("missing --rpc flag")
	}

	if startHTTP {
		pwd, _ := os.Getwd()
		path := pwd + "/datadir/withrawHash.json"
		fmt.Println("path:", path)

		database = make(map[string]WithdrawHashDatabaseItem, 0)

		withdrawHashData, err := os.ReadFile(path)
		if err != nil {
			fmt.Println("load withdraw hash database fail")
			panic(err)
		}

		err = json.Unmarshal(withdrawHashData, &database)
		if err != nil {
			fmt.Println("parse database json fail")
			panic(err)
		}
		fmt.Println("----------")
		fmt.Println(database)
		fmt.Println("----------")
		err = writeFile(path + time.Now().String())
		if err != nil {
			fmt.Println("save database fail")
			panic(err)
		}

		r := gin.Default()
		config := cors.DefaultConfig()
		config.AllowOrigins = []string{"*"}                                        // 允许什么域名访问，支持多个域名
		config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}  // 允许的 HTTP 方法
		config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type"} // 允许的 HTTP 头
		config.AllowCredentials = true
		// 设置cors中间件
		r.Use(cors.New(config))

		r.POST("/getL1ToL2Hash", getL1ToL2Hash)
		r.POST("/getL2ToL1Hash", getL2ToL1Hash)
		r.POST("/writel1ToL2Hash", writeL1ToL2Hash)
		r.POST("/writeWithdrawHash", writeWithdrawHash)
		r.POST("/writeProveHash", writeProveHash)
		r.POST("/writeFinalizeHash", writeFinalizeHash)

		r.POST("/getProveWithdrawalPara", getProveWithdrawalPara)
		r.POST("/getFinalizePara", getFinalizePara)

		r.Run(":10003")
	} else {
		if withdrawalFlag == "" {
			log.Fatalf("missing --withdrawal flag")
		}
		withdrawal := common.HexToHash(withdrawalFlag)

		ctx := context.Background()
		l1Client, err := ethclient.DialContext(ctx, rpcFlag)
		if err != nil {
			log.Fatalf("Error dialing L1 client: %v", err)
		}
		l2Client, err := rpc.DialContext(ctx, n.l2RPC)
		if err != nil {
			log.Fatalf("Error dialing L2 client: %v", err)
		}

		portal, err := bindings.NewOptimismPortal(common.HexToAddress(n.portalAddress), l1Client)
		if err != nil {
			log.Fatalf("Error binding OptimismPortal contract: %v", err)
		}

		l2oo, err := bindings.NewL2OutputOracle(common.HexToAddress(n.l2OOAddress), l1Client)
		if err != nil {
			log.Fatalf("Error binding L2OutputOracle contract: %v", err)
		}

		isFinalized, err := withdraw.ProofFinalized(ctx, portal, withdrawal)
		if err != nil {
			log.Fatalf("Error querying withdrawal finalization status: %v", err)
		}
		if isFinalized {
			fmt.Println("Withdrawal already finalized")
			return
		}
		finalizationPeriod, err := l2oo.FINALIZATIONPERIODSECONDS(&bind.CallOpts{})
		if err != nil {
			log.Fatalf("Error querying withdrawal finalization period: %v", err)
		}
		submissionInterval, err := l2oo.SUBMISSIONINTERVAL(&bind.CallOpts{})
		if err != nil {
			log.Fatalf("Error querying output proposal submission interval: %v", err)
		}
		l2BlockTime, err := l2oo.L2BLOCKTIME(&bind.CallOpts{})
		if err != nil {
			log.Fatalf("Error querying output proposal L2 block time: %v", err)
		}
		l2OutputBlock, err := l2oo.LatestBlockNumber(&bind.CallOpts{})
		if err != nil {
			log.Fatalf("Error querying latest proposed block: %v", err)
		}
		l2WithdrawalBlock, err := withdraw.TxBlock(ctx, l2Client, withdrawal)
		if err != nil {
			log.Fatalf("Error querying withdrawal tx block: %v", err)
		}
		proof, err := withdraw.ProvenWithdrawal(ctx, l2Client, portal, withdrawal)
		if err != nil {
			log.Fatalf("Error querying withdrawal proof: %v", err)
		}
		if l2OutputBlock.Uint64() < l2WithdrawalBlock.Uint64() {
			log.Fatalf("The latest L2 output is %d and is not past L2 block %d that includes the withdrawal, no withdrawal can be proved yet.\nPlease wait for the next proposal submission to %s, which happens every %v.",
				l2OutputBlock.Uint64(), l2WithdrawalBlock.Uint64(), n.l2OOAddress, time.Duration(submissionInterval.Int64()*l2BlockTime.Int64())*time.Second)
		}

		s, err := signer.CreateSigner(privateKey)
		if err != nil {
			log.Fatalf("Error creating signer: %v", err)
		}

		l1ChainID, err := l1Client.ChainID(ctx)
		if err != nil {
			log.Fatalf("Error querying chain ID: %v", err)
		}

		l1Nonce, err := l1Client.PendingNonceAt(ctx, s.Address())
		if err != nil {
			log.Fatalf("Error querying nonce: %v", err)
		}

		l1opts := &bind.TransactOpts{
			From:    s.Address(),
			Signer:  s.SignerFn(l1ChainID),
			Context: ctx,
			Nonce:   big.NewInt(int64(l1Nonce) - 1), // subtract 1 because we add 1 each time newl1opts is called
		}
		newl1opts := func() *bind.TransactOpts {
			l1opts.Nonce = big.NewInt(0).Add(l1opts.Nonce, big.NewInt(1))
			return l1opts
		}
		if proof.Timestamp.Uint64() == 0 {
			err = withdraw.ProveWithdrawal(ctx, l1Client, l2Client, l2oo, portal, withdrawal, newl1opts())
			if err != nil {
				log.Fatalf("Error proving withdrawal: %v", err)
			}
			fmt.Printf("The withdrawal can be completed after the finalization period, in approximately %v\n", time.Duration(finalizationPeriod.Int64())*time.Second)
			return
		}

		err = withdraw.CompleteWithdrawal(ctx, l1Client, l2Client, l2oo, portal, withdrawal, finalizationPeriod, newl1opts())
		if err != nil {
			log.Fatalf("Error completing withdrawal: %v", err)
		}
	}
}

func initWork(withdrawalFlag string) (l1 *ethclient.Client, l2c *rpc.Client, l2oo *bindings.L2OutputOracle, portal *bindings.OptimismPortal, l2TxHash common.Hash, finalizationPeriod *big.Int, err error) {
	if len(withdrawalFlag) == 0 || withdrawalFlag[:2] != "0x" || len(withdrawalFlag) != 66 {
		return nil, nil, nil, nil, common.Hash{}, nil, errors.New("withdrawHash is invalid")
	}
	withdrawal := common.HexToHash(withdrawalFlag)
	l2TxHash = withdrawal

	ctx := context.Background()
	l1Client, err := ethclient.DialContext(ctx, rpcFlag)
	if err != nil {
		return nil, nil, nil, nil, common.Hash{}, nil, errors.New(fmt.Sprintf("Error dialing L1 client: %v", err))
	}
	n := networks[networkFlag]
	l2Client, err := rpc.DialContext(ctx, n.l2RPC)
	if err != nil {
		return nil, nil, nil, nil, common.Hash{}, nil, errors.New(fmt.Sprintf("Error dialing L1 client: %v", err))
	}

	for {
		l2 := ethclient.NewClient(l2Client)
		fmt.Println("debug00", l2TxHash)
		if l2 == nil {
			fmt.Println("debug00, l2 == nil")
		}
		receipt, err := l2.TransactionReceipt(ctx, l2TxHash)
		if err != nil {
			fmt.Println("debug00", err.Error())
			break
		}
		fmt.Println("debug01")
		fmt.Println("debug01", receipt.BlockHash)
		break
	}

	portal, err = bindings.NewOptimismPortal(common.HexToAddress(n.portalAddress), l1Client)
	if err != nil {
		return nil, nil, nil, nil, common.Hash{}, nil, errors.New(fmt.Sprintf("Error binding OptimismPortal contract: %v", err))
	}

	l2oo, err = bindings.NewL2OutputOracle(common.HexToAddress(n.l2OOAddress), l1Client)
	if err != nil {
		return nil, nil, nil, nil, common.Hash{}, nil, errors.New(fmt.Sprintf("Error binding L2OutputOracle contract: %v", err))
	}

	isFinalized, err := withdraw.ProofFinalized(ctx, portal, withdrawal)
	if err != nil {
		return nil, nil, nil, nil, common.Hash{}, nil, errors.New(fmt.Sprintf("Error querying withdrawal finalization status: %v", err))
	}
	if isFinalized {
		return nil, nil, nil, nil, common.Hash{}, nil, errors.New(fmt.Sprintf("Withdrawal already finalized"))
	}
	finalizationPeriod, err = l2oo.FINALIZATIONPERIODSECONDS(&bind.CallOpts{})
	if err != nil {
		log.Fatalf("Error querying withdrawal finalization period: %v", err)
	}

	submissionInterval, err := l2oo.SUBMISSIONINTERVAL(&bind.CallOpts{})
	if err != nil {
		return nil, nil, nil, nil, common.Hash{}, nil, errors.New(fmt.Sprintf("Error querying output proposal submission interval: %v", err))
	}
	l2BlockTime, err := l2oo.L2BLOCKTIME(&bind.CallOpts{})
	if err != nil {
		return nil, nil, nil, nil, common.Hash{}, nil, errors.New(fmt.Sprintf("Error querying output proposal L2 block time: %v", err))
	}

	l2OutputBlock, err := l2oo.LatestBlockNumber(&bind.CallOpts{})
	if err != nil {
		return nil, nil, nil, nil, common.Hash{}, nil, errors.New(fmt.Sprintf("Error querying latest proposed block: %v", err))
	}
	l2WithdrawalBlock, err := withdraw.TxBlock(ctx, l2Client, withdrawal)

	if err != nil {
		return nil, nil, nil, nil, common.Hash{}, nil, errors.New(fmt.Sprintf("Error querying withdrawal tx block: %v", err))
	}

	if l2OutputBlock.Uint64() < l2WithdrawalBlock.Uint64() {
		return nil, nil, nil, nil, common.Hash{}, nil, errors.New(fmt.Sprintf("The latest L2 output is %d and is not past L2 block %d that includes the withdrawal, no withdrawal can be proved yet.\nPlease wait for the next proposal submission to %s, which happens every %v.",
			l2OutputBlock.Uint64(), l2WithdrawalBlock.Uint64(), n.l2OOAddress, time.Duration(submissionInterval.Int64()*l2BlockTime.Int64())*time.Second))
	}
	return l1Client, l2Client, l2oo, portal, l2TxHash, finalizationPeriod, nil
}

type Error struct {
	Err string `json:"error"`
}
type Success struct {
	Suc string `json:"result"`
}

func wrapError(error string) string {
	var a = Error{
		Err: error,
	}
	ret, _ := json.Marshal(a)
	return string(ret)
}

func wrapSuccess(success string) string {
	var a = Success{
		Suc: success,
	}
	ret, _ := json.Marshal(a)
	return string(ret)
}

func writeFinalizeHash(c *gin.Context) {
	pwd, _ := os.Getwd()
	path := pwd + "/datadir/withrawHash.json"
	fmt.Println("path:", path)
	type Resp struct {
		UserAddr     string `json:"userAddr"`
		WithdrawHash string `json:"txHash"`
	}
	var res = Resp{}
	c.Header("Content-Type", "application/json")

	dataB, err := c.GetRawData()
	err = json.Unmarshal(dataB, &res)
	if err != nil {
		c.String(200, wrapError("parse json fail. req:"+string(dataB)))
		return
	}

	tmp := database[res.UserAddr].L2ToL1Hash.FinalizeHash
	for _, item := range tmp {
		if item == res.WithdrawHash {
			c.String(200, wrapError("already contain this withdraw hash"))
			return
		}
	}
	tmp = append(tmp, res.WithdrawHash)
	database[res.UserAddr] = WithdrawHashDatabaseItem{
		UserAddr:   res.UserAddr,
		L1ToL2Hash: database[res.UserAddr].L1ToL2Hash,
		L2ToL1Hash: L2ToL1{
			WithdrawHash: database[res.UserAddr].L2ToL1Hash.WithdrawHash,
			ProveHash:    database[res.UserAddr].L2ToL1Hash.ProveHash,
			FinalizeHash: tmp,
		},
	}
	err = writeFile(path)
	if err != nil {
		c.String(200, wrapError("write fail"))
		return
	}
	c.String(200, wrapSuccess("success"))
	return
}

func writeProveHash(c *gin.Context) {
	pwd, _ := os.Getwd()
	path := pwd + "/datadir/withrawHash.json"
	fmt.Println("path:", path)
	type Resp struct {
		UserAddr     string `json:"userAddr"`
		WithdrawHash string `json:"txHash"`
	}
	var res = Resp{}
	c.Header("Content-Type", "application/json")

	dataB, err := c.GetRawData()
	err = json.Unmarshal(dataB, &res)
	if err != nil {
		c.String(200, wrapError("parse json fail. req:"+string(dataB)))
		return
	}

	tmp := database[res.UserAddr].L2ToL1Hash.ProveHash
	for _, item := range tmp {
		if item == res.WithdrawHash {
			c.String(200, wrapError("already contain this withdraw hash"))
			return
		}
	}
	tmp = append(tmp, res.WithdrawHash)
	database[res.UserAddr] = WithdrawHashDatabaseItem{
		UserAddr:   res.UserAddr,
		L1ToL2Hash: database[res.UserAddr].L1ToL2Hash,
		L2ToL1Hash: L2ToL1{
			WithdrawHash: database[res.UserAddr].L2ToL1Hash.WithdrawHash,
			ProveHash:    tmp,
			FinalizeHash: database[res.UserAddr].L2ToL1Hash.FinalizeHash,
		},
	}
	err = writeFile(path)
	if err != nil {
		c.String(200, wrapError("write fail"))
		return
	}
	c.String(200, wrapSuccess("success"))
	return
}

func writeL1ToL2Hash(c *gin.Context) {
	pwd, _ := os.Getwd()
	path := pwd + "/datadir/withrawHash.json"
	fmt.Println("path:", path)
	type Resp struct {
		UserAddr     string `json:"userAddr"`
		WithdrawHash string `json:"txHash"`
	}
	var res = Resp{}
	c.Header("Content-Type", "application/json")

	dataB, err := c.GetRawData()
	err = json.Unmarshal(dataB, &res)
	if err != nil {
		c.String(200, wrapError("parse json fail. req:"+string(dataB)))
		return
	}

	tmp := database[res.UserAddr].L1ToL2Hash
	for _, item := range tmp {
		if item == res.WithdrawHash {
			c.String(200, wrapError("already contain this withdraw hash"))
			return
		}
	}
	tmp = append(tmp, res.WithdrawHash)
	database[res.UserAddr] = WithdrawHashDatabaseItem{
		UserAddr:   res.UserAddr,
		L1ToL2Hash: tmp,
		L2ToL1Hash: L2ToL1{
			WithdrawHash: database[res.UserAddr].L2ToL1Hash.WithdrawHash,
			ProveHash:    database[res.UserAddr].L2ToL1Hash.ProveHash,
			FinalizeHash: database[res.UserAddr].L2ToL1Hash.FinalizeHash,
		},
	}
	err = writeFile(path)
	if err != nil {
		c.String(200, wrapError("write fail"))
		return
	}
	c.String(200, wrapSuccess("success"))
	return
}
func writeWithdrawHash(c *gin.Context) {
	pwd, _ := os.Getwd()
	path := pwd + "/datadir/withrawHash.json"
	fmt.Println("path:", path)
	type Resp struct {
		UserAddr     string `json:"userAddr"`
		WithdrawHash string `json:"txHash"`
	}
	var res = Resp{}
	c.Header("Content-Type", "application/json")

	dataB, err := c.GetRawData()
	err = json.Unmarshal(dataB, &res)
	if err != nil {
		c.String(200, wrapError("parse json fail. req:"+string(dataB)))
		return
	}

	tmp := database[res.UserAddr].L2ToL1Hash.WithdrawHash
	for _, item := range tmp {
		if item == res.WithdrawHash {
			c.String(200, wrapError("already contain this withdraw hash"))
			return
		}
	}
	tmp = append(tmp, res.WithdrawHash)
	database[res.UserAddr] = WithdrawHashDatabaseItem{
		UserAddr:   res.UserAddr,
		L1ToL2Hash: database[res.UserAddr].L1ToL2Hash,
		L2ToL1Hash: L2ToL1{
			WithdrawHash: tmp,
			ProveHash:    database[res.UserAddr].L2ToL1Hash.ProveHash,
			FinalizeHash: database[res.UserAddr].L2ToL1Hash.FinalizeHash,
		},
	}
	err = writeFile(path)
	if err != nil {
		c.String(200, wrapError("write fail"))
		return
	}
	c.String(200, wrapSuccess("success"))
	return
}

func getL1ToL2Hash(c *gin.Context) {
	type Req struct {
		UserAddr string
	}
	var req = Req{}
	c.Header("Content-Type", "application/json")

	dataA, err := c.GetRawData()
	err = json.Unmarshal(dataA, &req)
	if err != nil {
		c.String(200, wrapError(`{"error":"parse json fail"}`))
		return
	}
	dataB := req.UserAddr
	var userAddr string
	userAddr = dataB
	if len(userAddr) != 42 {
		c.String(200, wrapError("len(str) is must be 42"))
		return
	}
	withdrashHashes, is := database[userAddr]
	if !is {
		c.String(200, wrapError("no withdraw hash"))
		return
	}
	withdrashHashesBytes, err := json.Marshal(struct {
		L1ToL2 []string `json:"l1ToL2Hash"`
	}{
		L1ToL2: withdrashHashes.L1ToL2Hash,
	})
	if err != nil {
		c.String(200, wrapError("internal fail:"+err.Error()))
		return
	}
	c.Data(200, "application/json", withdrashHashesBytes)
	return
}

func getL2ToL1Hash(c *gin.Context) {
	type Req struct {
		UserAddr string
	}
	var req = Req{}
	c.Header("Content-Type", "application/json")

	dataA, err := c.GetRawData()
	err = json.Unmarshal(dataA, &req)
	if err != nil {
		c.String(200, wrapError(`{"error":"parse json fail"}`))
		return
	}
	dataB := req.UserAddr
	var userAddr string
	userAddr = dataB
	if len(userAddr) != 42 {
		c.String(200, wrapError("len(str) is must be 42"))
		return
	}
	withdrashHashes, is := database[userAddr]
	if !is {
		c.String(200, wrapError("no withdraw hash"))
		return
	}
	withdrashHashesBytes, err := json.Marshal(struct {
		L2ToL1 L2ToL1 `json:"l2ToL1Hash"`
	}{
		L2ToL1: withdrashHashes.L2ToL1Hash,
	})
	if err != nil {
		c.String(200, wrapError("internal fail:"+err.Error()))
		return
	}
	c.Data(200, "application/json", withdrashHashesBytes)
	return
}

func getProveWithdrawalPara(c *gin.Context) {

	type Req struct {
		WithdrawHash string
	}
	var req = Req{}
	c.Header("Content-Type", "application/json")

	dataA, err := c.GetRawData()
	err = json.Unmarshal(dataA, &req)
	if err != nil {
		c.String(200, wrapError(`{"error":"parse json fail"}`))
		return
	}
	dataB := req.WithdrawHash
	//fmt.Println("debug0:", dataB, len(dataB))

	l1, l2c, l2oo, portal, l2TxHash, _, err := initWork(dataB)
	if err != nil {
		c.String(200, wrapError(err.Error()))
		return
	}
	ret, err := withdraw.ProveWithdrawal2(context.Background(), l1, l2c, l2oo, portal, l2TxHash)
	if err != nil {
		c.String(200, wrapError(err.Error()))
		return
	}
	c.String(200, ret)
	return
}

func getFinalizePara(c *gin.Context) {
	type Req struct {
		WithdrawHash string
	}
	var req = Req{}
	c.Header("Content-Type", "application/json")
	//fmt.Println("debug0:", string(dataA), len(dataA))

	dataA, err := c.GetRawData()
	err = json.Unmarshal(dataA, &req)
	if err != nil {
		c.String(200, wrapError(`{"error":"parse json fail"}`))
		return
	}
	dataB := req.WithdrawHash

	l1, l2c, l2oo, portal, l2TxHash, finalizationPeriod, err := initWork(dataB)
	if err != nil {
		c.String(200, wrapError(err.Error()))
		return
	}
	ret, err := withdraw.CompleteWithdrawal2(context.Background(), l1, l2c, l2oo, portal, l2TxHash, finalizationPeriod)
	if err != nil {
		c.String(200, wrapError(err.Error()))
		return
	}
	c.String(200, ret)
}
