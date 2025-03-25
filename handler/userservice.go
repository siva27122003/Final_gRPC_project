package Server

import (
	Config "GRPC/Config"
	"GRPC/model"
	"GRPC/pb"
	"context"
	"errors"
	"log"
)

type Server struct {
	pb.UnimplementedUserServiceServer
}

func (s *Server) CreateUser(ctx context.Context, in *pb.User) (*pb.UserResponse, error) {
	log.Printf("Received: %v", in.GetUserName())
	user := &model.User{
		Name:     in.GetUserName(),
		Email:    in.GetEmail(),
		Phone:    in.GetPhoneNumber(),
		Password: in.GetPassword(),
		Role:     in.GetRole(),
		Location: in.GetLocation(),
	}
	res := Config.DB.Create(user)
	if res.RowsAffected == 0 {
		return nil, errors.New("error saving user details")
	}
	return &pb.UserResponse{User: in}, nil
}
