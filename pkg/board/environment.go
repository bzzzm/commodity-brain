package board

import (
	"github.com/bzzzm/commodity-brain/pkg/db"
	"github.com/bzzzm/commodity-brain/pkg/utils"
	"go.etcd.io/bbolt"
	"go.uber.org/zap"
	"time"
)

const (
	EnvReg     = 0xa1
	EnvDataLen = 12

	TemperatureKey = "env_temp"
	PressureKey    = "env_press"
)

func (b CommodityBoard) SyncEnv() {

	var err error
	var temperature, pressure, humidity float32
	var start time.Time

	buf := make([]byte, EnvDataLen)
	tick := time.Tick(b.config.Update.Env)

	zap.S().Debugf("SyncEnv started with update rate: %v", b.config.Update.Env)

	for {
		select {
		case <-tick:

			// read from i2c using the mutex lock and
			// start the timer only after the lock was aquired
			b.mux.Lock()

			start = time.Now()
			err = b.Dev.Tx([]byte{EnvReg}, buf)
			if err != nil {
				b.ec <- err
				continue
			}

			b.mux.Unlock()

			// convert to floats
			temperature = utils.Float32FromBytes(buf[0:4])
			pressure = utils.Float32FromBytes(buf[4:8])
			humidity = utils.Float32FromBytes(buf[8:])

			// update db (ignore humidity, chip is not capable to measure it)
			err = b.db.Update(func(tx *bbolt.Tx) error {
				bucket := tx.Bucket([]byte(db.BoardBucket))
				err = bucket.Put([]byte(TemperatureKey), buf[0:4])
				if err != nil {
					return err
				}

				err = bucket.Put([]byte(PressureKey), buf[4:8])
				if err != nil {
					return err
				}
				return nil
			})

			if err != nil {
				b.ec <- err
				continue
			}

			// log.debug the result
			zap.S().Debugf("SyncEnv finished in %v - %v C; %v ??; %v %%", time.Now().Sub(start), temperature, pressure, humidity)

		case <-b.ctx.Done():
			zap.S().Debugf("SyncEnv stopped")
			return
		}
	}
}
