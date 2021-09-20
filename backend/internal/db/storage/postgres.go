package storage

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/hackfeed/remrratality/backend/internal/db/models"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PostgresClient struct {
	client *pgxpool.Pool
}

type Options struct {
	Host     string
	Port     int
	User     string
	Password string
	DB       string
}

var (
	postgresClient *PostgresClient
	lock           = &sync.Mutex{}
)

func NewPostgresClient(ctx context.Context, options *Options) (*PostgresClient, error) {
	lock.Lock()
	defer lock.Unlock()

	if postgresClient == nil {
		client, err := getPostgresClient(ctx, options)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize postgres client, error is: %s", err)
		}
		postgresClient = &PostgresClient{
			client: client,
		}
		return postgresClient, nil
	}

	return postgresClient, nil
}

func getPostgresClient(ctx context.Context, options *Options) (*pgxpool.Pool, error) {
	dbURL := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s",
		options.Host, options.Port, options.User, options.DB, options.Password)
	client, err := pgxpool.Connect(ctx, dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect postgres client, error is: %s", err)
	}

	return client, nil
}

func (pc *PostgresClient) Insert(ctx context.Context, table string, fields []string, invoices []models.Invoice) error {
	tx, err := pc.client.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	data := make([][]interface{}, len(invoices))
	for i := range data {
		periodStart, _ := time.Parse("2006-01-02", invoices[i].PeriodStart)
		periodEnd, _ := time.Parse("2006-01-02", invoices[i].PeriodEnd)
		data[i] = []interface{}{
			invoices[i].UserID,
			invoices[i].FileID,
			invoices[i].CustomerID,
			periodStart,
			invoices[i].PaidPlan,
			invoices[i].PaidAmount,
			periodEnd,
		}
	}
	_, err = pc.client.CopyFrom(ctx, pgx.Identifier{table}, fields, pgx.CopyFromRows(data))
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
