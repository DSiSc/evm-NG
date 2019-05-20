package evm

import (
	"github.com/DSiSc/craft/types"
	"github.com/DSiSc/crypto-suite/util"
)

const (
	EvmWordSize = 32
)

//system contract address
var (
	SolidityBufferAddr = util.HexToAddress("0000000000000000000000000000000000100000")
	TencentCosAddr     = util.HexToAddress("0000000000000000000000000000000000011111")
)

// SysContractExecutionFunc system contract execute function
type SysContractExecutionFunc func(interpreter *EVM, contract ContractRef, input []byte) ([]byte, error)

// system call routes
var routes = make(map[types.Address]SysContractExecutionFunc)

// RegisterRoutes register new contract to global routes
func RegisterRoutes(addr types.Address, execFunc SysContractExecutionFunc) {
	routes[addr] = execFunc
}

//IsSystemContract check the contract with specified address is system contract
func IsSystemContract(addr types.Address) bool {
	return routes[addr] != nil
}

// GetSystemContractExecFunc get system contract execution function by address
func GetSystemContractExecFunc(addr types.Address) SysContractExecutionFunc {
	return routes[addr]
}
