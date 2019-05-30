package evm

import (
	"github.com/DSiSc/evm-NG/common/hexutil"
	"github.com/DSiSc/evm-NG/common/math"
	"github.com/DSiSc/evm-NG/constant"
	"github.com/DSiSc/evm-NG/system/contract/util"
	"github.com/justitia/common"
	"io"
	"math/big"
)

const (
	writeMethodName = "write(bytes)"
	readMethodName  = "read(uint256,uint256)"
	closeMethodName = "close()"
)

// SolidityBuffer solidity contract 'SolidityBuffer' call tool
type SolidityBuffer struct {
	evm            *EVM        // evm instance
	callerContract ContractRef // callerContract
	cursor         int         // cursor of the buffer
}

// create a new instance
func NewSolidityBuffer(evm *EVM, callerContract ContractRef) *SolidityBuffer {
	return &SolidityBuffer{
		evm:            evm,
		callerContract: callerContract,
	}
}

// Read read the data from the buffer
func (this *SolidityBuffer) Read(p []byte) (n int, err error) {
	ret, err := this.readFromSolidityBuffer(big.NewInt(int64(this.cursor)), big.NewInt(int64(len(p))))
	if err != nil {
		return 0, err
	}
	if len(ret) == 0 {
		return 0, io.EOF
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
	methodParam := util.Hash([]byte(closeMethodName))[:4]
	_, _, err := this.evm.Call(this.callerContract, this.callerContract.Address(), methodParam, math.MaxUint64, big.NewInt(0))
	return err
}

// read data from solidity buffer
func (this *SolidityBuffer) readFromSolidityBuffer(pos, size *big.Int) ([]byte, error) {
	methodParam := util.Hash([]byte(readMethodName))[:4]
	posParam := math.PaddedBigBytes(pos, constant.EvmWordSize)
	sizeParam := math.PaddedBigBytes(size, constant.EvmWordSize)
	input := append(methodParam, posParam...)
	input = append(input, sizeParam...)
	ret, _, err := this.evm.Call(this.callerContract, this.callerContract.Address(), input, math.MaxUint64, big.NewInt(0))
	if err != nil {
		return nil, err
	}
	dataLen, _ := math.ParseUint64(hexutil.Encode(ret[constant.EvmWordSize : 2*constant.EvmWordSize]))
	return ret[2*constant.EvmWordSize : 2*constant.EvmWordSize+dataLen], nil
}

// write data to solidity contract buffer
func (this *SolidityBuffer) writeToSolidityBuffer(data []byte) error {
	method := util.Hash([]byte(writeMethodName))[:4]
	dataLen := math.PaddedBigBytes(big.NewInt(int64(len(data))), constant.EvmWordSize)
	offset := math.PaddedBigBytes(big.NewInt(int64(len(dataLen))), constant.EvmWordSize)
	encodedData := common.RightPadBytes(data, constant.EvmWordSize)

	input := append(method, offset...)
	input = append(input, dataLen...)
	input = append(input, encodedData...)
	_, _, err := this.evm.Call(this.callerContract, this.callerContract.Address(), input, math.MaxUint64, big.NewInt(0))
	return err
}
