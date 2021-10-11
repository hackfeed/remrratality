package storagerepo

import (
	"errors"
	"time"

	"github.com/hackfeed/remrratality/backend/internal/domain"
)

type StorageRepositoryMock struct{}

func (prm *StorageRepositoryMock) AddInvoices(invoices []domain.Invoice) ([]domain.Invoice, error) {
	if len(invoices) == 0 {
		return invoices, errors.New("error while adding invoices")
	}
	return invoices, nil
}

func (prm *StorageRepositoryMock) GetInvoicesByPeriod(userID, _ string, _, _ time.Time) ([]domain.Invoice, error) {
	if userID == "errorGetInvoicesByPeriod" {
		return nil, errors.New("error while getting invoices by period")
	}
	if userID == "emptyGetInvoicesByPeriod" {
		return make([]domain.Invoice, 0), nil
	}
	return []domain.Invoice{
		{
			UserID:      "",
			FileID:      "",
			CustomerID:  0,
			PeriodStart: "2021-10-01",
			PaidPlan:    "monthly",
			PaidAmount:  100.0,
			PeriodEnd:   "2021-10-31",
		},
	}, nil
}

func (prm *StorageRepositoryMock) DeleteInvoices(userID, _ string) error {
	if userID == "errorDeleteInvoices" {
		return errors.New("error while deleting invoices")
	}
	return nil
}
