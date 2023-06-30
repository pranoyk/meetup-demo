package main

import (
	"context"
	"fmt"

	pb "github.com/pranoyk/meetup-demo/common/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userService struct {
	dbClient *mongo.Client
	pb.UnimplementedUsersServer
}

func NewUserService(dbClient *mongo.Client) *userService {
	return &userService{dbClient: dbClient}
}

func (us *userService) Validate(ctx context.Context, in *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	fmt.Println("Validating user")
	usersCollection := us.dbClient.Database("user-management").Collection("users")
	filter := bson.M{"email": in.Email, "name": in.Name}
	err := usersCollection.FindOne(context.Background(), filter).Decode(&in)
	if err != nil {
		fmt.Printf("Invalid user or email: %v", err)
		return &pb.ValidateResponse{Valid: false}, nil
	}
	return &pb.ValidateResponse{Valid: true}, nil
}

func (us *userService) GetUsers(_ *pb.Empty, stream pb.Users_GetUsersServer) error {
	fmt.Println("Getting users")
	usersCollection := us.dbClient.Database("user-management").Collection("users")
	cursor, err := usersCollection.Find(context.Background(), bson.M{})
	if err != nil {
		fmt.Printf("Error getting users: %v", err)
		return err
	}
	var users []*pb.User
	for cursor.Next(context.Background()) {
		var user *pb.User
		err := cursor.Decode(&user)
		if err != nil {
			fmt.Printf("Error decoding user: %v", err)
			return err
		}
		users = append(users, user)
	}
	for _, result := range users {
		if err := stream.Send(result); err != nil {
			fmt.Println("Error sending user")
			return status.Error(codes.Internal, "Error sending user")
		}
	}
	return nil
}