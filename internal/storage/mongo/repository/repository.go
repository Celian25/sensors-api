package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"time"
)

type MongoRepository struct {
	*mongo.Database
}

func (m *MongoRepository) GetDevices(ctx context.Context) ([]Sensor, error) {
	c, err := m.Database.Collection("devices").Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error getting devices collection: %w", err)
	}
	defer c.Close(ctx)

	devices := make([]Sensor, 0)

	for c.Next(ctx) {
		result := Sensor{}
		if err := c.Decode(&result); err != nil {
			return nil, fmt.Errorf("error decoding devices collection: %w", err)
		}
		devices = append(devices, result)
	}

	return devices, nil
}

func (m *MongoRepository) InsertDevice(ctx context.Context, sensor Sensor) error {
	_, err := m.Database.Collection("devices").InsertOne(ctx, sensor)
	if err != nil {
		return fmt.Errorf("error inserting devices collection: %w", err)
	}
	return nil
}

func (m *MongoRepository) InsertData(ctx context.Context, data SensorData) (SensorData, error) {
	c, err := m.Database.Collection("sensors").InsertOne(ctx, data)
	if err != nil {
		return SensorData{}, fmt.Errorf("error inserting devices collection: %w", err)
	}

	id := c.InsertedID

	result := SensorData{}
	err = m.Database.Collection("sensors").FindOne(ctx, bson.M{"_id": id}).Decode(&result)
	if err != nil {
		return SensorData{}, fmt.Errorf("error inserting devices collection: %w", err)
	}
	return result, nil
}

func (m *MongoRepository) GetSearchData(ctx context.Context, sensorID string, from time.Time, to time.Time) ([]SensorData, error) {

	sensor, err := bson.ObjectIDFromHex(sensorID)
	if err != nil {
		return nil, fmt.Errorf("error converting sensor ID from bson.ObjectID: %w", err)
	}

	filter := bson.M{
		"sensor_id": sensor,
		"timestamp": bson.M{
			"$gte": from,
			"$lte": to,
		},
	}
	c, err := m.Database.Collection("sensors").Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error fetching devices data collection: %w", err)
	}
	defer c.Close(ctx)

	sensorsData := make([]SensorData, 0)
	for c.Next(ctx) {
		result := SensorData{}
		if err := c.Decode(&result); err != nil {
			return nil, fmt.Errorf("error decoding devices data collection: %w", err)
		}
		sensorsData = append(sensorsData, result)
	}
	return sensorsData, nil
}

func NewMongoRepository(client *mongo.Client) *MongoRepository {
	db := client.Database("sensor_storage")
	return &MongoRepository{db}
}
