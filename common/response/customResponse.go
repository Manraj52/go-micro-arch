package response

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/status"
)

type CustomResponse interface {
	Error(ctx *gin.Context, err error)
	Success(ctx *gin.Context, data any)
}

type customResponse struct {
	logger *logrus.Entry
}

func NewCustomResponse(logger *logrus.Entry) *customResponse {
	return &customResponse{
		logger: logger,
	}
}

func (cr *customResponse) Error(ctx *gin.Context, err error) {
	cerr, ok := status.FromError(err)
	if !ok {
		cr.logger.Error("err", err)
		ctx.JSON(200, gin.H{
			"success": false,
			"err": gin.H{
				"code": 500,
				// "message": "something went wrong, please try again later",
				"message": err.Error(),
			},
		})
		return
	}

	ctx.JSON(200, gin.H{
		"success": false,
		"err": gin.H{
			"code":    cerr.Code(),
			"message": cerr.Message(),
		},
	})
}

func (cr *customResponse) Success(ctx *gin.Context, data any) {
	ctx.JSON(200, gin.H{
		"success": true,
		"data":    data,
	})
}
