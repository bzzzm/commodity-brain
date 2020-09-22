package camera

import (
	"context"
	"fmt"
	"github.com/bzzzm/commodity-brain/pkg/utils"
	"go.uber.org/zap"
	"gocv.io/x/gocv"
	"image"
	"sync"
	"time"
)

const (
	LogFramesSkipped = 10
	LogReadErrors    = 10
)

// CaptureImage is a representation of a captured image with annotations. Image is currently not ued.
type CaptureImage struct {
	Mat     gocv.Mat
	Image   image.Image
	Time    time.Time
	FrameId int
}

// CaptureInterface implements Start, Stop and Shutter for Camera Captures
type CaptureInterface interface {
	Start()
	Stop()
	Shutter()
}

// OpenCVCapture is the capture struct using opencv. Implements CaptureInterface.
type OpenCVCapture struct {
	Camera       *gocv.VideoCapture
	ImageChannel chan CaptureImage
	ctx          context.Context
	config       utils.CameraConfig
	mux          *sync.Mutex
}

// NewCapture creates a OpenCVCapture
func NewCapture(ctx context.Context, config utils.CameraConfig, ic chan CaptureImage) (*OpenCVCapture, error) {

	zap.S().Infof("opening capture device id=%v", config.CameraID)

	// open webcam
	camera, err := gocv.OpenVideoCapture(config.CameraID)
	if err != nil {
		zap.S().Errorf("error opening capture device %v %v", config.CameraID, err)
		return nil, err
	}
	zap.S().Infof("opened capture using codec id=%v", camera.CodecString())

	// set height and width
	// todo: add more settings to ensure consistent white balance, etc
	camera.Set(gocv.VideoCaptureFrameWidth, float64(config.Width))
	camera.Set(gocv.VideoCaptureFrameHeight, float64(config.Height))

	capture := &OpenCVCapture{
		Camera:       camera,
		ImageChannel: ic,
		ctx:          ctx,
		config:       config,
		mux:          &sync.Mutex{},
	}
	return capture, nil
}

// Start start collecting images from c.Camera and sends them as CaptureImage on c.ImageChannel
func (c *OpenCVCapture) Start() {
	var skips int
	var reads int
	var frameId int

	// define a new opencv matrix
	frame := gocv.NewMat()
	flipMat := gocv.NewMat()
	defer frame.Close()
	defer flipMat.Close()

	// wait for images
	for {
		select {
		case <-c.ctx.Done():
			zap.S().Debug("closing camera capture")
			return
		default:

			// used by image consumers to generate end to end timing
			startTime := time.Now()

			// if the channel is full, skip the frame and log every 10 skipped frames as an event
			if len(c.ImageChannel) == cap(c.ImageChannel) {
				if skips%LogFramesSkipped == 0 {
					skips++
					zap.S().Warnf("camera skipped %v frames", LogFramesSkipped)
				}
				continue
			}

			// read a frame into the mat
			if ok := c.Camera.Read(&frame); !ok {
				if reads%LogReadErrors == 0 {
					reads++
					zap.S().Warnf("camera could not read %v frames", LogReadErrors)
				}
				continue
			}

			if frame.Empty() {
				continue
			}

			// create the default capturedImage
			capturedImage := CaptureImage{
				Mat:     frame,
				Time:    startTime,
				FrameId: frameId,
			}

			// if flip is required
			if c.config.VFlip && c.config.HFlip {
				if c.config.VFlip {
					gocv.Flip(frame, &flipMat, -1)
					capturedImage.Mat = flipMat
				}
			} else {
				if c.config.VFlip {
					gocv.Flip(frame, &flipMat, 0)
					capturedImage.Mat = flipMat
				}

				if c.config.HFlip {
					gocv.Flip(frame, &flipMat, 1)
					capturedImage.Mat = flipMat
				}
			}

			c.ImageChannel <- capturedImage
			frameId++
		}
	}
}

// Shutter get the first available image from c.ImageChannel, encodes it in JPEG and writes it to config.Camera.SaveDir
func (c *OpenCVCapture) Shutter() (string, error) {
	img := <-c.ImageChannel

	filename := buildFileName(c.config.SaveDir)

	if ok := gocv.IMWrite(filename, img.Mat); !ok {
		return "", fmt.Errorf("could not write shutter image in %v", filename)
	}

	zap.S().Infof("captured shutter image in %v", filename)
	return filename, nil
}

// Close stops collecting images
func (c *OpenCVCapture) Close() {
	_ = c.Camera.Close()
}

// buildFileName builds the path for the captured image
func buildFileName(path string) string {
	return fmt.Sprintf("%v/capture-%v.jpg", path, time.Now().Format("20060102150405"))
}
