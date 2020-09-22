package http

import (
	"context"
	"github.com/bzzzm/commodity-brain/pkg/utils"
	"go.uber.org/zap"
	"net/http"
)

type HTTPServer struct {
	Stream *MJPEGStream
	config utils.HttpConfig
	ctx    context.Context
}

func NewHTTPServer(ctx context.Context, stream *MJPEGStream, config utils.HttpConfig) *HTTPServer {
	return &HTTPServer{
		Stream: stream,
		config: config,
		ctx:    ctx,
	}
}

func (s *HTTPServer) Start() {

	// mjpeg video stream
	http.HandleFunc("/mjpeg", s.Stream.Stream.ServeHTTP)

	// dashboard goes here
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`Commodity v1`))
	})

	server := &http.Server{Addr: s.config.Address}

	go func() {
		<-s.ctx.Done()
		_ = server.Shutdown(s.ctx)
	}()

	err := server.ListenAndServe()
	if err != nil {
		zap.S().Errorf("cannot start http server: %v", err)
	}
}
