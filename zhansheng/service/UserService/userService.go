package UserService

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var user_proto = UserService{}

type UserService struct {
}

func (u UserService) UserLogin(ctx context.Context, request *UserLoginRequest) (*UserLoginResponse, error) {
	if request.Username == "admin" && request.Password == "admin123456" {
		return &UserLoginResponse{Status: 0, Msg: "login ok"}, nil
	} else {
		return &UserLoginResponse{Status: 1, Msg: "login failed"}, status.Error(codes.InvalidArgument, "params is wrong")
	}
}

func (u UserService) mustEmbedUnimplementedUserServiceServer() {
	return
}
