package repository

import (
	"be/graph/model"
	"be/log"
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type VideoRepository interface {
	Save(video *model.Video)
	FindAll() []*model.Video
}

type database struct {
	client *mongo.Client
}

var (
	DATABASE = "dev"
)

// FindAll implements VideoRepository.
func (db *database) FindAll() []*model.Video {
	collection := db.client.Database(DATABASE).Collection("videos")
	cursor, err := collection.Find(context.TODO(), bson.D{})

	if err != nil {
		log.Fatal("FindAll video fail", log.Any("err", err))
	}

	defer cursor.Close(context.TODO())

	var results []*model.Video

	for cursor.Next(context.TODO()) {
		var v *model.Video
		err := cursor.Decode(&v)

		if err != nil {
			log.Fatal("FindAll video fail", log.Any("err", err))
		}

		results = append(results, v)
	}

	return results
}

// Save implements VideoRepository.
func (db *database) Save(video *model.Video) {
	collection := db.client.Database(DATABASE).Collection("videos")
	_, err := collection.InsertOne(context.TODO(), video)

	if err != nil {
		log.Fatal("Save video fail", log.Any("err", err))
	}
}

func New() VideoRepository {

	MONGODB := os.Getenv("MONGODB")

	clientOptions := options.Client().ApplyURI(MONGODB)

	clientOptions = clientOptions.SetMaxPoolSize(50)

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	dbClient, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal("Fail to connect db", log.Any("err", err))
	}

	fmt.Println("Connect to MongoDB")

	return &database{
		client: dbClient,
	}
}
