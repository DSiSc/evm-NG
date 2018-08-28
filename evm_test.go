// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
package evm

import (
	"testing"

	"encoding/hex"
	"fmt"
	"github.com/DSiSc/blockchain"
	"github.com/DSiSc/blockchain/config"
	"github.com/DSiSc/craft/types"
	"github.com/DSiSc/statedb-NG/common/crypto"
	"math/big"
)

// GenesisAlloc specifies the initial state that is part of the genesis block.
type GenesisAlloc map[types.Address]GenesisAccount

// GenesisAccount is an account in the state of the genesis block.
type GenesisAccount struct {
	Code       []byte                    `json:"code,omitempty"`
	Storage    map[types.Hash]types.Hash `json:"storage,omitempty"`
	Balance    *big.Int                  `json:"balance" gencodec:"required"`
	Nonce      uint64                    `json:"nonce,omitempty"`
	PrivateKey []byte                    `json:"secretKey,omitempty"` // for tests
}

var callerAddress = types.HexToAddress("0x8a8c58e424f4a6d2f0b2270860c96dfe34f10c78")
var contractAddress = types.HexToAddress("0xf74cc8824a00bcb96e8546bf3b4dc47ace9cab2c")

//var code,_ = hex.DecodeString("0x6000600001600055");
var code, _ = hex.DecodeString("6080604052348015600f57600080fd5b5060998061001e6000396000f300608060405260043610603e5763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416634f2be91f81146043575b600080fd5b348015604e57600080fd5b5060556067565b60408051918252519081900360200190f35b610378905600a165627a7a723058205d540f3e87376532c076a230eb73eee4aa46c0df1a71cdba5a33cda64a8e6f400029")
var input1, _ = hex.DecodeString("4f2be91f")
var input2, _ = hex.DecodeString("4f2be91f")

func TestVM(t *testing.T) {
	//init statedb state
	testStateDb := MakePreState(GenesisAlloc{})
	insertData(testStateDb)
	fmt.Println(testStateDb.GetCode(contractAddress))

	//init chain config
	callerRef := AccountRef(callerAddress)

	//execute contract code
	evmInst := newEVM(testStateDb)
	resp, leftgas, error := evmInst.Call(callerRef, contractAddress, input1, 3000, big.NewInt(0))
	fmt.Println("Resp:", resp, " Left Gas:", leftgas, " Error:", error)
	resp, leftgas, error = evmInst.Call(callerRef, contractAddress, input2, 3000, big.NewInt(0))
	fmt.Println("Resp:", resp, " Left Gas:", leftgas, " Error:", error)
}

func MakePreState(accounts GenesisAlloc) *blockchain.BlockChain {
	blockchain.InitBlockChain(config.BlockChainConfig{
		PluginName:    blockchain.PLUGIN_MEMDB,
		StateDataPath: "",
		BlockDataPath: "",
	})
	state, _ := blockchain.NewLatestStateBlockChain()
	for addr, a := range accounts {
		state.SetCode(addr, a.Code)
		state.SetNonce(addr, a.Nonce)
		state.SetBalance(addr, a.Balance)
		for k, v := range a.Storage {
			state.SetState(addr, k, v)
		}
	}
	// Commit and re-open to start with a clean statedb.
	root, _ := state.Commit(false)
	state, _ = blockchain.NewBlockChainByHash(root)
	return state
}

func insertData(db *blockchain.BlockChain) {
	//create caller account
	db.CreateAccount(callerAddress)
	db.AddBalance(callerAddress, big.NewInt(1000))

	//create contract account
	db.CreateAccount(contractAddress)
	db.SetCode(contractAddress, code)

	db.Commit(false)
}

func newEVM(statedb1 *blockchain.BlockChain) *EVM {
	canTransfer := func(db *blockchain.BlockChain, address types.Address, amount *big.Int) bool {
		return true
	}
	transfer := func(db *blockchain.BlockChain, sender, recipient types.Address, amount *big.Int) {}
	context := Context{
		CanTransfer: canTransfer,
		Transfer:    transfer,
		GetHash:     vmTestBlockHash,
		Origin:      callerAddress,
		Coinbase:    callerAddress,
		BlockNumber: big.NewInt(0x00),
		Time:        new(big.Int).SetUint64(0x01),
		GasLimit:    0x0f4240,
		Difficulty:  big.NewInt(0x0100),
		GasPrice:    big.NewInt(0x5af3107a4000),
	}
	return NewEVM(context, statedb1)
}

func vmTestBlockHash(n uint64) types.Hash {
	return types.BytesToHash(crypto.Keccak256([]byte(big.NewInt(int64(n)).String())))
}
