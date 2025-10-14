package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sensor-api/internal/config"
	"sensor-api/internal/database"
	"sensor-api/internal/storage/mongo"
	repository2 "sensor-api/internal/storage/mongo/repository"
	"sensor-api/internal/transport/mux"
	"sensor-api/internal/transport/rest"
	"sensor-api/internal/transport/websocket"

	"sensor-api/internal/service"
)

func main() {
	var rootDir string
	flag.StringVar(&rootDir, "root", "./../../", "Root directory")
	flag.Parse()

	cfg, err := config.Load(rootDir)
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	dbConf := database.New(&cfg.Database)

	err = dbConf.Connect(rootDir)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}

	defer dbConf.Close(context.Background())

	err = dbConf.ExecuteMigrations(context.Background(), rootDir)
	if err != nil {
		log.Fatalf("error executing migrations: %v", err)
	}

	//queries := repository.New(dbConf.DB)
	//storage := sqlite.NewSQLiteStorage(queries)
	repo := repository2.NewMongoRepository(dbConf.Mongo)

	storage := mongo.NewMongoStorage(repo)
	sensorService := service.New(storage, dbConf)

	err = sensorService.CollectDevices(context.Background(), cfg.Devices)
	if err != nil {
		log.Fatalf("error collecting devices: %v", err)
	}

	s := rest.NewServer(sensorService)

	ws := websocket.NewWs(sensorService)

	r := mux.GetMux(rootDir, "8000", &cfg.Server, ws)
	h := rest.HandlerFromMux(s, r)
	serv := &http.Server{
		Handler: h,
		Addr:    "0.0.0.0:" + fmt.Sprint(cfg.Server.Port),
	}

	log.Printf("serving on %s", serv.Addr)

	// And we serve HTTP until the world ends.
	err = serv.ListenAndServe()
	if err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}
