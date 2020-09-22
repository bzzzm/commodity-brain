package http

import (
	"context"
	"github.com/bzzzm/commodity-brain/pkg/camera"
	"github.com/hybridgroup/mjpeg"
	"go.uber.org/zap"
	"gocv.io/x/gocv"
	"time"
)

// LogFrameTime represents the number of frames between logging (debug) how much time it takes
// to process a frame
const LogFrameTime = 100

// StreamInterface
type StreamInterface interface {
	Start()
	Stop()
}

// MJPEGStream struct is handling images coming on ImageChannel and saves the in the mjpeg buffer
// so HTTP server can send it to the clients
type MJPEGStream struct {
	Stream *mjpeg.Stream
	ic     <-chan camera.CaptureImage
	ctx    context.Context
}

// NewStream creates new MPJEGStream
func NewStream(ctx context.Context, ic <-chan camera.CaptureImage) *MJPEGStream {
	stream := mjpeg.NewStream()
	return &MJPEGStream{
		Stream: stream,
		ic:     ic,
		ctx:    ctx,
	}
}

// Start will start collecting images from ImageChannel and save them in mjpeg buffer
func (s *MJPEGStream) Start() {
	var updates int
	for {
		select {
		case <-s.ctx.Done():
			zap.S().Debug("closing camera stream")
			return
		case img := <-s.ic:
			buf, err := gocv.IMEncode(gocv.JPEGFileExt, img.Mat)
			if err != nil {
				zap.S().Errorf("cannot encode stream frame: %v", err)
				continue
			}
			s.Stream.UpdateJPEG(buf)

			if updates%LogFrameTime == 0 {
				zap.S().Debugf("sent stream frame in %v", time.Since(img.Time))
			}
			updates++
		}
	}
}

func (s *MJPEGStream) Stop() {

}
