package utils

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"sensor-api/internal/database"
)

func GenerateID(db *database.Database) string {
	switch db.DatabaseConfig.Type {
	case "mongo":
		return bson.NewObjectID().Hex()
	case "sqlite":
		id, _ := uuid.NewV7()
		return id.String()
	default:
		panic("Unknown database type")
	}
}
