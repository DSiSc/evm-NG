package contract

import (
	"github.com/DSiSc/craft/types"
	"github.com/DSiSc/evm-NG/system"
	"github.com/DSiSc/evm-NG/system/contract/storage"
)

// system call routes
var routes map[types.Address]map[string]*system.SysCallFunc

// init system call routes
func init() {
	routes[system.TencentCosAddr] = storage.CosRoutes
}

// Call call system contract with specified args
func Call(addr types.Address, args []byte) (ret []byte, err error) {
	//TODO
	return nil, nil
}
