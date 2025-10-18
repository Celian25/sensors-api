package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sensor-api/internal/domain"
	"sensor-api/internal/service"
	_ "sensor-api/internal/storage/sqlite"
	"sensor-api/internal/transport/structs"
	"sensor-api/internal/utils"
	"time"
)

var _ ServerInterface = (*Server)(nil)

type Server struct {
	*service.SensorService
}

func NewServer(s *service.SensorService) *Server {
	return &Server{s}
}

func (s *Server) GetSensors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	sensors, err := s.Repo.GetDevices(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, "internal server error", err)
		return
	}

	var serverSensors []structs.Sensor
	for _, sensor := range sensors {
		s := structs.Sensor{
			Id:       sensor.ID,
			Location: sensor.Location,
			Name:     sensor.Name,
			Token:    sensor.Token,
		}
		serverSensors = append(serverSensors, s)
	}

	sensorItems := &structs.SensorsList{
		Items: serverSensors,
	}

	if err := json.NewEncoder(w).Encode(sensorItems); err != nil {
		respondError(w, http.StatusInternalServerError, "internal server error", err)
		return
	}

}

func (s *Server) PostSend(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	data, err := io.ReadAll(r.Body)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Cannot close body", err)
		return
	}

	var req structs.PostSendJSONBody
	err = json.Unmarshal(data, &req)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	token := r.Header.Get("Token")
	if token == "" {
		respondError(w, http.StatusUnauthorized, "Missing Token header", nil)
		return
	}

	device, ok := s.Devices.Load(token)
	if !ok {
		respondError(w, http.StatusForbidden, "unknown device", nil)
		return
	}

	id := utils.GenerateID(s.DB)
	_, err = s.Repo.InsertData(r.Context(), domain.SensorData{
		ID:          id,
		SensorID:    device.(string),
		Temperature: req.Temperature,
		Humidity:    req.Humidity,
		Timestamp:   time.Now(),
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "error inserting sensor data", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	var msg structs.SuccessMessage
	msg.Message = "success"
	err = json.NewEncoder(w).Encode(&msg)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "error encoding success message", err)
		return
	}

	log.Println("Successfully created sensor data message to database")
}

func (s *Server) GetDataSearch(w http.ResponseWriter, r *http.Request, params GetDataSearchParams) {
	w.Header().Set("Content-Type", "application/json")

	token := r.Header.Get("Token")
	if token == "" {
		respondError(w, http.StatusUnauthorized, "Missing Token header", nil)
		return
	}

	device, ok := s.Devices.Load(token)
	if !ok {
		respondError(w, http.StatusBadRequest, "invalid device id", nil)
		return
	}

	if params.From.IsZero() {
		respondError(w, http.StatusBadRequest, "invalid parameter", nil)
		return
	}

	if params.To.IsZero() {
		respondError(w, http.StatusBadRequest, "invalid parameter", nil)
		return
	}

	deviceData, err := s.Repo.GetSearchData(r.Context(), device.(string), params.From, params.To)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "error getting device data", err)
		return
	}

	if len(deviceData) == 0 {
		respondError(w, http.StatusNotFound, "data for device not found", err)
		return
	}

	items := make([]*structs.SensorData, 0)
	for _, v := range deviceData {
		item := &structs.SensorData{
			ID:          v.ID,
			SensorID:    v.SensorID,
			Temperature: v.Temperature,
			Humidity:    v.Humidity,
			Timestamp:   v.Timestamp.String(),
		}
		items = append(items, item)
	}
	deviceResponse := structs.SensorDataSearchResponse{
		Items: items,
	}
	err = json.NewEncoder(w).Encode(&deviceResponse)
	if err != nil {
		respondError(w, http.StatusBadRequest, "error encoding response", err)
		return
	}

}

func respondError(w http.ResponseWriter, code int, message string, err error) {
	w.WriteHeader(code)
	errorMessage := structs.ErrorMessage{
		Error:   true,
		Message: message,
	}
	if err != nil {
		log.Println(fmt.Sprintf("%s:%v", errorMessage.Message, err))
	} else {
		log.Println(errorMessage.Message)
	}
	_ = json.NewEncoder(w).Encode(&errorMessage)
}
