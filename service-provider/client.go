package main

import (
	"context"

	pb "github.com/pranoyk/meetup-demo/common/proto"
	"google.golang.org/grpc"
)

type userClient struct {
	conn pb.UserClient
}

func NewUserClient(conn *grpc.ClientConn) *userClient {
	return &userClient{
		conn: pb.NewUserClient(conn),
	}
}

func (uc *userClient) Validate(ctx context.Context, in *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	return uc.conn.Validate(ctx, in)
	// return &pb.ValidateResponse{
	// 	Valid: true,
	// }, nil
}
