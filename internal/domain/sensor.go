package domain

import "time"

type SensorData struct {
	ID          string
	SensorID    string
	Temperature float64
	Humidity    float64
	Timestamp   time.Time
}

type Sensor struct {
	ID       string
	Name     string
	Location string
	Token    string
}
