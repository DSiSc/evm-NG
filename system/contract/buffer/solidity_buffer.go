package buffer

import (
	"github.com/DSiSc/craft/types"
	"github.com/DSiSc/evm-NG"
	"github.com/DSiSc/evm-NG/common/hexutil"
	"github.com/DSiSc/evm-NG/common/math"
	"github.com/DSiSc/evm-NG/system"
	"github.com/justitia/common"
	"golang.org/x/crypto/sha3"
	"math/big"
)

const (
	writeMethodName = "write(bytes)"
	readMethodName  = "read(uint256,uint256)"
	closeMethodName = "close()"
	wordSize        = 32
)

// SolidityBuffer solidity contract 'SolidityBuffer' call tool
type SolidityBuffer struct {
	evm    *evm.EVM        // evm instance
	caller evm.ContractRef // caller
	addr   types.Address   // relative solidity contract address
	cursor int             // cursor of the buffer
}

// Read read the data from the buffer
func (this *SolidityBuffer) Read(p []byte) (n int, err error) {
	ret, err := this.readFromSolidityBuffer(big.NewInt(int64(this.cursor)), big.NewInt(int64(len(p))))
	if err != nil {
		return 0, err
	}
	this.cursor += len(ret)
	copy(p, ret)
	return len(ret), err
}

// Write write the data to buffer
func (this *SolidityBuffer) Write(p []byte) (n int, err error) {
	err = this.writeToSolidityBuffer(p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

// Close close the buffer
func (this *SolidityBuffer) Close() error {
	methodParam := hash([]byte(closeMethodName))[:4]
	_, _, err := this.evm.Call(this.caller, system.SystemBufferAddr, methodParam, math.MaxUint64, big.NewInt(0))
	return err
}

// read data from solidity buffer
func (this *SolidityBuffer) readFromSolidityBuffer(pos, size *big.Int) ([]byte, error) {
	methodParam := hash([]byte(readMethodName))[:4]
	posParam := math.PaddedBigBytes(pos, wordSize)
	sizeParam := math.PaddedBigBytes(size, wordSize)
	input := append(methodParam, posParam...)
	input = append(input, sizeParam...)
	ret, _, err := this.evm.Call(this.caller, system.SystemBufferAddr, input, math.MaxUint64, big.NewInt(0))
	if err != nil {
		return nil, err
	}
	dataLen, _ := math.ParseUint64(hexutil.Encode(ret[wordSize : 2*wordSize]))
	return ret[2*wordSize : 2*wordSize+dataLen], nil
}

// write data to solidity contract buffer
func (this *SolidityBuffer) writeToSolidityBuffer(data []byte) error {
	method := hash([]byte(writeMethodName))[:4]
	dataLen := math.PaddedBigBytes(big.NewInt(int64(len(data))), wordSize)
	offset := math.PaddedBigBytes(big.NewInt(int64(len(dataLen))), wordSize)
	encodedData := common.RightPadBytes(data, wordSize)

	input := append(method, offset...)
	input = append(input, dataLen...)
	input = append(input, encodedData...)
	_, _, err := this.evm.Call(this.caller, system.SystemBufferAddr, input, math.MaxUint64, big.NewInt(0))
	return err
}

// return the hash of the data
func hash(data []byte) []byte {
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(data)
	return hasher.Sum(nil)
}
