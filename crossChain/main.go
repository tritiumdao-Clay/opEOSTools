package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/gin-contrib/cors"
	"io"
	"log"
	"math/big"
	"net/http"
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
		l2RPC:         "http://13.228.210.115:8545",
		portalAddress: "0xdd52D429c7c85d2122EbEB3C5808fbf73caBe927",
		l2OOAddress:   "0xfAEFE87de2A01F26583B3922cfdea6fE2f285641",
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

type WithdrawHashDatabaseItem struct {
	UserAddr     string   `json:"userAddr"`
	WithdrawHash []string `json:"withdrawHash"`
}

type WithdrawHashDatabase struct {
	Database []WithdrawHashDatabaseItem `json:"database"`
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

		r.POST("/getProveWithdrawalPara", getUserTxHash)

		r.Run(":10003")

		//http.HandleFunc("/getProveWithdrawalPara", getProveWithdrawalPara)
		//http.HandleFunc("/getFinalizePara", getFinalizePara)
		//http.HandleFunc("/writeTxHash", writeTxHash)

		//log.Println("Go!")
		//http.ListenAndServe(":10003", nil)
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

func writeTxHash(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "X-PINGOTHER, Content-Type")
	w.Header().Set("Access-Control-Max-Age", "86400")
	switch r.Method {
	case "POST":
		pwd, _ := os.Getwd()
		path := pwd + "/datadir/withrawHash.json"
		fmt.Println("path:", path)

		dataA := make([]byte, 512)
		n, _ := r.Body.Read(dataA)
		dataB := string(dataA[:n])
		type Resp struct {
			UserAddr     string `json:"userAddr"`
			WithdrawHash string `json:"withdrawHash"`
		}
		var res = Resp{}
		err := json.Unmarshal([]byte(dataB), &res)
		if err != nil {
			io.WriteString(w, wrapError("parse json fail. req:"+dataB))
			return
		}

		tmp := database[res.UserAddr].WithdrawHash
		for _, item := range tmp {
			if item == res.WithdrawHash {
				io.WriteString(w, wrapError("already contain this withdraw hash"))
				return
			}
		}
		tmp = append(tmp, res.WithdrawHash)
		database[res.UserAddr] = WithdrawHashDatabaseItem{
			UserAddr:     res.UserAddr,
			WithdrawHash: tmp,
		}
		err = writeFile(path)
		if err != nil {
			io.WriteString(w, wrapError("write fail"))
			return
		}
		io.WriteString(w, wrapSuccess("success"))
		return

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		io.WriteString(w, wrapError("I can't do that"))
	}
}

func getUserTxHash(c *gin.Context) {
	type Req struct {
		UserAddr string
	}
	var req = Req{}

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
	withdrashHashesBytes, err := json.Marshal(withdrashHashes)
	if err != nil {
		c.String(200, wrapError("internal fail:"+err.Error()))
		return
	}
	c.Data(200, "application/json", withdrashHashesBytes)
	return
}

func getProveWithdrawalPara(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "X-PINGOTHER, Content-Type")
	w.Header().Set("Access-Control-Max-Age", "86400")
	switch r.Method {
	case "POST":
		dataA := make([]byte, 512)
		n, _ := r.Body.Read(dataA)
		type Req struct {
			WithdrawHash string
		}
		var req = Req{}
		err := json.Unmarshal(dataA[:n], &req)
		if err != nil {
			io.WriteString(w, wrapError(`{"error":"parse json fail"}`))
			return
		}
		dataB := req.WithdrawHash
		//fmt.Println("debug0:", dataB, len(dataB))

		l1, l2c, l2oo, portal, l2TxHash, _, err := initWork(dataB)
		if err != nil {
			io.WriteString(w, wrapError(err.Error()))
			return
		}
		ret, err := withdraw.ProveWithdrawal2(context.Background(), l1, l2c, l2oo, portal, l2TxHash)
		if err != nil {
			fmt.Println("debug11")
			io.WriteString(w, wrapError(err.Error()))
			return
		}
		io.WriteString(w, ret)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "I can't do that.")
	}
}

func getFinalizePara(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "X-PINGOTHER, Content-Type")
	w.Header().Set("Access-Control-Max-Age", "86400")
	switch r.Method {
	case "POST":
		dataA := make([]byte, 512)
		n, _ := r.Body.Read(dataA)
		type Req struct {
			WithdrawHash string
		}
		var req = Req{}
		//fmt.Println("debug0:", string(dataA), len(dataA))
		err := json.Unmarshal(dataA[:n], &req)
		if err != nil {
			io.WriteString(w, wrapError(`{"error":"parse json fail"}`))
			return
		}
		dataB := req.WithdrawHash

		l1, l2c, l2oo, portal, l2TxHash, finalizationPeriod, err := initWork(dataB)
		if err != nil {
			io.WriteString(w, wrapError(err.Error()))
			return
		}
		fmt.Println("debug1")
		ret, err := withdraw.CompleteWithdrawal2(context.Background(), l1, l2c, l2oo, portal, l2TxHash, finalizationPeriod)
		if err != nil {
			fmt.Println("debug11")
			io.WriteString(w, wrapError(err.Error()))
			return
		}
		fmt.Println("debug2")
		io.WriteString(w, ret)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "I can't do that.")
	}
}
