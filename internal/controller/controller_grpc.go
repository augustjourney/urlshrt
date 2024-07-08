package controller

import (
	"context"
	"errors"
	pb "github.com/augustjourney/urlshrt/internal/proto"
	"github.com/augustjourney/urlshrt/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Grpc-контроллер
type GrpcController struct {
	service service.IService
	pb.UnimplementedURLServiceServer
}

// Получает полную ссылку по короткой через grpc
func (c *GrpcController) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	var res pb.GetResponse

	originalURL, err := c.service.FindOriginal(req.ShortUrl)
	if errors.Is(err, service.ErrIsDeleted) {
		return &res, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if errors.Is(err, service.ErrNotFound) {
		return &res, status.Errorf(codes.NotFound, err.Error())
	}

	if err != nil {
		return &res, status.Errorf(codes.Internal, err.Error())
	}

	res.OriginalUrl = originalURL

	return &res, nil
}

// Создает короткую ссылку из оригинальной по grpc
func (c *GrpcController) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	var res pb.CreateResponse

	user, err := c.getUserFromMetadata(ctx)

	if err != nil {
		return &res, err
	}

	if req.OriginalUrl == "" {
		return &res, status.Errorf(codes.InvalidArgument, "original url is required")
	}

	result, err := c.service.Shorten(req.OriginalUrl, user)
	if err != nil {
		return &res, status.Errorf(codes.Internal, err.Error())
	}

	if result.AlreadyExists {
		return &res, status.Errorf(codes.AlreadyExists, "short url already exists")
	}

	res.ShortUrl = result.ResultURL

	return &res, nil
}

func (c *GrpcController) CreateBatch(ctx context.Context, req *pb.CreateBatchRequest) (*pb.CreateBatchResponse, error) {
	var res pb.CreateBatchResponse

	user, err := c.getUserFromMetadata(ctx)

	if err != nil {
		return &res, err
	}

	if len(req.Urls) == 0 {
		return &res, status.Errorf(codes.InvalidArgument, "no urls provided")
	}

	body := make([]service.BatchURL, 0)

	for _, url := range req.Urls {
		if url != nil {
			body = append(body, service.BatchURL{
				OriginalURL:   url.OriginalUrl,
				CorrelationID: url.CorrelationId,
			})
		}
	}

	result, err := c.service.ShortenBatch(body, user)
	if err != nil {
		return &res, status.Errorf(codes.Internal, err.Error())
	}

	for _, url := range result {
		res.Urls = append(res.Urls, &pb.BatchURLResult{
			ShortUrl:      url.ShortURL,
			CorrelationId: url.CorrelationID,
		})
	}

	return &res, nil
}

func (c *GrpcController) GetUserURLs(ctx context.Context, req *pb.GetUserURLsRequest) (*pb.GetUserURLsResponse, error) {
	var res pb.GetUserURLsResponse

	user, err := c.getUserFromMetadata(ctx)

	if err != nil {
		return &res, err
	}

	urls, err := c.service.GetUserURLs(context.Background(), user)

	if err != nil {
		return &res, status.Errorf(codes.Internal, err.Error())
	}

	for _, url := range urls {
		res.Urls = append(res.Urls, &pb.UserURL{
			ShortUrl:    url.ShortURL,
			OriginalUrl: url.OriginalURL,
		})
	}

	return &res, nil
}

func (c *GrpcController) DeleteBatch(ctx context.Context, req *pb.DeleteBatchRequest) (*pb.DeleteBatchResponse, error) {
	var res pb.DeleteBatchResponse

	user, err := c.getUserFromMetadata(ctx)

	if err != nil {
		return &res, err
	}

	err = c.service.DeleteBatch(ctx, req.ShortUrls, user)
	if err != nil {
		return &res, status.Errorf(codes.Internal, err.Error())
	}

	return &res, nil
}

func (c *GrpcController) GetStats(ctx context.Context, req *pb.GetStatsRequest) (*pb.GetStatsResponse, error) {
	var res pb.GetStatsResponse

	stats, err := c.service.GetStats(context.Background())
	if err != nil {
		return &res, status.Errorf(codes.Internal, err.Error())
	}

	res.Urls = int32(stats.Urls)
	res.Users = int32(stats.Users)

	return &res, nil
}

// получает пользователя из gprc metadata context
func (c *GrpcController) getUserFromMetadata(ctx context.Context) (string, error) {
	var user string

	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return user, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md.Get("user")

	if len(values) == 0 {
		return user, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	user = values[0]

	if user == "" {
		return user, status.Errorf(codes.Unauthenticated, "user is not provided")
	}

	return user, nil
}

// Создает новый экземпляр grpc-контроллера
func NewGrpcController(service service.IService) *GrpcController {
	return &GrpcController{
		service: service,
	}
}
