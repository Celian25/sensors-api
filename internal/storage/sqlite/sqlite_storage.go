package sqlite

import (
	"context"
	"sensor-api/internal/domain"
	"sensor-api/internal/mappers"
	"sensor-api/internal/storage/sqlite/repository"
	"time"
)

type SQLiteStorage struct {
	q *repository.Queries
}

func (s *SQLiteStorage) GetDevices(ctx context.Context) ([]domain.Sensor, error) {
	devices, err := s.q.GetDevices(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]domain.Sensor, len(devices))
	for i, d := range devices {
		result[i] = mappers.DeviceToDomain(d)
	}
	return result, nil
}

func (s *SQLiteStorage) InsertDevice(ctx context.Context, sensor domain.Sensor) error {
	params := mappers.DomainToInsertDevice(sensor)
	return s.q.InsertDevice(ctx, params)
}

func (s *SQLiteStorage) InsertData(ctx context.Context, data domain.SensorData) (domain.SensorData, error) {
	params := mappers.DomainToInsertSensorData(data)
	sd, err := s.q.InsertSensorData(ctx, params)
	if err != nil {
		return domain.SensorData{}, err
	}
	// Возвращаем Sensor, если нужно — можно сделать отдельный маппер
	return domain.SensorData{
		ID:          sd.ID,
		SensorID:    sd.SensorID,
		Temperature: sd.Temperature,
		Humidity:    sd.Humidity,
		Timestamp:   data.Timestamp,
	}, nil
}

func (s *SQLiteStorage) GetSearchData(ctx context.Context, sensorID string, from, to time.Time) ([]domain.SensorData, error) {
	rows, err := s.q.GetSensorDataByTime(ctx, repository.GetSensorDataByTimeParams{
		SensorID:    sensorID,
		Timestamp:   from.Format(time.RFC3339),
		Timestamp_2: to.Format(time.RFC3339),
	})
	if err != nil {
		return nil, err
	}
	result := make([]domain.SensorData, len(rows))
	for i, r := range rows {
		result[i] = mappers.SensorDatumToDomain(r)
	}
	return result, nil
}

func NewSQLiteStorage(q *repository.Queries) *SQLiteStorage {
	return &SQLiteStorage{q: q}
}
