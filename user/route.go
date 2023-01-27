package user

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/salman-pathan/go-micro-arch/common/middlewares"
	"github.com/salman-pathan/go-micro-arch/common/response"
	user "github.com/salman-pathan/go-micro-arch/user/pb"
)

type Routes struct {
	logger         *log.Entry
	routerGroup    *gin.RouterGroup
	userClient     user.UserClient
	authMiddleware middlewares.AuthMiddleware
	custResponse   response.CustomResponse
}

func NewUserRoutes(
	logger *log.Entry,
	routerGroup *gin.RouterGroup,
	userClient user.UserClient,
	authMiddleware middlewares.AuthMiddleware,
	custResponse response.CustomResponse,
) Routes {
	return Routes{
		logger:         logger,
		routerGroup:    routerGroup,
		userClient:     userClient,
		authMiddleware: authMiddleware,
		custResponse:   custResponse,
	}
}

func (r *Routes) RegisterRoutes() {
	r.routerGroup.POST("/signup", r.signup)
}

func (r *Routes) signup(ctx *gin.Context) {
	var signupReq SignupRequest
	if err := ctx.ShouldBindJSON(&signupReq); err != nil {
		r.custResponse.Error(ctx, err)
		return
	}

	req := &user.SignupRequest{
		FirstName: signupReq.FirstName,
		LastName:  signupReq.LastName,
		Email:     signupReq.Email,
		Password:  signupReq.Password,
	}

	res, err := r.userClient.Signup(ctx, req)
	if err != nil {
		r.logger.Error("err", err)
		r.custResponse.Error(ctx, err)
		return
	}

	r.logger.Debug("res", res.Id)
	r.custResponse.Success(ctx, res)
}
