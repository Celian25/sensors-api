package service

import (
	"context"
	"fmt"
	"sensor-api/internal/config"
	"sensor-api/internal/database"
	"sensor-api/internal/domain"
	"sensor-api/internal/storage"
	"sensor-api/internal/utils"
	"sync"
)

type SensorService struct {
	DB      *database.Database
	Repo    storage.Storage
	Devices sync.Map
}

func (s *SensorService) CollectDevices(ctx context.Context, sensors []config.Device) error {
	devices, err := s.Repo.GetDevices(ctx)
	if err != nil {
		return fmt.Errorf("error getting devices from db: %w", err)
	}

	devicesFromDB := make(map[string]string)
	for _, device := range devices {
		devicesFromDB[device.Token] = device.ID
	}

	for _, sensor := range sensors {
		if _, ok := devicesFromDB[sensor.Token]; !ok {
			id := utils.GenerateID(s.DB)
			err = s.Repo.InsertDevice(ctx, domain.Sensor{
				ID:       id,
				Name:     sensor.Name,
				Token:    sensor.Token,
				Location: sensor.Location,
			})
			if err != nil {
				return fmt.Errorf("error inserting device into db: %w", err)
			}
			devicesFromDB[sensor.Token] = id
		}
	}

	for key, value := range devicesFromDB {
		s.Devices.Store(key, value)
	}

	return nil
}

func New(st storage.Storage, db *database.Database) *SensorService {
	return &SensorService{
		DB:      db,
		Repo:    st,
		Devices: sync.Map{},
	}
}
