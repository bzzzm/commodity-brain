package utils

import "time"

type Config struct {
	Database DatabaseConfig
	HTTP     HttpConfig
	Board    BoardConfig
	Camera   CameraConfig
}

type DatabaseConfig struct {
	Path string `desc:"Path to BoltDB file"`
}

type HttpConfig struct {
	Address string `desc:"HTTP address (127.0.0.1:8080)"`
	SSL     bool
	Timeout time.Duration
}

type BoardConfig struct {
	BusID   int           `desc:"I2C Bus to use to communicate with AStar HAT"`
	Addr    byte          `desc:"I2C Address for the AStar HAT"`
	Refresh time.Duration `desc:"How many times per second to refresh the buffer from I2C"`
}

type CameraConfig struct {
	CameraID int    `desc:"Camera ID"`
	Width    int    `desc:"Resolution width"`
	Height   int    `desc:"Resolution height"`
	HFlip    bool   `desc:"Horizontal flip for camera"`
	VFlip    bool   `desc:"Vertical flip for camera"`
	SaveDir  string `desc:"Path where camera saves images and recordings"`
}

// NewConfig returns the default configuration that will be ammended by gflag/flag
func NewConfig() Config {
	return Config{
		Database: DatabaseConfig{
			Path: "/var/commodity/db",
		},
		HTTP: HttpConfig{
			Address: "0.0.0.0:8090",
			SSL:     false,
			Timeout: 3 * time.Second,
		},
		Board: BoardConfig{
			BusID:   1,
			Addr:    0x14,
			Refresh: 50,
		},
		Camera: CameraConfig{
			CameraID: 0,
			Width:    640,
			Height:   480,
			VFlip:    true,
			HFlip:    true,
			SaveDir:  "/data/camera",
		},
	}
}
