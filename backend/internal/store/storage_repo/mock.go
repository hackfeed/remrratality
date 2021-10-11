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

func (prm *StorageRepositoryMock) GetInvoicesByPeriod(_, _ string, _, _ time.Time) ([]domain.Invoice, error) {
	return make([]domain.Invoice, 0), nil
}

func (prm *StorageRepositoryMock) DeleteInvoices(userID, _ string) error {
	if userID == "errorDeleteInvoices" {
		return errors.New("error while deleting invoices")
	}
	return nil
}
