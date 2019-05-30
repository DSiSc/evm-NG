package util

import (
	"github.com/DSiSc/crypto-suite/util"
	"github.com/DSiSc/evm-NG/common/hexutil"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHash(t *testing.T) {
	assert := assert.New(t)
	data := []byte("Hello, World")
	exepectedHash := []byte{0xa0, 0x4a, 0x45, 0x10, 0x28, 0xd0, 0xf9, 0x28, 0x4c, 0xe8, 0x22, 0x43, 0x75, 0x5e, 0x24, 0x52, 0x38, 0xab, 0x1e, 0x4e, 0xcf, 0x7b, 0x9d, 0xd8, 0xbf, 0x47, 0x34, 0xd9, 0xec, 0xfd, 0x5, 0x29}
	actualHash := Hash(data)
	assert.Equal(exepectedHash, actualHash)
}

func TestExtractMethodHash(t *testing.T) {
	assert := assert.New(t)
	expectedHash := Hash([]byte("hello(string,string)"))[:4]
	input, _ := hexutil.Decode("0x939531c0000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000001610000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000016200000000000000000000000000000000000000000000000000000000000000")
	methodHash := ExtractMethodHash(input)
	assert.Equal(expectedHash, methodHash)
}

func TestExtractParam(t *testing.T) {
	assert := assert.New(t)
	arg1 := new(string)
	arg2 := new(string)
	expectedParam1 := "a"
	expectedParam2 := "b"
	input, _ := hexutil.Decode("0x939531c0000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000001610000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000016200000000000000000000000000000000000000000000000000000000000000")
	err := ExtractParam(input[4:], arg1, arg2)
	assert.Nil(err)
	assert.Equal(expectedParam1, *arg1)
	assert.Equal(expectedParam2, *arg2)
}

func TestExtractParam2(t *testing.T) {
	assert := assert.New(t)
	arg1 := make([]byte, 0)
	arg2 := make([]byte, 0)
	expectedParam1 := []byte("a")
	expectedParam2 := []byte("b")
	input, _ := hexutil.Decode("0x939531c0000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000001610000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000016200000000000000000000000000000000000000000000000000000000000000")
	err := ExtractParam(input[4:], &arg1, &arg2)
	assert.Nil(err)
	assert.Equal(expectedParam1, arg1)
	assert.Equal(expectedParam2, arg2)
}

func TestExtractParam3(t *testing.T) {
	assert := assert.New(t)
	arg1 := uint64(2)
	expectedParam1 := uint64(2)
	input, _ := hexutil.Decode("0xe05e91e00000000000000000000000000000000000000000000000000000000000000002")
	err := ExtractParam(input[4:], &arg1)
	assert.Nil(err)
	assert.Equal(expectedParam1, arg1)
}

func TestEncodeReturnValue(t *testing.T) {
	assert := assert.New(t)
	retVal1 := "a"
	retVal2 := "b"
	expect, _ := hexutil.Decode("0x000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000001610000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000016200000000000000000000000000000000000000000000000000000000000000")
	retB, err := EncodeReturnValue(retVal1, retVal2)
	assert.Nil(err)
	assert.Equal(expect, retB)
}

func TestEncodeReturnValue2(t *testing.T) {
	assert := assert.New(t)
	retVal1 := []byte("a")
	retVal2 := []byte("b")
	expect, _ := hexutil.Decode("0x000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000001610000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000016200000000000000000000000000000000000000000000000000000000000000")
	retB, err := EncodeReturnValue(retVal1, retVal2)
	assert.Nil(err)
	assert.Equal(expect, retB)
}

func TestEncodeReturnValue3(t *testing.T) {
	assert := assert.New(t)
	retVal1 := uint64(2)
	expect, _ := hexutil.Decode("0x0000000000000000000000000000000000000000000000000000000000000002")
	retB, err := EncodeReturnValue(retVal1)
	assert.Nil(err)
	assert.Equal(expect, retB)
}

func TestEncodeReturnValue4(t *testing.T) {
	assert := assert.New(t)
	addr := util.HexToAddress("0000000000000000000000000000000000011110")
	expect, _ := hexutil.Decode("0x0000000000000000000000000000000000000000000000000000000000011110")
	retB, err := EncodeReturnValue(addr)
	assert.Nil(err)
	assert.Equal(expect, retB)
}
