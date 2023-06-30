package main

import (
	"context"
	"io"

	pb "github.com/pranoyk/meetup-demo/common/proto"
	"google.golang.org/grpc"
)

type userClient struct {
	conn pb.UsersClient
}

func NewUserClient(conn *grpc.ClientConn) *userClient {
	return &userClient{
		conn: pb.NewUsersClient(conn),
	}
}

func (uc *userClient) Validate(ctx context.Context, in *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	return uc.conn.Validate(ctx, in)
}

func (uc *userClient) GetUsers(ctx context.Context, in *pb.Empty) ([]*pb.User, error) {
	results := []*pb.User{}
	stream, err :=  uc.conn.GetUsers(ctx, in)
	if err != nil {
		return results, err
	}
	for {
		result, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return results, err
		}
		results = append(results, result)
	}
	return results, nil
}