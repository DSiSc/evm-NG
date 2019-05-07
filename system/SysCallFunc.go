package system

import (
	"github.com/DSiSc/evm-NG/util"
	"reflect"
	"strings"
)

//system contract address
var (
	TencentCosAddr = util.HexToAddress("0000000000000000000000000000000000010000")
)

// SysCallFunc contains the introspected type information for a function
type SysCallFunc struct {
	f        reflect.Value  // underlying rpc function
	args     []reflect.Type // type of each function arg
	returns  []reflect.Type // type of each return arg
	argNames []string       // name of each argument
}

// NewSysCallFunc wraps a function for introspection.
// f is the function, args are comma separated argument names
func NewSysCallFunc(f interface{}, args string) *SysCallFunc {
	var argNames []string
	if args != "" {
		argNames = strings.Split(args, ",")
	}
	return &SysCallFunc{
		f:        reflect.ValueOf(f),
		args:     funcArgTypes(f),
		returns:  funcReturnTypes(f),
		argNames: argNames,
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
