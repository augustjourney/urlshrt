package controller

import (
	"context"
	"github.com/augustjourney/urlshrt/internal/app"
	"github.com/augustjourney/urlshrt/internal/config"
	"github.com/augustjourney/urlshrt/internal/logger"
	pb "github.com/augustjourney/urlshrt/internal/proto"
	"github.com/augustjourney/urlshrt/internal/service"
	"github.com/augustjourney/urlshrt/internal/storage"
	"github.com/augustjourney/urlshrt/internal/storage/inmemory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"net"
	"regexp"
	"testing"
)

func newGrpcAppInstance() (pb.URLServiceClient, storage.IRepo, service.Service, func()) {
	cfg := config.New()
	logger.New()

	repo := inmemory.New()
	urlService := service.New(repo, cfg)
	controller := NewGrpcController(&urlService)
	grpcServer := app.NewGrpcApp(controller)

	// Соединение для тестирования
	listener := bufconn.Listen(1024 * 1024)

	// Запуск grpc-сервера
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			panic(err)
		}
	}()

	// Создание grpc-клиента
	dialer := func(ctx context.Context, address string) (net.Conn, error) {
		return listener.Dial()
	}

	conn, err := grpc.DialContext(context.Background(), "bufconn", grpc.WithContextDialer(dialer), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client := pb.NewURLServiceClient(conn)

	cleanup := func() {
		listener.Close()
		grpcServer.Stop()
		conn.Close()
	}

	return client, repo, urlService, cleanup
}

func TestGrpcController_Get(t *testing.T) {
	t.Parallel()
	client, repo, _, cleanup := newGrpcAppInstance()

	t.Cleanup(cleanup)

	url1 := storage.URL{
		UUID:     "some-uuid-1cv23sdfadsfsf",
		Short:    "shrturl1cv23sdfadsfsf",
		Original: "http://google.com?q=1cv23sdfadsfsf",
	}

	url2 := storage.URL{
		UUID:     "some-uuid-211123232",
		Short:    "shrturl211123232",
		Original: "http://yandex.ru?q=211123232",
	}

	repo.Create(context.TODO(), url1)
	repo.Create(context.TODO(), url2)

	tests := []struct {
		name     string
		code     codes.Code
		shortURL string
		response string
	}{
		{
			name:     "Not found url",
			code:     codes.NotFound,
			shortURL: "3453",
			response: "",
		},
		{
			name:     "Found url",
			response: url1.Original,
			code:     codes.OK,
			shortURL: url1.Short,
		},
		{
			name:     "Found url 2",
			code:     codes.OK,
			response: url2.Original,
			shortURL: url2.Short,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := client.Get(context.Background(), &pb.GetRequest{
				ShortUrl: tt.shortURL,
			})
			if tt.code == codes.OK {
				require.NoError(t, err)
				assert.Equal(t, resp.OriginalUrl, tt.response)
			} else {
				errCode, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.code, errCode.Code())
			}
		})
	}
}

func TestGrpcController_Create(t *testing.T) {
	t.Parallel()
	client, _, _, cleanup := newGrpcAppInstance()

	t.Cleanup(cleanup)

	tests := []struct {
		name        string
		code        codes.Code
		originalURL string
	}{
		{
			name:        "URL created",
			code:        codes.OK,
			originalURL: "http://yandex.ru?q=09123123gg",
		},
		{
			name:        "URL created with conflict",
			originalURL: "http://yandex.ru?q=09123123gg",
			code:        codes.AlreadyExists,
		},
		{
			name:        "Empty body",
			originalURL: "",
			code:        codes.InvalidArgument,
		},
	}

	md := metadata.New(map[string]string{
		"user": "user-uuid-0123",
	})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := client.Create(ctx, &pb.CreateRequest{
				OriginalUrl: tt.originalURL,
			})
			if tt.code == codes.OK {
				require.NoError(t, err)
				shortMatch, err := regexp.Match(`\/\w+$`, []byte(resp.ShortUrl))
				require.NoError(t, err)
				assert.True(t, shortMatch)
			} else {
				errCode, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.code, errCode.Code())
			}
		})
	}
}

func TestGrpcController_GetStats(t *testing.T) {
	t.Parallel()
	client, repo, _, cleanup := newGrpcAppInstance()
	t.Cleanup(cleanup)

	ctx := context.Background()

	resp, err := client.GetStats(ctx, &pb.GetStatsRequest{})
	require.NoError(t, err)
	assert.Equal(t, int32(0), resp.Users)
	assert.Equal(t, int32(0), resp.Urls)

	repo.Create(ctx, storage.URL{
		Short:    "123123....=/",
		Original: "http://google.com?q=1cv23sdfadsfsf",
		UUID:     "uid-cxvxv",
		UserUUID: "user-uuid-0123",
	})

	resp, err = client.GetStats(ctx, &pb.GetStatsRequest{})
	require.NoError(t, err)
	assert.Equal(t, int32(1), resp.Users)
	assert.Equal(t, int32(1), resp.Urls)

	repo.Create(context.TODO(), storage.URL{
		Short:    "123123zxcv23.=cv",
		Original: "http://google.com?q=1cvzxccv4f45hbgf563sfsf",
		UUID:     "uid-cxvxv-zxvzcv",
		UserUUID: "user-uuid-0123",
	})

	resp, err = client.GetStats(ctx, &pb.GetStatsRequest{})
	require.NoError(t, err)
	assert.Equal(t, int32(1), resp.Users)
	assert.Equal(t, int32(2), resp.Urls)
}
