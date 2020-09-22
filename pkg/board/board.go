package board

import (
	"context"
	"github.com/bzzzm/commodity-brain/pkg/utils"
	"go.etcd.io/bbolt"
	"go.uber.org/zap"
	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"
	"sync"
	"time"
)

type CommodityBoard struct {
	Bus i2c.BusCloser
	Dev *i2c.Dev

	ctx    context.Context
	ec     chan error
	db     *bbolt.DB
	mux    *sync.Mutex
	config *utils.BoardConfig
	fake   bool
}

// NewI2CBoard creates a new Commodity Board instance that is able to communicate over I2C
func NewBoard(ctx context.Context, config *utils.BoardConfig, db *bbolt.DB, fake bool, ec chan error) (*CommodityBoard, error) {

	// create the standard board
	board := &CommodityBoard{
		db:     db,
		ec:     ec,
		mux:    &sync.Mutex{},
		ctx:    ctx,
		config: config,
		fake:   fake,
	}

	// if running on a device connect to the board, add the I2C Bus and Device
	// to the board object
	if !fake {
		bus, err := i2creg.Open(string(config.BusID))
		if err != nil {
			return nil, err
		}

		dev := &i2c.Dev{Addr: uint16(config.Addr), Bus: bus}

		// t3h board
		board.Dev = dev
		board.Bus = bus
	}

	return board, nil
}

// Start is used to start ... @todo
func (b *CommodityBoard) Start() {

	// if the board is not fake, start the board processes
	// otherwise start some dummy routines to fill random data in db
	if !b.fake {
		go b.Heartbeat()

		go b.SyncPower()
		go b.SyncEnv()
	} else {

	}
	//go b.SyncIMU(time.Second / 2)
	//go b.SyncTOF(time.Second / 2)
	//go b.SyncWifi(time.Second * 2)
	//go b.Watchdog(time.Second)
}

// Close is used to close ... @todo
func (b *CommodityBoard) Close() {
	if !b.fake {
		_ = b.Bus.Close()
	}
}

// IsFake returns b.fake attr
func (b *CommodityBoard) IsFake() bool {
	return b.fake
}

// Heartbeat shifts a bit on the board and blinks the blue led
func (b *CommodityBoard) Heartbeat() {
	tick := time.Tick(b.config.Update.Heartbeat)

	zap.S().Debugf("Heartbeat started with update rate: %v", b.config.Update.Heartbeat)

	for {
		select {
		case <-tick:

			// switch blue led
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
