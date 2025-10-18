package websocket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sensor-api/internal/service"
	"sensor-api/internal/transport/structs"
	"sync"
)

type Ws struct {
	*service.SensorService
	*websocket.Upgrader
	clients sync.Map
}

func (ws *Ws) SendData(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Token")
	if token == "" {
		http.Error(w, "Missing Token", http.StatusUnauthorized)
		log.Println("Missing Token")
		return
	}

	deviceID, ok := ws.Devices.Load(token)
	if !ok {
		http.Error(w, "Device not found", http.StatusUnauthorized)
		log.Println("Device not found")
		return
	}

	c, err := ws.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading to websocket: %v", err)
	}
	defer c.Close()

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Printf("Error reading from websocket: %v", err)
			break
		}
		if mt == websocket.BinaryMessage {
			log.Print("Unsupported binary websocket message type - binary message")
			break
		}

		log.Printf("got message: %s", string(message))

		data := &structs.PostSendJSONBody{}
		err = json.Unmarshal(message, data)
		if err != nil {
			log.Printf("Error unmarshalling from websocket: %v", err)
			break
		}
		ws.clients.Store(deviceID.(string), data)

		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Printf("Error writing from websocket: %v", err)
			break
		}
	}
}

func (ws *Ws) ReadData(w http.ResponseWriter, r *http.Request) {
	c, err := ws.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading to websocket: %v", err)
	}
	defer c.Close()

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Printf("Error reading from websocket: %v", err)
			break
		}
		if mt == websocket.BinaryMessage {
			log.Print("Unsupported binary websocket message type - binary message")
			break
		}

		sensorID := string(message)

		sensorData, ok := ws.clients.Load(sensorID)
		if !ok {
			log.Printf("Error reading from websocket: %v", err)
			break
		}

		data, ok := sensorData.(*structs.PostSendJSONBody)
		if !ok {
			log.Printf("Error reading from websocket: %v", err)
			break
		}

		bytes, err := json.Marshal(data)
		if err != nil {
			log.Printf("Error marshalling from websocket: %v", err)
			break
		}

		err = c.WriteMessage(mt, bytes)
		if err != nil {
			log.Printf("Error writing from websocket: %v", err)
			break
		}
	}
}

func NewWs(s *service.SensorService) *Ws {
	return &Ws{
		SensorService: s,
		Upgrader: &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		clients: sync.Map{},
	}
}
