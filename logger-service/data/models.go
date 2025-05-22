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

func New(mongo *mongo.Client) Models {
	client = mongo
	return Models{
		LogEntry: LogEntry{},
	}
}

type Models struct {
	LogEntry LogEntry
}

type LogEntry struct {
	ID        string    `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string    `json:"name,omitempty" bson:"name,omitempty"`
	Data      string    `json:"data,omitempty" bson:"data,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

func (l *LogEntry) Insert(entry LogEntry) error {
	collection := client.Database("logs").Collection("logs")

	_, err := collection.InsertOne(context.TODO(), LogEntry{
		Name:      entry.Name,
		Data:      entry.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		log.Println("Error inserting log entry:", err)
	}

	return nil
}

func (l *LogEntry) All() ([]*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	collection := client.Database("logs").Collection("logs")
	opts := options.Find()
	opts.SetSort(bson.D{{"created_at", -1}})
	opts.SetLimit(100)
	cursor, err := collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		log.Println("Error finding log entries:", err)
		return nil, err
	}
	defer cursor.Close(ctx)
	var logs []*LogEntry
	for cursor.Next(ctx) {
		var logEntry LogEntry
		if err := cursor.Decode(&logEntry); err != nil {
			log.Println("Error decoding log entry:", err)
			return nil, err
		}
		logs = append(logs, &logEntry)
	}
	return logs, nil
}

func (l *LogEntry) GetOne(id string) (*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	collection := client.Database("logs").Collection("logs")
	var logEntry LogEntry
	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Error converting ID to ObjectID:", err)
		return nil, err
	}
	err = collection.FindOne(ctx, bson.M{"_id": docID}).Decode(&logEntry)
	if err != nil {
		log.Println("Error finding log entry:", err)
		return nil, err
	}
	return &logEntry, nil
}

func (l *LogEntry) DropCollection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	collection := client.Database("logs").Collection("logs")
	err := collection.Drop(ctx)
	if err != nil {
		log.Println("Error dropping collection:", err)
		return err
	}
	return nil
}

func (l *LogEntry) Update(id string, entry LogEntry) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	collection := client.Database("logs").Collection("logs")
	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Error converting ID to ObjectID:", err)
		return nil, err
	}
	result, err := collection.UpdateOne(ctx, bson.M{"_id": docID}, bson.M{
		"$set": bson.M{
			"name":       entry.Name,
			"data":       entry.Data,
			"updated_at": time.Now(),
		},
	})
	if err != nil {
		log.Println("Error updating log entry:", err)
		return nil, err
	}
	return result, nil
}
