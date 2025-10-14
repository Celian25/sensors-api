package mappers

import (
	"fmt"
	"go.mongodb.org/mongo-driver/v2/bson"
	"sensor-api/internal/domain"
	"sensor-api/internal/storage/mongo/repository"
	"time"
)

func ToDomainSensor(s repository.Sensor) domain.Sensor {
	return domain.Sensor{
		ID:       s.ID.Hex(),
		Name:     s.Name,
		Location: s.Location,
		Token:    s.Token,
	}
}

func ToDomainSensorData(sd repository.SensorData) domain.SensorData {
	return domain.SensorData{
		ID:          sd.ID.Hex(),
		SensorID:    sd.SensorID.Hex(),
		Temperature: sd.Temperature,
		Humidity:    sd.Humidity,
		Timestamp:   time.Now(),
	}
}

func FromDomainSensor(s domain.Sensor) (repository.Sensor, error) {
	var id bson.ObjectID
	var err error

	if s.ID != "" {
		id, err = bson.ObjectIDFromHex(s.ID)
		if err != nil {
			return repository.Sensor{}, fmt.Errorf("error converting sensor ID to ObjectID: %w", err)
		}
	}

	return repository.Sensor{
		ID:       id,
		Name:     s.Name,
		Location: s.Location,
		Token:    s.Token,
	}, nil
}

func FromDomainSensorData(sd domain.SensorData) (repository.SensorData, error) {
	var id bson.ObjectID
	var sensorID bson.ObjectID
	var err error

	if sd.ID != "" {
		id, err = bson.ObjectIDFromHex(sd.ID)
		if err != nil {
			return repository.SensorData{}, fmt.Errorf("error converting id to ObjectID: %w", err)
		}
	}
	if sd.SensorID != "" {
		sensorID, err = bson.ObjectIDFromHex(sd.SensorID)
		if err != nil {
			return repository.SensorData{}, fmt.Errorf("error converting sensor ID to ObjectID: %w", err)
		}
	}

	return repository.SensorData{
		ID:          id,
		SensorID:    sensorID,
		Temperature: sd.Temperature,
		Humidity:    sd.Humidity,
		Timestamp:   sd.Timestamp,
	}, nil
}
