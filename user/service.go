package user

import (
	"context"

	log "github.com/sirupsen/logrus"

	user "github.com/salman-pathan/go-micro-arch/user/pb"
)

type Service interface {
	Signup(ctx context.Context, request *user.SignupRequest) (*user.SignupResponse, error)
}

type service struct {
	logger log.Logger
}

func NewService(logger log.Logger) Service {
	return &service{logger: logger}
}

func (s *service) Signup(ctx context.Context, request *user.SignupRequest) (*user.SignupResponse, error) {
	return &user.SignupResponse{
		Id: "sasasaa",
	}, nil
}
