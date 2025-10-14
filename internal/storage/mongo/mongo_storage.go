package mongo

import (
	"context"
	"sensor-api/internal/domain"
	"sensor-api/internal/mappers"
	"sensor-api/internal/storage/mongo/repository"
	"time"
)

type MongoStorage struct {
	repo *repository.MongoRepository
}

func (m *MongoStorage) GetDevices(ctx context.Context) ([]domain.Sensor, error) {
	devices, err := m.repo.GetDevices(ctx)
	if err != nil {
		return nil, err
	}

	var domainDevices []domain.Sensor
	for _, device := range devices {
		d := mappers.ToDomainSensor(device)
		domainDevices = append(domainDevices, d)
	}
	return domainDevices, nil
}

func (m *MongoStorage) InsertDevice(ctx context.Context, sensor domain.Sensor) error {
	s, err := mappers.FromDomainSensor(sensor)
	if err != nil {
		return err
	}
	err = m.repo.InsertDevice(ctx, s)
	if err != nil {
		return err
	}
	return nil
}

func (m *MongoStorage) InsertData(ctx context.Context, data domain.SensorData) (domain.SensorData, error) {
	d, err := mappers.FromDomainSensorData(data)
	if err != nil {
		return domain.SensorData{}, err
	}
	sensorData, err := m.repo.InsertData(ctx, d)
	if err != nil {
		return domain.SensorData{}, err
	}

	ds := mappers.ToDomainSensorData(sensorData)

	return ds, nil
}

func (m *MongoStorage) GetSearchData(ctx context.Context, sensorID string, from time.Time, to time.Time) ([]domain.SensorData, error) {
	sensorData, err := m.repo.GetSearchData(ctx, sensorID, from, to)
	if err != nil {
		return nil, err
	}

	var domainSensorData []domain.SensorData
	for _, d := range sensorData {
		ds := mappers.ToDomainSensorData(d)
		domainSensorData = append(domainSensorData, ds)
	}

	return domainSensorData, nil
}

func NewMongoStorage(repo *repository.MongoRepository) *MongoStorage {
	return &MongoStorage{repo: repo}
}
