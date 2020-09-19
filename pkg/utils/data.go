package utils

import (
	"encoding/binary"
	"math"
)

// ByteToBool returns false in case b is 0x0 else returns true
func ByteToBool(b byte) bool {
	if b == 0x0 {
		return false
	}

	return true
}

func BoolToByte(b bool) byte {
	if b {
		return 0x1
	}
	return 0x0
}

func BoolByteToString(b byte) string {
	if bl := ByteToBool(b); bl {
		return "on"
	} else {
		return "off"
	}

}

func Float32FromBytes(bytes []byte) float32 {
	b := binary.LittleEndian.Uint32(bytes)
	f := math.Float32frombits(b)
	return f
}
