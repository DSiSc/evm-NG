package evm

import (
	"github.com/DSiSc/craft/types"
	"github.com/DSiSc/crypto-suite/util"
	"github.com/DSiSc/evm-NG/system/contract/storage"
)

//system contract address
var (
	TencentCosAddr = util.HexToAddress("0000000000000000000000000000000000011111")
)

// SysContractExecutionFunc system contract execute function
type SysContractExecutionFunc func(interpreter *EVM, contract ContractRef, input []byte) ([]byte, error)

// system call routes
var routes = make(map[types.Address]SysContractExecutionFunc)

func init() {
	routes[TencentCosAddr] = func(execEvm *EVM, caller ContractRef, input []byte) ([]byte, error) {
		solidityBuffer := NewSolidityBuffer(execEvm, caller)
		tencentCos := storage.NewTencentCosContract(solidityBuffer)
		return storage.CosExecute(tencentCos, input)
	}
}

//IsSystemContract check the contract with specified address is system contract
func IsSystemContract(addr types.Address) bool {
	return routes[addr] != nil
}

// GetSystemContractExecFunc get system contract execution function by address
func GetSystemContractExecFunc(addr types.Address) SysContractExecutionFunc {
	return routes[addr]
}
