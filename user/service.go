package user

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	user "github.com/salman-pathan/go-micro-arch/user/pb"
	"github.com/salman-pathan/go-micro-arch/user/repositories/model"
	mongorepository "github.com/salman-pathan/go-micro-arch/user/repositories/mongo"
)

type Service interface {
	Signup(ctx context.Context, request *user.SignupRequest) (*user.SignupResponse, error)
}

type service struct {
	logger   log.Logger
	location *time.Location
	userRepo mongorepository.UserRepository
}

func NewService(
	logger log.Logger,
	location *time.Location,
	userRepo mongorepository.UserRepository,
) Service {
	return &service{
		logger:   logger,
		location: location,
		userRepo: userRepo,
	}
}

func (s *service) Signup(ctx context.Context, request *user.SignupRequest) (*user.SignupResponse, error) {
	userModel := model.NewUser(request.FirstName, request.LastName, request.Email, request.Password, time.Now().In(s.location))
	id, err := s.userRepo.AddUser(ctx, userModel)
	if err != nil {
		s.logger.Errorln("Signup | err : ", err)
		return nil, err
	}

	return &user.SignupResponse{
		Id: id,
	}, nil
}
