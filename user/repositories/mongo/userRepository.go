package mongo

import (
	"context"
	"fmt"
	"reflect"

	"github.com/salman-pathan/go-micro-arch/user/repositories/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	colUser = "User"
)

type UserRepository interface {
	AddUser(ctx context.Context, user model.User) (string, error)
}

type userRepository struct {
	database    string
	mongoClient mongo.Client
}

func NewUserRepository(database string, mongoClient mongo.Client) (UserRepository, error) {
	_, err := mongoClient.Database(database).Collection(colUser).Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	})

	if err != nil {
		return &userRepository{}, err
	}
	return &userRepository{
		database:    database,
		mongoClient: mongoClient,
	}, nil
}

func (r *userRepository) AddUser(ctx context.Context, user model.User) (string, error) {
	res, err := r.mongoClient.Database(r.database).Collection(colUser).InsertOne(ctx, user)
	if err != nil {
		return "", err
	}
	fmt.Println(res.InsertedID)
	fmt.Println(reflect.TypeOf(res.InsertedID))
	return res.InsertedID.(string), nil
}
