package util

import (
	"github.com/DSiSc/evm-NG"
	"github.com/DSiSc/evm-NG/common"
	"github.com/DSiSc/evm-NG/common/hexutil"
	"github.com/DSiSc/evm-NG/common/math"
	"github.com/pkg/errors"
	"golang.org/x/crypto/sha3"
	"math/big"
	"reflect"
)

// ExtractMethodHash extract method hash from input
func ExtractMethodHash(input []byte) []byte {
	return input[:4]
}

// Hash return the hash of the data
func Hash(data []byte) []byte {
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(data)
	return hasher.Sum(nil)
}

// ExtractParam extract string params from input
func ExtractParam(input []byte, argTypes ...reflect.Kind) ([]interface{}, error) {
	args := make([]interface{}, 0)
	for i := 0; i < len(argTypes); i++ {
		switch argTypes[i] {
		case reflect.String:
			offset, _ := math.ParseUint64(hexutil.Encode(input[i*evm.EvmWordSize : (i+1)*evm.EvmWordSize]))
			dataLen, _ := math.ParseUint64(hexutil.Encode(input[offset : offset+evm.EvmWordSize]))
			argStart := offset + evm.EvmWordSize
			argEnd := argStart + dataLen
			arg := string(input[argStart:argEnd])
			args = append(args, arg)
		default:
			return nil, errors.New("unsupported arg type")
		}
	}
	return args, nil
}

// EncodeReturnValue encode the return value to the format needed by evm
func EncodeReturnValue(retVals ...interface{}) ([]byte, error) {
	retPre := make([]byte, 0)
	retData := make([]byte, 0)
	preOffsetPadding := len(retVals) * evm.EvmWordSize
	for _, retVal := range retVals {
		switch reflect.TypeOf(retVal).Kind() {
		case reflect.String:
			offset := preOffsetPadding + len(retData)
			retPre = append(retPre, math.PaddedBigBytes(big.NewInt(int64(offset)), evm.EvmWordSize)...)
			retData = append(retData, encodeString(retVal.(string))...)
		default:
			return nil, errors.New("unsupported return type")
		}
	}
	return append(retPre, retData...), nil
}

// encode the string to the format needed by evm
func encodeString(val string) []byte {
	ret := make([]byte, 0)
	valB := []byte(val)
	ret = append(ret, math.PaddedBigBytes(big.NewInt(int64(len(valB))), evm.EvmWordSize)...)
	for i := 0; i < len(valB); {
		if (len(valB) - i) > evm.EvmWordSize {
			ret = append(ret, valB[i:i+evm.EvmWordSize]...)
			i += evm.EvmWordSize
		} else {
			ret = append(ret, common.RightPadBytes(valB[i:], evm.EvmWordSize)...)
			i += len(valB)
		}
	}
	return ret
}
