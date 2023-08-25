package service

import (
	"context"
)

// 接口逻辑

type Service struct {
	ctx context.Context
}

func NewService() *Service {
	return &Service{
		ctx: context.Background(),
	}
}
