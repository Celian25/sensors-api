package database

import (
	"context"
	"database/sql"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	_ "modernc.org/sqlite"
	"os"
	"sensor-api/internal/config"
)

type Database struct {
	DatabaseConfig *config.Database
	DB             *sql.DB
	Mongo          *mongo.Client
}

func (d *Database) Connect(path string) error {
	if d.DatabaseConfig.Type == "sqlite" {
		db, err := sql.Open("sqlite", fmt.Sprintf("%s/%s.db", path, d.DatabaseConfig.Name))
		if err != nil {
			return fmt.Errorf("error opening sqlite database: %w", err)
		}
		d.DB = db
	}

	if d.DatabaseConfig.Type == "mongo" {
		serverAPI := options.ServerAPI(options.ServerAPIVersion1)
		client, err := mongo.Connect(options.Client().ApplyURI(d.DatabaseConfig.URL).SetServerAPIOptions(serverAPI))
		if err != nil {
			return fmt.Errorf("error connecting to mongo: %w", err)
		}
		d.Mongo = client
	}

	return nil
}

func (d *Database) Close(ctx context.Context) {
	if d.DB != nil {
		d.DB.Close()
	}
	if d.Mongo != nil {
		d.Mongo.Disconnect(ctx)
	}
}

func (d *Database) ExecuteMigrations(ctx context.Context, path string) error {
	if d.DatabaseConfig.Type == "mongo" {
		return nil
	}

	bytes, err := os.ReadFile(fmt.Sprintf("%s/schema/schema.sql", path))
	if err != nil {
		return fmt.Errorf("error loading schema: %v", err)
	}

	sqlString := string(bytes)

	tx, err := d.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = tx.ExecContext(ctx, sqlString)
	if err != nil {
		return fmt.Errorf("error executing migrations: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error committing migrations: %w", err)
	}

	return nil
}

func New(database *config.Database) *Database {
	return &Database{
		DatabaseConfig: database,
	}
}
