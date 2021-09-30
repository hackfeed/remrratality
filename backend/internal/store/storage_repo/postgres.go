package storagerepo

import (
	"context"
	"fmt"
	"time"

	"github.com/hackfeed/remrratality/backend/internal/db/storage"
	"github.com/hackfeed/remrratality/backend/internal/domain"
)

type postgresRepo struct {
	storageClient storage.PostgresClient
}

func NewPostgresRepo(storageClient storage.PostgresClient) StorageRepository {
	return &postgresRepo{
		storageClient: storageClient,
	}
}

func (pr *postgresRepo) AddInvoices(invoices []domain.Invoice) ([]domain.Invoice, error) {
	mappedInvoices := make([]storage.Invoice, 0)

	for _, invoice := range invoices {
		mappedInvoice := storage.Invoice{
			UserID:      invoice.UserID,
			FileID:      invoice.FileID,
			CustomerID:  invoice.CustomerID,
			PeriodStart: mapDate(invoice.PeriodStart),
			PaidPlan:    invoice.PaidPlan,
			PaidAmount:  invoice.PaidAmount,
			PeriodEnd:   mapDate(invoice.PeriodEnd),
		}
		mappedInvoices = append(mappedInvoices, mappedInvoice)
	}

	err := pr.storageClient.Create(context.Background(), "invoices", storage.AllFields, mappedInvoices)
	if err != nil {
		return nil, fmt.Errorf("failed to insert invoices, error is: %s", err)
	}

	return invoices, nil
}

func (pr *postgresRepo) GetInvoicesByPeriod(userID, fileID string, periodStart, periodEnd time.Time) ([]domain.Invoice, error) {
	invoices, err := pr.storageClient.ReadByPeriod(
		context.Background(),
		"invoices",
		storage.AllFields,
		userID,
		fileID,
		periodStart,
		periodEnd,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to read invoices by period from %v to %v with user_id %s, file_id %s, error is: %s",
			periodStart, periodEnd, userID, fileID, err,
		)
	}

	mappedInvoices := make([]domain.Invoice, 0)

	for _, invoice := range invoices {
		mappedInvoice := domain.Invoice{
			UserID:      invoice.UserID,
			FileID:      invoice.FileID,
			CustomerID:  invoice.CustomerID,
			PeriodStart: invoice.PeriodStart.Format("2006-01-02"),
			PaidPlan:    invoice.PaidPlan,
			PaidAmount:  invoice.PaidAmount,
			PeriodEnd:   invoice.PeriodEnd.Format("2006-01-02"),
		}
		mappedInvoices = append(mappedInvoices, mappedInvoice)
	}

	return mappedInvoices, nil
}

func (pr *postgresRepo) DeleteInvoices(userID, fileID string) error {
	return pr.storageClient.Delete(context.Background(), "invoices", userID, fileID)
}

func mapDate(date string) time.Time {
	parsed, _ := time.Parse("02.01.2006", date)

	return parsed
}
