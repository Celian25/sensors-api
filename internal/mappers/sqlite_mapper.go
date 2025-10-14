package mappers

import (
	"sensor-api/internal/domain"
	"sensor-api/internal/storage/sqlite/repository"
	"time"
)

// DeviceToDomain Device SQLC -> domain.Sensor
func DeviceToDomain(d repository.Device) domain.Sensor {
	return domain.Sensor{
		ID:       d.ID,
		Name:     d.Name,
		Location: d.Location,
		Token:    d.Token,
	}
}

// DomainToInsertDevice domain.Sensor -> InsertDeviceParams
func DomainToInsertDevice(d domain.Sensor) repository.InsertDeviceParams {
	return repository.InsertDeviceParams{
		ID:       d.ID,
		Name:     d.Name,
		Location: d.Location,
		Token:    d.Token,
	}
}

// SensorDatumToDomain SensorDatum SQLC -> domain.SensorData
func SensorDatumToDomain(d repository.SensorDatum) domain.SensorData {
	ts, _ := time.Parse(time.RFC3339, d.Timestamp)
	return domain.SensorData{
		ID:          d.ID,
		SensorID:    d.SensorID,
		Temperature: d.Temperature,
		Humidity:    d.Humidity,
		Timestamp:   ts,
	}
}

// DomainToInsertSensorData domain.SensorData -> InsertSensorDataParams
func DomainToInsertSensorData(d domain.SensorData) repository.InsertSensorDataParams {
	return repository.InsertSensorDataParams{
		ID:          d.SensorID,
		SensorID:    d.SensorID,
		Temperature: d.Temperature,
		Humidity:    d.Humidity,
		Timestamp:   d.Timestamp.Format(time.RFC3339),
	}
}
