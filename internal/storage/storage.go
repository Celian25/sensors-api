package storage

import (
	"context"
	"sensor-api/internal/domain"
	"time"
)

type Storage interface {
	GetDevices(ctx context.Context) ([]domain.Sensor, error)
	InsertDevice(ctx context.Context, sensor domain.Sensor) error
	InsertData(ctx context.Context, data domain.SensorData) (domain.SensorData, error)
	GetSearchData(ctx context.Context, sensorID string, from time.Time, to time.Time) ([]domain.SensorData, error)
}
