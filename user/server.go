package main

import (
	"context"
	"fmt"

	pb "github.com/pranoyk/meetup-demo/common/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type userService struct {
	dbClient *mongo.Client
	pb.UnimplementedUserServer
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
