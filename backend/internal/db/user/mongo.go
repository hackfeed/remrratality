package user

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	client *mongo.Client
}

type Options struct {
	Host     string
	Port     string
	User     string
	Password string
	DB       string
}

type File struct {
	Name       string    `bson:"name"`
	UploadedAt time.Time `bson:"uploaded_at"`
}

type User struct {
	ID           primitive.ObjectID `bson:"_id"`
	UserID       string             `bson:"user_id"`
	Email        *string            `bson:"email" validate:"email,required"`
	Password     *string            `bson:"password" validate:"required,min=6"`
	Token        *string            `bson:"token"`
	RefreshToken *string            `bson:"refresh_token"`
	CreatedAt    time.Time          `bson:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at"`
	Files        []File             `bson:"files"`
}

var mongoClient *MongoClient

func NewMongoClient(ctx context.Context, options *Options) (*MongoClient, error) {
	if mongoClient == nil {
		client, err := getMongoClient(ctx, options)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize mongo client, error is: %s", err)
		}
		mongoClient = &MongoClient{
			client: client,
		}
		return mongoClient, nil
	}

	return mongoClient, nil
}

func getMongoClient(ctx context.Context, opts *Options) (*mongo.Client, error) {
	dbURL := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", opts.User, opts.Password, opts.Host, opts.Port, opts.DB)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURL))
	if err != nil {
		return nil, fmt.Errorf("failed to connect mongo client, error is: %s", err)
	}

	return client, nil
}

func (mc *MongoClient) Create(user User) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	if _, err := mc.client.Database("mrr").Collection("user").InsertOne(ctx, user); err != nil {
		return User{}, fmt.Errorf("failed to run mongo insertOne method, error is: %s", err)
	}

	return user, nil
}

func (mc *MongoClient) Read(key string, value interface{}) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var user User

	if err := mc.client.Database("mrr").Collection("user").FindOne(ctx, bson.M{key: value}).Decode(&user); err != nil {
		return User{}, fmt.Errorf("failed to run mongo decode method, error is: %s", err)
	}

	return user, nil
}

func (mc *MongoClient) Update(obj interface{}, key string, value interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	upsert := true
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	if _, err := mc.
		client.
		Database("mrr").
		Collection("user").
		UpdateOne(
			ctx,
			bson.M{key: value},
			bson.D{{Key: "$set", Value: obj}},
			&opt,
		); err != nil {
		return fmt.Errorf("failed to run mongo updateOne method, error is: %s", err)
	}

	return nil
}
