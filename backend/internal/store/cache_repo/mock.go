package cacherepo

import (
	"errors"

	"github.com/hackfeed/remrratality/backend/internal/domain"
)

type CacheRepositoryMock struct{}

func (crm *CacheRepositoryMock) GetMRR(key string) (domain.TotalMRR, error) {
	if key == "user.file-2021-01-02-2021-02-02" {
		return domain.TotalMRR{}, errors.New("error while fetching mrr from cache")
	}
	if key == "user.file-2021-10-01-2021-10-31" {
		return domain.TotalMRR{Total: []float32{0, 0}}, nil
	}
	return domain.TotalMRR{}, nil
}

func (crm *CacheRepositoryMock) SetMRR(key string, mrr domain.TotalMRR) (domain.TotalMRR, error) {
	if key == "errorSetMRR.file-2021-10-01-2021-10-31" {
		return domain.TotalMRR{}, errors.New("error while setting mrr to cache")
	}
	return mrr, nil
}
