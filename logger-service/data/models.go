package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func New(mongoClient *mongo.Client) *Models {
	client = mongoClient

	return &Models{}
}

type Models struct {
	LogEntry LogEntry
}

type LogEntry struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string    `bson:"name" json:"name"`
	Data      string    `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_at" json:"created"`
	UpdatedAt time.Time `bson:"updated _at" json:"updated "`
}

func (le *LogEntry) Insert(entry LogEntry) error {
	collection := client.Database("logs").Collection("logs")

	_, err := collection.InsertOne(context.TODO(), entry)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (le *LogEntry) GetAll() ([]*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	opts := options.Find()
	opts.SetSort(bson.D{{"created_at", -1}})

	cursor, err := collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		log.Println("error finding all records", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	logs := []*LogEntry{}
	for cursor.Next(ctx) {
		entry := LogEntry{}

		err = cursor.Decode(&entry)
		if err != nil {
			log.Println("error decoding entry", err)
			return nil, err
		} else {
			logs = append(logs, &entry)
		}
	}

	return logs, nil
}

func (le *LogEntry) GetOne(id string) (*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("error formatting object ID", err)
		return nil, err
	}

	entry := LogEntry{}
	err = collection.FindOne(ctx, bson.M{"_id": docID}).Decode(&entry)
	if err != nil {
		log.Println("error finding entry", err)
		return nil, err
	}

	return &entry, nil
}

func (le *LogEntry) DropCollection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	err := collection.Drop(ctx)
	if err != nil {
		log.Println("error dropping collection", err)
		return err
	}

	return nil
}

func (le *LogEntry) Update() (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	docID, err := primitive.ObjectIDFromHex(le.ID)
	if err != nil {
		log.Println("error formatting object ID", err)
		return nil, err
	}

	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": docID},
		bson.D{
			{"$set", bson.D{
				{"name", le.Name},
				{"data", le.Data},
				{"updated_at", le.UpdatedAt},
			}},
		},
	)

	if err != nil {
		log.Println("error updating log entry", err)
		return nil, err
	}

	return result, nil
}
