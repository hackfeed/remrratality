package cacherepo

import "github.com/hackfeed/remrratality/backend/internal/domain"

type CacheRepository interface {
	GetMRR(string) (domain.TotalMRR, error)
	SetMRR(string, domain.TotalMRR) (domain.TotalMRR, error)
}
