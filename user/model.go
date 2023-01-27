package user

import "google.golang.org/grpc/status"

var (
	ErrAuthTokenNotFound  = status.Error(140201, "Authorization token not found")
	ErrAuthTokenExpired   = status.Error(140202, "Authorization token expired")
	ErrInvalidHeaders     = status.Error(140203, "Invalid Headers")
	ErrUnAuthorizedAccess = status.Error(140204, "Unauthorized Access")
)

type SignupRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required"`
}
