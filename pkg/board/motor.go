package board

import (
	"go.uber.org/zap"
	"time"
)

const (
	MotorReg     = 0xb2
	MotorDataLen = 4
	MotorLF      = 0
	MotorRF      = 1
	MotorLB      = 2
	MotorRB      = 3
)

func (b *CommodityBoard) SetMotors(lf int8, rf int8, lb int8, rb int8) error {

	var err error
	var start time.Time

	// lock the mutex and write the new values to i2c
	// old values are not relevant
	// @todo: maybe read stop flag from TOF
	b.mux.Lock()

	start = time.Now()
	buf := []byte{byte(lf), byte(rf), byte(lb), byte(rb)}

	// write the data
	_, err = b.Dev.Write(append([]byte{MotorReg}, buf...))
	if err != nil {
		return err
	}

	b.mux.Unlock()

	// stop the timer and log.debug the new values
	zap.S().Debugf("SetMotors finished in %v - lf %v, rf %v, lb %v, rb %v", time.Now().Sub(start), lf, rf, lb, rb)

	return nil
}
