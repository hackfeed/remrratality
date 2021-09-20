package main

import (
	"context"
	"fmt"
	"log"

	"github.com/hackfeed/remrratality/backend/internal/db/models"
	"github.com/hackfeed/remrratality/backend/internal/db/storage"
)

var (
	ctx            context.Context
	postgresClient *storage.PostgresClient
)

func init() {
	var err error

	ctx = context.Background()

	postgresClient, err = storage.NewPostgresClient(ctx, &storage.Options{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "pass",
		DB:       "postgres",
	})
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	data := []models.Invoice{
		{"1", "1", 1, "1980-01-01", "flex", 100.0, "1980-01-01"},
		{"1", "1", 1, "1980-01-01", "flex", 100.0, "1980-01-01"},
		{"1", "1", 1, "1980-01-01", "flex", 100.0, "1980-01-01"},
	}
	err := postgresClient.Insert(ctx, "invoices", models.AllFields, data)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("hello")
}
