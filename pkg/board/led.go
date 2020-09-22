package board

import (
	"go.uber.org/zap"
	"time"
)

const (
	LedReg     = 0xb2
	LedDataLen = 3
	RedLed     = 0
	YellowLed  = 1
	BlueLed    = 2
)

func (b *CommodityBoard) SwitchLed(led int) error {
	var err error
	var start time.Time
	buf := make([]byte, LedDataLen)

	// read from i2c
	b.mux.Lock()

	start = time.Now()
	err = b.Dev.Tx([]byte{LedReg}, buf)
	if err != nil {
		return err
	}

	// error correction
	for i, s := range buf {
		if s > 1 {
			buf[i] = 1
		}
	}

	// switch the led in buffer
	if buf[led] == 0 {
		buf[led] = 1
	} else {
		buf[led] = 0
	}

	// write back the data
	_, err = b.Dev.Write(append([]byte{LedReg}, buf...))
	if err != nil {
		return err
	}

	b.mux.Unlock()

	zap.S().Debugf("SwitchLed finished in %v - Led %v: %v", time.Now().Sub(start), led, int(buf[led]))

	return nil
}
