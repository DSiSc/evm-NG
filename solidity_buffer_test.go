package evm

import (
	"fmt"
	"github.com/DSiSc/craft/types"
	"github.com/DSiSc/evm-NG/common/hexutil"
	"github.com/DSiSc/evm-NG/util"
	"github.com/DSiSc/monkey"
	"github.com/stretchr/testify/assert"
	"math/big"
	"reflect"
	"testing"
)

func mockSolidityBuffer() *SolidityBuffer {
	return &SolidityBuffer{
		evm:            &EVM{},
		callerContract: AccountRef(util.HexToAddress("0x8a8c58e424f4a6d2f0b2270860c96dfe34f10c78")),
	}
}

func TestSolidityBuffer_Read(t *testing.T) {
	defer monkey.UnpatchAll()
	assert := assert.New(t)
	data := make([]byte, 32)
	expectInputStr := "7508099700000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000020"
	acctualInputStr := ""
	expectRet, _ := hexutil.Decode("0x111111")
	solidityBuffer := mockSolidityBuffer()
	monkey.PatchInstanceMethod(reflect.TypeOf(solidityBuffer.evm), "Call", func(vm *EVM, caller ContractRef, addr types.Address, input []byte, gas uint64, value *big.Int) (ret []byte, leftOverGas uint64, err error) {
		acctualInputStr = fmt.Sprintf("%x", input)
		ret, _ = hexutil.Decode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000031111110000000000000000000000000000000000000000000000000000000000")
		return ret, 0, nil
	})
	n, err := solidityBuffer.Read(data)
	assert.Equal(expectInputStr, acctualInputStr)
	assert.Equal(expectRet, data[:n])
	assert.Nil(err)
}

func TestSolidityBuffer_Write(t *testing.T) {
	defer monkey.UnpatchAll()
	assert := assert.New(t)
	expectInputStr := "7ed0c3b2000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000011100000000000000000000000000000000000000000000000000000000000000"
	acctualInputStr := ""
	data, _ := hexutil.Decode("0x11")
	solidityBuffer := mockSolidityBuffer()
	monkey.PatchInstanceMethod(reflect.TypeOf(solidityBuffer.evm), "Call", func(vm *EVM, caller ContractRef, addr types.Address, input []byte, gas uint64, value *big.Int) (ret []byte, leftOverGas uint64, err error) {
		acctualInputStr = fmt.Sprintf("%x", input)
		return nil, 0, nil
	})
	_, err := solidityBuffer.Write(data)
	assert.Equal(expectInputStr, acctualInputStr)
	assert.Nil(err)
}

func TestSolidityBuffer_Close(t *testing.T) {
	defer monkey.UnpatchAll()
	assert := assert.New(t)
	expectInputStr := "43d726d6"
	acctualInputStr := ""
	solidityBuffer := mockSolidityBuffer()
	monkey.PatchInstanceMethod(reflect.TypeOf(solidityBuffer.evm), "Call", func(vm *EVM, caller ContractRef, addr types.Address, input []byte, gas uint64, value *big.Int) (ret []byte, leftOverGas uint64, err error) {
		acctualInputStr = fmt.Sprintf("%x", input)
		fmt.Println(acctualInputStr)
		return nil, 0, nil
	})
	err := solidityBuffer.Close()
	assert.Equal(expectInputStr, acctualInputStr)
	assert.Nil(err)
}
