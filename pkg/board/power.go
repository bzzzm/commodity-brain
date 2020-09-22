package board

import (
	"github.com/bzzzm/commodity-brain/pkg/db"
	"github.com/bzzzm/commodity-brain/pkg/utils"
	"go.etcd.io/bbolt"
	"go.uber.org/zap"
	"time"
)

const (
	PowerReg     = 0xa0
	PowerDataLen = 12

	VoltageKey = "power_volt"
	CurrentKey = "power_amp"
	PowerKey   = "power_watt"
)

func (b CommodityBoard) SyncPower() {

	var err error
	var voltage, current, power float32
	var start time.Time

	buf := make([]byte, PowerDataLen)
	tick := time.Tick(b.config.Update.Power)

	zap.S().Debugf("SyncPower started with update rate: %v", b.config.Update.Power)

	for {
		select {
		case <-tick:

			// read from i2c
			b.mux.Lock()

			start = time.Now()

			err = b.Dev.Tx([]byte{PowerReg}, buf)
			if err != nil {
				b.ec <- err
				continue
			}
			b.mux.Unlock()

			// convert to float32
			voltage = utils.Float32FromBytes(buf[0:4])
			current = utils.Float32FromBytes(buf[4:8])
			power = utils.Float32FromBytes(buf[8:])

			// update db
			err = b.db.Update(func(tx *bbolt.Tx) error {
				bucket := tx.Bucket([]byte(db.BoardBucket))
				err = bucket.Put([]byte(VoltageKey), buf[0:4])
				if err != nil {
					return err
				}

				err = bucket.Put([]byte(CurrentKey), buf[4:8])
				if err != nil {
					return err
				}

				err = bucket.Put([]byte(PowerKey), buf[8:])
				if err != nil {
					return err
				}
				return nil
			})

			if err != nil {
				b.ec <- err
				continue
			}

			zap.S().Debugf("SyncPower finished in %v - %v mV; %v mA; %v mW", time.Now().Sub(start), voltage, current, power)

		case <-b.ctx.Done():
			zap.S().Debugf("SyncPower stopped")
			return
		}
	}
}
