package commodity

import (
	"context"
	"github.com/bzzzm/commodity-brain/pkg/board"
	"github.com/bzzzm/commodity-brain/pkg/db"
	"github.com/bzzzm/commodity-brain/pkg/utils"
	"go.etcd.io/bbolt"
	"go.uber.org/zap"
	"os"
	"periph.io/x/periph/host"
	"runtime"
)

type Commodity struct {
	Board      *board.CommodityBoard
	Db         *bbolt.DB
	Config     *utils.Config
	ContextBag *utils.ContextBag

	//Camera     cc.OpenCVCapture
	//Stream     cc.MJPEGStream
	//HTTP       http.HTTPServer

	ctx context.Context
	ec  chan error
	//ic  chan cc.CaptureImage
}

func NewCommodity(ctx context.Context, config *utils.Config, ec chan error, qc chan os.Signal) *Commodity {

	// create the context bag
	cb := &utils.ContextBag{}

	// periph host
	if _, err := host.Init(); err != nil {
		zap.S().Panicf("cannot init host: %v", err)
	}

	// create db
	database, err := db.SetupDB(config.Database.Path)
	if err != nil {
		zap.S().Panicf("cannot initiate db: %v", err)
	}

	// create board object
	boardCtx, cancel := context.WithCancel(ctx)
	cb.Board = utils.ContextPair{
		Context: boardCtx,
		Cancel:  cancel,
	}

	var fake bool
	if runtime.GOARCH != "arm" {
		fake = true
	}

	cboard, err := board.NewBoard(boardCtx, &config.Board, database, fake, ec)
	if err != nil {
		zap.S().Panicf("cannot initiate board: %v", err)
	}

	return &Commodity{cboard, database, config, cb, ctx, ec}

}

func (c *Commodity) Close() {
	c.Board.Close()
}
func (c *Commodity) Start() {

	// start the board and keep the database in sync
	c.Board.Start()
}

//
//// NewCommodity creates a new Commodity instance
//func NewCommodity(ctx context.Context, config *utils.Config, ec chan error, qc chan os.Signal) *Commodity {
//
//	// disable debug messages fro i2c library
//	_ = logger.ChangePackageLogLevel("i2c", logger.InfoLevel)
//
//	// create the context bag
//	cb := &utils.ContextBag{}
//
//	// create i2c master
//	twi, err := i2c.NewI2C(config.AStar.Addr, config.AStar.BusID)
//	if err != nil {
//		if runtime.GOARCH == "arm" {
//			zap.S().Panicf("cannot initiate i2c: %v", err)
//		}
//
//	}
//
//	// create the database
//	db, err := cdb.SetupDB(config.Database.Path)
//	if err != nil {
//		zap.S().Panicf("cannot initiate db: %v", err)
//	}
//
//	// create astar object and start a goroutine to keep the DB in sync
//	astarCtx, cancel := context.WithCancel(ctx)
//	cb.Astar = utils.ContextPair{
//		Context: astarCtx,
//		Cancel:  cancel,
//	}
//	astar := ca.NewAStar(astarCtx, twi, db, bufferSize, ec)
//
//	// create the camera capture
//	captureCtx, cancel := context.WithCancel(ctx)
//	cb.Camera = utils.ContextPair{
//		Context: captureCtx,
//		Cancel:  cancel,
//	}
//	ic := make(chan cc.CaptureImage, 5)
//	camera, err := cc.NewCapture(captureCtx, config.Camera, ic)
//	if err != nil {
//		zap.S().Panicf("cannot start capture: %v", err)
//	}
//
//	// create mjpeg stream
//	streamCtx, cancel := context.WithCancel(ctx)
//	cb.Stream = utils.ContextPair{
//		Context: streamCtx,
//		Cancel:  cancel,
//	}
//	stream := cc.NewStream(streamCtx, ic)
//
//	// http server
//	httpCtx, cancel := context.WithCancel(ctx)
//	cb.HTTP = utils.ContextPair{
//		Context: httpCtx,
//		Cancel:  cancel,
//	}
//	httpd := http.NewHTTPServer(httpCtx, stream, config.HTTP)
//
//	return &Commodity{*astar, *camera, *stream, *httpd, db, config, cb, ctx, ec, ic}
//}
//
//// Start will start all peripherial jobs
//func (c *Commodity) Start() {
//
//	// Sync the database from i2C
//	if runtime.GOARCH == "arm" {
//		go c.Astar.SyncDB(time.Second / c.Config.AStar.Refresh)
//	}
//
//	// start batter logger
//	if runtime.GOARCH == "arm" {
//		go c.Astar.BatteryLogger()
//	}
//
//	// Start video capturing
//	go c.Camera.Start()
//
//	// Start MJPEG stream
//	go c.Stream.Start()
//
//	// Start HTTP Server
//	go c.HTTP.Start()
//}
//
//// Close will be used in defer to close underlying objects
//func (c *Commodity) Close() {
//	if runtime.GOARCH == "arm" {
//		_ = c.Astar.I2C.Close()
//	}
//
//	c.ContextBag.Astar.Cancel()
//	c.ContextBag.Camera.Cancel()
//	c.ContextBag.HTTP.Cancel()
//}
//
//// WatchButtons creates watchButton goroutines for each button.
//func (c *Commodity) WatchButtons() {
//	//fn1 := []func(interface{}, interface{}){c.Camera.Shutter}
//	//go c.Astar.WatchButton(1, fn1)
//
//	//fn2 := []func(){c.Astar.BatteryButton, c.Astar.DumpDB}
//	//go c.Astar.WatchButton(2, fn2)
//	//
//	//fn3 := []func(){c.Astar.TestLed}
//	//go c.Astar.WatchButton(3, fn3)
//}
