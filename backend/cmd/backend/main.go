package main

import (
	"context"
	"fmt"
	"log"

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
	// time, _ := time.Parse("2006-01-02", "1980-01-01")
	// data := []models.Invoice{
	// 	{"1", "1", 1, time, "flex", 100.0, time},
	// 	{"1", "1", 1, time, "flex", 100.0, time},
	// 	{"1", "1", 1, time, "flex", 100.0, time},
	// }
	// err := postgresClient.Insert(ctx, "invoices", models.AllFields, data)
	// data, err := postgresClient.SelectByPeriod(ctx, "invoices", models.AllFields, "1", "1", time, time)
	// data, err := postgresClient.SelectDynamic(ctx, "invoices")
	err := postgresClient.Delete(ctx, "invoices", "1", "1")
	if err != nil {
		fmt.Println(err)
	}
}
