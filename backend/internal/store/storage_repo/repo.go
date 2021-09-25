package storagerepo

import (
	"time"

	"github.com/hackfeed/remrratality/backend/internal/domain"
)

type StorageRepository interface {
	AddInvoices([]domain.Invoice) ([]domain.Invoice, error)
	GetInvoicesByPeriod(string, string, time.Time, time.Time) ([]domain.Invoice, error)
	DeleteInvoices(string, string) error
}
