package repository

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

type SensorData struct {
	ID          bson.ObjectID `bson:"_id,omitempty"`
	SensorID    bson.ObjectID `bson:"sensor_id"`
	Temperature float64       `bson:"temperature"`
	Humidity    float64       `bson:"humidity"`
	Timestamp   time.Time     `bson:"timestamp"`
}

type Sensor struct {
	ID       bson.ObjectID `bson:"_id,omitempty"`
	Name     string        `bson:"name"`
	Location string        `bson:"location"`
	Token    string        `bson:"token"`
}
