package Interaction

import (
	"fmt"
	atypes "github.com/DSiSc/apigateway/core/types"
	"github.com/DSiSc/craft/monitor"
	"github.com/DSiSc/craft/rlp"
	"github.com/DSiSc/craft/types"
	"github.com/DSiSc/crypto-suite/crypto"
	cutil "github.com/DSiSc/crypto-suite/util"
	"github.com/DSiSc/evm-NG/system/contract/util"
	"github.com/DSiSc/repository"
	sutil "github.com/DSiSc/statedb-NG/util"
	"github.com/DSiSc/txpool"
	wtypes "github.com/DSiSc/wallet/core/types"
	wutils "github.com/DSiSc/wallet/utils"
	"math/big"
	"github.com/DSiSc/craft/log"
)

const(
	FAILED = iota
	SUCCESS
	PENDING
)

type Status uint64

var CrossChainAddr = cutil.HexToAddress("0000000000000000000000000000000000011100")
var (
	forwardFundsMethodHash = string(util.ExtractMethodHash(util.Hash([]byte("forwardFunds(string, uint64, string)"))))
	getTxStateMethodHash = string(util.ExtractMethodHash(util.Hash([]byte("getTxState(string,string)"))))
)

type CrossChainContract struct {
	records map[types.Address]Status
}

func NewCrossChainContract() *CrossChainContract {
	return new(CrossChainContract)
}

//如何获得合约的调用者？？？，保证资金安全性
func (this *CrossChainContract) forwardFunds(toAddr types.Address, amount uint64, chainFlag string) (types.Hash, bool) {
	//调用apigateway的receiveCrossTx交易
	from, err := getPubliceAcccount()
	if err != nil {
		return types.Hash{}, false
	}

	hash, err := CallCrossRawTransactionReq(from, toAddr, amount, chainFlag)
	if err != nil {
		return types.Hash{}, false
	}
	return hash, false
}

func (this *CrossChainContract) getTxState(address types.Address, chainFlag string) (uint64, bool) {
	switch chainFlag {
		case "chainA":
			//123
			fmt.Println()
		default:
	}

	return SUCCESS, true
}

func CallCrossRawTransactionReq(from types.Address, to types.Address, amount uint64, chainFlag string) (types.Hash, error) {
	monitor.JTMetrics.ApigatewayReceivedTx.Add(1)

	tx := new(types.Transaction)

	// Patchwork tx，fix from -- get publicAccount
	addr, err := getPubliceAcccount()
	if err != nil {
		return types.Hash{}, err
	}
	tx.Data.From = &addr

	// like sendTransaction, need sig
	private := "29ad43a4ebb4a65436d9fb116d471d96516b3d5cc153e045b384664bed5371b9"

	//get nonce
	bc, _ := repository.NewLatestStateRepository()
	noncePool := txpool.GetPoolNonce(*tx.Data.From)
	nonceChain := bc.GetNonce(*tx.Data.From)
	nonce := uint64(0)
	if noncePool > nonceChain {
		nonce = noncePool + 1
	} else {
		nonce = nonceChain
	}
	tx.Data.AccountNonce = nonce
	tx.Data.Price = big.NewInt(100)
	tx.Data.GasLimit = 6721975
	tx.Data.Recipient = &to
	tx.Data.Amount = big.NewInt(int64(amount))

	//payload
	tx.Data.Payload = []byte(sutil.AddressToHex(from))

	//sign tx
	priKey, err := crypto.HexToECDSA(private)
	if err != nil {
		return types.Hash{}, err
	}

	//chainId, err := config.GetChainIdFromConfig()
	//if err != nil {
	//	return types.Hash{}, err
	//}

	chainID := big.NewInt(int64(5777))
	tx, err = wtypes.SignTx(tx, wtypes.NewEIP155Signer(chainID), priKey)
	if err != nil {
		return types.Hash{}, err
	}

	from123, _ := wtypes.Sender(wtypes.NewEIP155Signer(big.NewInt(5777)), tx)

	log.Info("from %x", addr)
	log.Info("from_123 %x", from123)
	log.Info("tx recipient: %x", tx.Data.Recipient)
	log.Info("tx to: %x", to)

	//call the broadcast the tx
	var port string
	switch chainFlag {
	case "chainA":
		port = "47768"
		break
	case "chainB":
		port = "47769"
		break
	default:
		port = ""
	}

	//web, err := wutils.NewWeb3("47.92.1.248", "47768", false)
	web, err := wutils.NewWeb3("127.0.0.1", port, false)
	if err != nil {
		return types.Hash{}, err
	}

	txBytes, _ := rlp.EncodeToBytes(tx)
	txHash, err := web.Eth.SendRawTransaction(txBytes)
	if err != nil {
		return types.Hash{}, err
	}

	return (types.Hash)(txHash), nil
}

func getPubliceAcccount() (types.Address, error){
	//get from config or genesis ?
	addr := "0x0fA3E9c7065Cf9b5f513Fb878284f902d167870c"
	address := atypes.HexToAddress(addr)

	return types.Address(address), nil
}

