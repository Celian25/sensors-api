package structs

import "encoding/json"

type Sensor struct {
	Id       string `json:"id"`
	Location string `json:"location"`
	Name     string `json:"name"`
	Token    string `json:"token"`
}

type SensorsList struct {
	Items []Sensor `json:"items"`
}

type PostSendJSONBody struct {
	Humidity    float64 `json:"humidity"`
	Temperature float64 `json:"temperature"`
}

type SuccessMessage struct {
	Message string `json:"message"`
}

type ErrorMessage struct {
	Message string `json:"message"`
	Error   bool   `json:"error"`
}

type SensorDataSearchResponse struct {
	Items []*SensorData `json:"items"`
}

type SensorData struct {
	ID          string  `json:"id"`
	SensorID    string  `json:"sensor_id"`
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Timestamp   string  `json:"timestamp"`
}

func (nf *PostSendJSONBody) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		return nil
	}
	type Alias PostSendJSONBody
	aux := &Alias{}
	if err := json.Unmarshal(b, aux); err != nil {
		return err
	}
	nf.Humidity = aux.Humidity
	nf.Temperature = aux.Temperature
	return nil
}
