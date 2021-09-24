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
	Port     int
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
	dbURL := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s",
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
		return err
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
	_, err = pc.client.CopyFrom(ctx, pgx.Identifier{table}, fields, pgx.CopyFromRows(data))
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (pc *PostgresClient) CreateDynamic(ctx context.Context, table string, data []interface{}) error {
	vals := ""
	for _, val := range data {
		vals += fmt.Sprintf("%v,", val)
	}
	vals = strings.TrimSuffix(vals, ",")

	_, err := pc.client.Query(ctx, fmt.Sprintf("INSERT INTO %s VALUES (%v)", table, vals))

	return err
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
		"SELECT %s FROM %s WHERE period_start <= '%s' AND period_end >= '%s' AND user_id = '%s' and file_id = '%s'",
		cols,
		table,
		periodStart.Format("2006-01-02"),
		periodEnd.Format("2006-01-02"),
		userID,
		fileID,
	)
	fmt.Println(query)
	rows, err := pc.client.Query(
		ctx,
		query,
	)
	if err != nil {
		return nil, err
	}

	data := []Invoice{}
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
			return nil, err
		}
		data = append(data, invoice)
	}

	return data, nil
}

func (pc *PostgresClient) ReadDynamic(ctx context.Context, table string) ([][]interface{}, error) {
	rows, err := pc.client.Query(ctx, fmt.Sprintf("SELECT * FROM %s", table))
	if err != nil {
		return nil, err
	}

	res := [][]interface{}{}

	cols := rows.FieldDescriptions()
	if err != nil {
		return nil, err
	}
	count := len(cols)
	vals := make([]interface{}, count)
	valsPtrs := make([]interface{}, count)

	for i := range cols {
		valsPtrs[i] = &vals[i]
	}

	for rows.Next() {
		if err := rows.Scan(valsPtrs...); err != nil {
			return nil, err
		}
		row := []interface{}{}
		row = append(row, vals...)
		res = append(res, row)
	}

	return res, nil
}

func (pc *PostgresClient) Delete(ctx context.Context, table, userID, fileID string) error {
	_, err := pc.client.Query(
		ctx,
		fmt.Sprintf(
			"DELETE FROM %s WHERE user_id = '%s' AND file_id = '%s'",
			table,
			userID,
			fileID,
		),
	)

	return err
}
