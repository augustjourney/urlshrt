package interceptors

import (
	"context"
	"fmt"
	"github.com/augustjourney/urlshrt/internal/logger"
	"google.golang.org/grpc"
	"time"
)

func LogRequests(ctx context.Context, req any,
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	duration := time.Since(start).Milliseconds()
	requestLog := fmt.Sprintf("gRPC Request – Method: %s, Duration: %dms", info.FullMethod, duration)
	if err != nil {
		logger.Log.Error(requestLog, err)
	} else {
		logger.Log.Info(requestLog)
	}

	return resp, err
}
