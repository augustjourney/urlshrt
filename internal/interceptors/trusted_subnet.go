package interceptors

import (
	"context"
	"github.com/augustjourney/urlshrt/internal/config"
	"github.com/augustjourney/urlshrt/internal/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"net"
)

// интерсептор, который проверяет находится ли ip-адрес клиента
// в доверенной подсети из конфига TrustedSubnet
func IPInTrustedSubnet(ctx context.Context, req any,
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	requestPeer, ok := peer.FromContext(ctx)
	if !ok {
		return nil, status.Error(codes.FailedPrecondition, "peer not found")
	}

	ip := requestPeer.Addr.String()

	if ip == "" {
		return nil, status.Errorf(codes.PermissionDenied, "forbidden")
	}

	cfg := config.New()

	if cfg.TrustedSubnet == "" {
		return handler(ctx, req)
	}

	_, subnet, err := net.ParseCIDR(cfg.TrustedSubnet)
	if err != nil {
		logger.Log.Error("could not parse cidr in IPInTrustedSubnet", err.Error())
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	if subnet == nil || !subnet.Contains(net.ParseIP(ip)) {
		return nil, status.Errorf(codes.PermissionDenied, "forbidden")
	}

	return handler(ctx, req)
}
