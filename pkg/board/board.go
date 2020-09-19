package board

import (
	"context"
	"go.etcd.io/bbolt"
	"go.uber.org/zap"
	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/conn/physic"
	"sync"
	"time"
)

const (
	bufferSize = 32
	I2CAddr    = 0x20
	I2CBus     = "1"
	I2CSpeed   = 400 * physic.KiloHertz
)

type CommodityBoard struct {
	db  *bbolt.DB
	mux *sync.Mutex

	Bus i2c.BusCloser
	Dev *i2c.Dev

	ctx context.Context
	ec  chan error
}

// NewI2CBoard creates a new Commodity Board instance that is able to communicate over I2C
func NewBoard(ctx context.Context, db *bbolt.DB, ec chan error) (*CommodityBoard, error) {

	// open I2C bus+device
	b, err := i2creg.Open(I2CBus)
	if err != nil {
		return nil, err
	}

	d := &i2c.Dev{Addr: uint16(I2CAddr), Bus: b}

	// t3h board
	board := &CommodityBoard{
		db:  db,
		ec:  ec,
		mux: &sync.Mutex{},
		ctx: ctx,

		Bus: b,
		Dev: d,
	}
	return board, nil
}

// Start is used to start ... @todo
func (b *CommodityBoard) Start() {

	go b.Heartbeat(time.Second)

	go b.SyncPower(time.Second)
	go b.SyncEnv(time.Second)

	//go b.SyncIMU(time.Second / 2)
	//go b.SyncTOF(time.Second / 2)
	//go b.SyncWifi(time.Second * 2)
	//go b.Watchdog(time.Second)
}

// Close is used to close ... @todo
func (b *CommodityBoard) Close() {
	b.Bus.Close()
}

// Heartbeat shifts a bit on the board and blinks the blue led
func (b *CommodityBoard) Heartbeat(t time.Duration) {
	tick := time.Tick(t)

	for {
		select {
		case <-tick:

			// switch blueled
			err := b.SwitchLed(BlueLed)
			if err != nil {
				b.ec <- err
				continue
			}

		case <-b.ctx.Done():
			zap.S().Debugf("Heartbeat stopped")
			return
		}
	}
}
