package rpc

import (
	"errors"
	"fmt"
	cutil "github.com/DSiSc/crypto-suite/util"
	"github.com/DSiSc/evm-NG/system/contract/Interaction"
	"github.com/DSiSc/evm-NG/system/contract/util"
	wutils "github.com/DSiSc/wallet/utils"
	wcmn "github.com/DSiSc/web3go/common"
	"reflect"
)

var RpcContractAddr = cutil.HexToAddress("0000000000000000000000000000000000011101")

// rpc routes
var routes = map[string]*RPCFunc{
	string(util.ExtractMethodHash(util.Hash([]byte("ForwardFunds(string,uint64,string)")))): NewRPCFunc(ForwardFunds),
	string(util.ExtractMethodHash(util.Hash([]byte("GetTxState(string,string)")))): NewRPCFunc(GetTxState),
}

// 0 means failed, 1 means success
func ForwardFunds(toAddr string, amount uint64, chainFlag string) (error, string, uint64) {
	//from := RpcContractAddr
	from, _ := Interaction.GetPubliceAcccount()
	to := cutil.HexToAddress(toAddr)
	hash, err := Interaction.CallCrossRawTransactionReq(from, to, amount, chainFlag)
	if err != nil {
		return err, "", 0
	}

	hashBytes := cutil.HashToBytes(hash)
	return err,  wcmn.BytesToHex(hashBytes), 1
}

// GetCross Tx state
func GetTxState(txHash string, chainFlag string) (error, uint64){

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

	web, err := wutils.NewWeb3("127.0.0.1", port, false)
	if err != nil {
		return err, 0
	}

	hash := cutil.HexToHash(txHash)
	receipt, err := web.Eth.GetTransactionReceipt(wcmn.Hash(hash))
	if err != nil || receipt == nil {
		return err, 0
	}

	status := uint64(receipt.Status.Int64())

	return err, status
}

// Register register a rpc route
func Register(methodName string, f *RPCFunc) error {
	paramStr := ""
	for _, arg := range f.args {
		switch arg.Kind() {
		case reflect.Uint64:
			paramStr += "uint64,"
		case reflect.String:
			paramStr += "string,"
		case reflect.Slice:
			if reflect.Uint8 != arg.Elem().Kind() {
				return errors.New("unsupported arg type")
			}
		}
	}
	if len(paramStr) > 0 {
		paramStr = paramStr[:len(paramStr)-1]
	}
	methodHash := util.Hash([]byte(methodName + "(" + paramStr + ")"))[:4]
	routes[string(methodHash)] = f
	return nil
}

func Handler(input []byte) ([]byte, error) {
	method := util.ExtractMethodHash(input)
	rpcFunc := routes[string(method)]
	if rpcFunc == nil {
		return nil, errors.New("routes not found")
	}
	args, err := inputParamsToArgs(rpcFunc, input[len(method):])
	if err != nil {
		return nil, err
	}
	returns := rpcFunc.f.Call(args)
	return encodeResult(returns)
}

// Covert an http query to a list of properly typed values.
// To be properly decoded the arg must be a concrete type from tendermint (if its an interface).
func inputParamsToArgs(rpcFunc *RPCFunc, input []byte) ([]reflect.Value, error) {
	args := make([]interface{}, 0)
	for _, argT := range rpcFunc.args {
		args = append(args, reflect.New(argT).Interface())
	}
	err := util.ExtractParam(input, args...)
	if err != nil {
		return nil, err
	}

	argVs := make([]reflect.Value, 0)
	for _, arg := range args {
		argVs = append(argVs, reflect.ValueOf(arg).Elem())
	}
	return argVs, nil
}

// NOTE: assume returns is result struct and error. If error is not nil, return it
func encodeResult(returns []reflect.Value) ([]byte, error) {
	errV := returns[0]
	if errV.Interface() != nil {
		return nil, errors.New(fmt.Sprintf("%v", errV.Interface()))
	}
	returns = returns[1:]
	rvs := make([]interface{}, 0)
	for _, rv := range returns {
		// the result is a registered interface,
		// we need a pointer to it so we can marshal with type byte
		rvp := reflect.New(rv.Type())
		rvp.Elem().Set(rv)
		rvs = append(rvs, rvp.Elem().Interface())
	}
	return util.EncodeReturnValue(rvs...)
}

// RPCFunc contains the introspected type information for a function
type RPCFunc struct {
	f       reflect.Value  // underlying rpc function
	args    []reflect.Type // type of each function arg
	returns []reflect.Type // type of each return arg
}

// NewRPCFunc create a new RPCFunc instance
func NewRPCFunc(f interface{}) *RPCFunc {
	return &RPCFunc{
		f:       reflect.ValueOf(f),
		args:    funcArgTypes(f),
		returns: funcReturnTypes(f),
	}
}

// return a function's argument types
func funcArgTypes(f interface{}) []reflect.Type {
	t := reflect.TypeOf(f)
	n := t.NumIn()
	typez := make([]reflect.Type, n)
	for i := 0; i < n; i++ {
		typez[i] = t.In(i)
	}
	return typez
}

// return a function's return types
func funcReturnTypes(f interface{}) []reflect.Type {
	t := reflect.TypeOf(f)
	n := t.NumOut()
	typez := make([]reflect.Type, n)
	for i := 0; i < n; i++ {
		typez[i] = t.Out(i)
	}
	return typez
}
