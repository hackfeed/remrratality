package storage

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PostgresClient struct {
	client *pgxpool.Pool
}

type Options struct {
	Host     string
	Port     string
	User     string
	Password string
	DB       string
}

type Invoice struct {
	UserID      string
	FileID      string
	CustomerID  uint32
	PeriodStart time.Time
	PaidPlan    string
	PaidAmount  float32
	PeriodEnd   time.Time
}

var (
	postgresClient *PostgresClient
	layout         = "2006-01-02"
	AllFields      = []string{
		"user_id",
		"file_id",
		"customer_id",
		"period_start",
		"paid_plan",
		"paid_amount",
		"period_end",
	}
)

func NewPostgresClient(ctx context.Context, options *Options) (*PostgresClient, error) {
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
	dbURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s",
		options.Host, options.Port, options.User, options.DB, options.Password)
	client, err := pgxpool.Connect(ctx, dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect postgres client, error is: %s", err)
	}

	return client, nil
}

func (pc *PostgresClient) Create(ctx context.Context, table string, fields []string, invoices []Invoice) error {
	tx, err := pc.client.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin postgres transaction, error is: %s", err)
	}
	defer tx.Rollback(ctx)

	data := make([][]interface{}, len(invoices))
	for i := range data {
		data[i] = []interface{}{
			invoices[i].UserID,
			invoices[i].FileID,
			invoices[i].CustomerID,
			invoices[i].PeriodStart,
			invoices[i].PaidPlan,
			invoices[i].PaidAmount,
			invoices[i].PeriodEnd,
		}
	}
	if _, err = pc.client.CopyFrom(ctx, pgx.Identifier{table}, fields, pgx.CopyFromRows(data)); err != nil {
		return fmt.Errorf("failed to copy records to postgres from prepared data, error is: %s", err)
	}

	return tx.Commit(ctx)
}

func (pc *PostgresClient) ReadByPeriod(
	ctx context.Context,
	table string,
	fields []string,
	userID, fileID string,
	periodStart, periodEnd time.Time) ([]Invoice, error) {

	cols := ""
	for _, field := range fields {
		cols += fmt.Sprintf("%v,", field)
	}
	cols = strings.TrimSuffix(cols, ",")

	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE period_start >= '%s' AND period_end <= '%s' AND user_id = '%s' and file_id = '%s'",
		cols,
		table,
		periodStart.Format(layout),
		periodEnd.Format(layout),
		userID,
		fileID,
	)
	rows, err := pc.client.Query(
		ctx,
		query,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to run postgres query, error is: %s", err)
	}

	var data []Invoice
	for rows.Next() {
		invoice := Invoice{}
		if err := rows.Scan(
			&invoice.UserID,
			&invoice.FileID,
			&invoice.CustomerID,
			&invoice.PeriodStart,
			&invoice.PaidPlan,
			&invoice.PaidAmount,
			&invoice.PeriodEnd,
		); err != nil {
			return nil, fmt.Errorf("failed to map row to data, error is: %s", err)
		}
		data = append(data, invoice)
	}

	return data, nil
}

func (pc *PostgresClient) Delete(ctx context.Context, table, userID, fileID string) error {
	if _, err := pc.client.Query(
		ctx,
		fmt.Sprintf(
			"DELETE FROM %s WHERE user_id = '%s' AND file_id = '%s'",
			table,
			userID,
			fileID,
		),
	); err != nil {
		return fmt.Errorf("failed to run postgres query, error is: %s", err)
	}

	return nil
}
