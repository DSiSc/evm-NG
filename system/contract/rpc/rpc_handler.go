package rpc

import (
	"errors"
	"fmt"
	"github.com/DSiSc/evm-NG/system/contract/util"
	"reflect"
)

// rpc routes
var routes = map[string]*RPCFunc{}

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
