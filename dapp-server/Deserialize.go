//[128,9,2,1,2,0,16,95,83,69,76,76,95,95,95,79,78,71,95,79,78,84,95,0,20,164,200,17,159,236,16,104,78,212,62,178,208,75,222,28,36,63,63,113,56,2,4,128,240,250,2,2,1,1,2,1,2,2,0,2,0,2,0]

package main

import (
	"math/big"

	"github.com/ontio/ontology/common"
)

const (
	ARRAY = 128
	BYTES = 0
	INT   = 2
)

// 反转HexString
func reverseString2(s string) string {
	runes := []rune(s)
	for from, to := 0, len(runes)-1; from < to; from, to = from+2, to-2 {
		runes[from], runes[to-1] = runes[to-1], runes[from]
		runes[from+1], runes[to] = runes[to], runes[from+1]
	}
	return string(runes)
}

// 反转Bytearray
func reverseBytes(res []byte) []byte {
	for from, to := 0, len(res)-1; from < to; from, to = from+1, to-1 {
		res[from], res[to] = res[to], res[from]
	}
	return res
}

func Deserialize(source []byte) interface{} {
	switch source[0] {
	case ARRAY:
		temp := 2
		res := []interface{}{}
		for i := 0; i < int(source[1]); i++ {
			interface_tmp := Deserialize(source[temp : int(source[temp+1])+temp+2])
			temp = int(source[temp+1]) + temp + 2
			res = append(res, interface_tmp)
		}
		return res
	case BYTES:
		if source[1] == 20 {
			addr, err := common.AddressParseFromBytes(source[2:])
			if err != nil {

				return string(source[2:])
			}
			return addr.ToBase58()
		}
		return string(source[2:])
	case INT:
		iTmp := big.NewInt(0)
		if source[1] != 0 {
			iTmp.SetBytes(reverseBytes(source[2:]))
		}
		return iTmp
	default:
		return nil
	}
}
