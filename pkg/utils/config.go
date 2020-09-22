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
	Address string `desc:"HTTP address"`
	SSL     bool
	Timeout time.Duration
}

type BoardConfig struct {
	BusID  int               `desc:"I2C Bus to use to communicate with Commodity Board"`
	Addr   byte              `desc:"I2C Address for the Commodity Board"`
	Update BoardUpdateConfig `desc:"map[string]int with the refresh time of each task"`
}

type BoardUpdateConfig struct {
	Power     time.Duration `desc:"Power sensor"`
	Env       time.Duration `desc:"Environmental sensor"`
	IMU       time.Duration `desc:"9 axis IMU sensor"`
	Wifi      time.Duration `desc:"wifi scanner"`
	Heartbeat time.Duration `desc:"Heartbeat"`
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
			BusID: 1,
			Addr:  0x20,
			Update: BoardUpdateConfig{
				Power:     time.Second,
				Env:       time.Second * 3,
				Heartbeat: time.Second,
			},
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
