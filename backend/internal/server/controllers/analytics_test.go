package controllers

import (
	"errors"
	"testing"
	"time"

	"github.com/hackfeed/remrratality/backend/internal/domain"
	cacherepo "github.com/hackfeed/remrratality/backend/internal/store/cache_repo"
	storagerepo "github.com/hackfeed/remrratality/backend/internal/store/storage_repo"
	"github.com/stretchr/testify/assert"
)

func TestCreateAnalytics(t *testing.T) {
	type testInput struct {
		userID, fileID, periodStart, periodEnd string
	}
	type testWant struct {
		months []string
		mrr    domain.TotalMRR
		err    error
	}

	tests := []struct {
		input testInput
		want  testWant
	}{
		{
			input: testInput{
				userID:      "",
				fileID:      "",
				periodStart: "wrongPeriod",
				periodEnd:   "",
			},
			want: testWant{
				months: nil,
				mrr:    domain.TotalMRR{},
				err:    errors.New("failed parse period start date, error is: parsing time \"wrongPeriod\" as \"2006-01-02\": cannot parse \"wrongPeriod\" as \"2006\""),
			},
		},
		{
			input: testInput{
				userID:      "",
				fileID:      "",
				periodStart: "2006-01-02",
				periodEnd:   "wrongPeriod",
			},
			want: testWant{
				months: nil,
				mrr:    domain.TotalMRR{},
				err:    errors.New("failed parse period end date, error is: parsing time \"wrongPeriod\" as \"2006-01-02\": cannot parse \"wrongPeriod\" as \"2006\""),
			},
		},
		{
			input: testInput{
				userID:      "",
				fileID:      "",
				periodStart: "2021-02-02",
				periodEnd:   "2021-01-02",
			},
			want: testWant{
				months: nil,
				mrr:    domain.TotalMRR{},
				err:    errors.New("period start should be less than period end"),
			},
		},
		{
			input: testInput{
				userID:      "user",
				fileID:      "file",
				periodStart: "2021-01-02",
				periodEnd:   "2021-02-02",
			},
			want: testWant{
				months: nil,
				mrr:    domain.TotalMRR{},
				err:    errors.New("failed to get mrr from cache, error is: error while fetching mrr from cache"),
			},
		},
		{
			input: testInput{
				userID:      "user",
				fileID:      "file",
				periodStart: "2021-10-01",
				periodEnd:   "2021-10-31",
			},
			want: testWant{
				months: []string{"10.2021"},
				mrr:    domain.TotalMRR{Total: []float32{0, 0}},
				err:    nil,
			},
		},
		{
			input: testInput{
				userID:      "errorGetInvoicesByPeriod",
				fileID:      "file",
				periodStart: "2021-10-01",
				periodEnd:   "2021-10-31",
			},
			want: testWant{
				months: []string{"10.2021"},
				mrr:    domain.TotalMRR{},
				err:    errors.New("failed to form mpp, error is: failed to get invoices from storage, error is: error while getting invoices by period"),
			},
		},
		{
			input: testInput{
				userID:      "errorSetMRR",
				fileID:      "file",
				periodStart: "2021-10-01",
				periodEnd:   "2021-10-31",
			},
			want: testWant{
				months: []string{"10.2021"},
				mrr: domain.TotalMRR{
					New:          []float32{100},
					Old:          []float32{0},
					Reactivation: []float32{0},
					Expansion:    []float32{0},
					Contraction:  []float32{0},
					Churn:        []float32{0},
					Total:        []float32{100}},
				err: errors.New("failed to set mrr to cache, error is: error while setting mrr to cache"),
			},
		},
		{
			input: testInput{
				userID:      "userGood",
				fileID:      "file",
				periodStart: "2021-10-01",
				periodEnd:   "2021-10-31",
			},
			want: testWant{
				months: []string{"10.2021"},
				mrr: domain.TotalMRR{
					New:          []float32{100},
					Old:          []float32{0},
					Reactivation: []float32{0},
					Expansion:    []float32{0},
					Contraction:  []float32{0},
					Churn:        []float32{0},
					Total:        []float32{100}},
				err: nil,
			},
		},
	}

	storageMock := &storagerepo.StorageRepositoryMock{}
	cacheMock := &cacherepo.CacheRepositoryMock{}

	for _, test := range tests {
		months, mrr, err := createAnalytics(storageMock, cacheMock, test.input.userID, test.input.fileID, test.input.periodStart, test.input.periodEnd)
		assert.Equal(t, test.want.months, months)
		assert.Equal(t, test.want.mrr, mrr)
		assert.Equal(t, test.want.err, err)
	}
}

func TestConvertRawMRR(t *testing.T) {
	type testInput struct {
		mrr []domain.MRR
	}
	type testWant struct {
		mrr domain.TotalMRR
	}

	tests := []struct {
		input testInput
		want  testWant
	}{
		{
			input: testInput{
				mrr: []domain.MRR{
					{New: 200},
					{Old: 200},
					{Churn: -200},
					{Reactivation: 200},
					{Expansion: 40},
					{Contraction: -40},
				},
			},
			want: testWant{
				mrr: domain.TotalMRR{
					New:          []float32{200, 0, 0, 0, 0, 0},
					Old:          []float32{0, 200, 0, 0, 0, 0},
					Reactivation: []float32{0, 0, 0, 200, 0, 0},
					Expansion:    []float32{0, 0, 0, 0, 40, 0},
					Contraction:  []float32{0, 0, 0, 0, 0, -40},
					Churn:        []float32{0, 0, -200, 0, 0, 0},
					Total:        []float32{200, 200, -200, 200, 40, -40},
				},
			},
		},
	}

	for _, test := range tests {
		mrr := convertRawMRR(test.input.mrr)
		assert.Equal(t, test.want.mrr, mrr)
	}
}

func TestCalculateTotalMRR(t *testing.T) {
	type testInput struct {
		mpp []domain.MPP
	}
	type testWant struct {
		mrr []domain.MRR
	}

	tests := []struct {
		input testInput
		want  testWant
	}{
		{
			input: testInput{
				mpp: []domain.MPP{
					{
						CustomerID: 0,
						Months:     []float32{100, 100, 0, 100, 120, 100},
					},
					{
						CustomerID: 1,
						Months:     []float32{100, 100, 0, 100, 120, 100},
					},
				},
			},
			want: testWant{
				mrr: []domain.MRR{
					{New: 200},
					{Old: 200},
					{Churn: -200},
					{Reactivation: 200},
					{Expansion: 40},
					{Contraction: -40},
				},
			},
		},
	}

	for _, test := range tests {
		mrr := calculateTotalMRR(test.input.mpp)
		assert.Equal(t, test.want.mrr, mrr)
	}
}

func TestCalculateClientMRR(t *testing.T) {
	type testInput struct {
		mpp domain.MPP
	}
	type testWant struct {
		mrr []domain.MRR
	}

	tests := []struct {
		input testInput
		want  testWant
	}{
		{
			input: testInput{
				mpp: domain.MPP{
					CustomerID: 0,
					Months:     []float32{100, 100, 0, 100, 120, 100},
				},
			},
			want: testWant{
				mrr: []domain.MRR{
					{New: 100},
					{Old: 100},
					{Churn: -100},
					{Reactivation: 100},
					{Expansion: 20},
					{Contraction: -20},
				},
			},
		},
	}

	for _, test := range tests {
		mrr := calculateClientMRR(test.input.mpp)
		assert.Equal(t, test.want.mrr, mrr)
	}
}

func TestFormMPP(t *testing.T) {
	type testInput struct {
		months                 []string
		userID, fileID         string
		periodStart, periodEnd time.Time
	}
	type testWant struct {
		mpp []domain.MPP
		err error
	}

	tests := []struct {
		input testInput
		want  testWant
	}{
		{
			input: testInput{
				months:      []string{"10.2021"},
				userID:      "errorGetInvoicesByPeriod",
				fileID:      "",
				periodStart: time.Now(),
				periodEnd:   time.Now(),
			},
			want: testWant{
				mpp: nil,
				err: errors.New("failed to get invoices from storage, error is: error while getting invoices by period"),
			},
		},
		{
			input: testInput{
				months:      []string{"10.2021"},
				userID:      "emptyGetInvoicesByPeriod",
				fileID:      "",
				periodStart: time.Now(),
				periodEnd:   time.Now(),
			},
			want: testWant{
				mpp: nil,
				err: errors.New("no data found for given period"),
			},
		},
		{
			input: testInput{
				months:      []string{"10.2021"},
				userID:      "userID",
				fileID:      "",
				periodStart: time.Date(2021, time.October, 1, 0, 0, 0, 0, time.UTC),
				periodEnd:   time.Now(),
			},
			want: testWant{
				mpp: []domain.MPP{
					{
						CustomerID: 0,
						Months:     []float32{100.0},
					},
				},
				err: nil,
			},
		},
	}

	storageMock := &storagerepo.StorageRepositoryMock{}

	for _, test := range tests {
		mpp, err := formMPP(storageMock, test.input.months, test.input.userID, test.input.fileID, test.input.periodStart, test.input.periodEnd)
		assert.Equal(t, test.want.mpp, mpp)
		assert.Equal(t, test.want.err, err)
	}
}

func TestFixMPP(t *testing.T) {
	type testInput struct {
		mpp []domain.MPP
	}
	type testWant struct {
		mpp []domain.MPP
	}

	tests := []struct {
		input testInput
		want  testWant
	}{
		{
			input: testInput{
				mpp: []domain.MPP{
					{
						CustomerID: 0,
						Months:     []float32{0, 0, 100.0},
					},
					{
						CustomerID: 0,
						Months:     []float32{100.0, 0, 0},
					},
				},
			},
			want: testWant{
				mpp: []domain.MPP{
					{
						CustomerID: 0,
						Months:     []float32{100.0, 0, 100.0},
					},
				},
			},
		},
	}

	for _, test := range tests {
		mpp := fixMPP(test.input.mpp)
		assert.Equal(t, test.want.mpp, mpp)
	}
}

func TestFormMPPEntries(t *testing.T) {
	type testInput struct {
		invoices    []domain.Invoice
		monthsCount int
		periodStart time.Time
	}
	type testWant struct {
		mpp []domain.MPP
	}

	tests := []struct {
		input testInput
		want  testWant
	}{
		{
			input: testInput{
				invoices: []domain.Invoice{
					{
						UserID:      "",
						FileID:      "",
						CustomerID:  0,
						PeriodStart: "2021-10-01",
						PaidPlan:    "monthly",
						PaidAmount:  100.0,
						PeriodEnd:   "2021-10-31",
					},
				},
				monthsCount: 1,
				periodStart: time.Date(2021, time.October, 1, 0, 0, 0, 0, time.UTC),
			},
			want: testWant{
				mpp: []domain.MPP{
					{
						CustomerID: 0,
						Months:     []float32{100.0},
					},
				},
			},
		},
		{
			input: testInput{
				invoices: []domain.Invoice{
					{
						UserID:      "",
						FileID:      "",
						CustomerID:  0,
						PeriodStart: "2021-09-01",
						PaidPlan:    "monthly",
						PaidAmount:  100.0,
						PeriodEnd:   "2021-09-31",
					},
				},
				monthsCount: 1,
				periodStart: time.Date(2021, time.October, 1, 0, 0, 0, 0, time.UTC),
			},
			want: testWant{
				mpp: []domain.MPP{
					{
						CustomerID: 0,
						Months:     []float32{0.0},
					},
				},
			},
		},
		{
			input: testInput{
				invoices: []domain.Invoice{
					{
						UserID:      "",
						FileID:      "",
						CustomerID:  0,
						PeriodStart: "2021-09-01",
						PaidPlan:    "annually",
						PaidAmount:  60.0,
						PeriodEnd:   "2022-09-31",
					},
				},
				monthsCount: 3,
				periodStart: time.Date(2021, time.October, 1, 0, 0, 0, 0, time.UTC),
			},
			want: testWant{
				mpp: []domain.MPP{
					{
						CustomerID: 0,
						Months:     []float32{5.0, 5.0, 5.0},
					},
				},
			},
		},
	}

	for _, test := range tests {
		mpp := formMPPEntries(test.input.invoices, test.input.monthsCount, test.input.periodStart)
		assert.Equal(t, test.want.mpp, mpp)
	}
}

func TestGetMonthsBetween(t *testing.T) {
	type testInput struct {
		fdate, sdate time.Time
	}
	type testWant struct {
		months []string
	}

	tests := []struct {
		input testInput
		want  testWant
	}{
		{
			input: testInput{
				fdate: timeInLocation(time.Now(), "Europe/Moscow"),
				sdate: timeInLocation(time.Now(), "Europe/Monaco"),
			},
			want: testWant{
				months: []string{"10.2021"},
			},
		},
		{
			input: testInput{
				fdate: time.Now().AddDate(0, 2, 0),
				sdate: time.Now().AddDate(0, 0, 0),
			},
			want: testWant{
				months: []string{"10.2021", "11.2021", "12.2021"},
			},
		},
		{
			input: testInput{
				fdate: time.Now().AddDate(0, 3, 0),
				sdate: time.Now().AddDate(0, 0, 0),
			},
			want: testWant{
				months: []string{"10.2021", "11.2021", "12.2021", "1.2022"},
			},
		},
	}

	for _, test := range tests {
		months := getMonthsBetween(test.input.fdate, test.input.sdate)
		assert.Equal(t, test.want.months, months)
	}
}

func TestGetMonthsDiff(t *testing.T) {
	type testInput struct {
		fdate, sdate time.Time
	}
	type testWant struct {
		count int
	}

	tests := []struct {
		input testInput
		want  testWant
	}{
		{
			input: testInput{
				fdate: timeInLocation(time.Now(), "Europe/Moscow"),
				sdate: timeInLocation(time.Now(), "Europe/Monaco"),
			},
			want: testWant{
				count: 0,
			},
		},
		{
			input: testInput{
				fdate: time.Now().AddDate(1, 0, 0),
				sdate: time.Now().AddDate(0, 0, -1),
			},
			want: testWant{
				count: 12,
			},
		},
	}

	for _, test := range tests {
		count := getMonthsDiff(test.input.fdate, test.input.sdate)
		assert.Equal(t, test.want.count, count)
	}
}

func timeInLocation(tm time.Time, loc string) time.Time {
	loadedLoc, _ := time.LoadLocation(loc)
	return tm.In(loadedLoc)
}
