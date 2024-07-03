package controller

import (
	"context"
	"errors"
	pb "github.com/augustjourney/urlshrt/internal/proto"
	"github.com/augustjourney/urlshrt/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Grpc-контроллер
type GrpcController struct {
	service service.IService
	pb.UnimplementedURLServiceServer
}

// Получает полную ссылку по короткой через grpc
func (c *GrpcController) GetURL(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
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
func (c *GrpcController) CreateURL(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	var res pb.CreateResponse

	if req.OriginalUrl == "" {
		return &res, status.Errorf(codes.InvalidArgument, "original url is required")
	}

	result, err := c.service.Shorten(req.OriginalUrl, req.UserId)
	if err != nil {
		return &res, status.Errorf(codes.Internal, err.Error())
	}

	if result.AlreadyExists {
		return &res, status.Errorf(codes.AlreadyExists, "short url already exists")
	}

	res.ShortUrl = result.ResultURL

	return &res, nil
}

func (c *GrpcController) CreateURLBatch() {}

// Создает новый экземпляр grpc-контроллера
func NewGrpcController(service service.IService) *GrpcController {
	return &GrpcController{
		service: service,
	}
}
