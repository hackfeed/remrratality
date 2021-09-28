package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hackfeed/remrratality/backend/internal/domain"
	"github.com/hackfeed/remrratality/backend/internal/server/models"
	cacherepo "github.com/hackfeed/remrratality/backend/internal/store/cache_repo"
	storagerepo "github.com/hackfeed/remrratality/backend/internal/store/storage_repo"
)

var (
	layout = "2006-01-02"
)

// GetAnalytics godoc
// @Summary Get MRR analytics data
// @Description Getting MRR analytics data with all components for given period
// @Tags analytics
// @Accept  json
// @Produce  json
// @Success 200 {object} models.ResponseSuccessAnalytics
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Security ApiKeyAuth
// @Param request body models.Period true "Parameters for MRR analytics"
// @Router /analytics/mrr [post]
func GetAnalytics(c *gin.Context) {
	userID, ok := c.MustGet("user_id").(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
			Message: "Unable to determine logged in user",
		})
		return
	}
	storageRepo, ok := c.MustGet("storage_repo").(storagerepo.StorageRepository)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
			Message: "Failed to get storage_repo",
		})
		return
	}
	cacheRepo, ok := c.MustGet("cache_repo").(cacherepo.CacheRepository)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
			Message: "Failed to get cache_repo",
		})
		return
	}

	var req models.Period

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.Response{
			Message: "Failed to parse request body",
		})
		return
	}

	months, mrr, err := getAnalytics(storageRepo, cacheRepo, userID, req.Filename, req.PeriodStart, req.PeriodEnd)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.ResponseSuccessAnalytics{
		Message: "Analytics is loaded",
		Months:  months,
		MRR:     mrr,
	})
}

func getAnalytics(storageRepo storagerepo.StorageRepository, cacheRepo cacherepo.CacheRepository, userID, fileID, periodStart, periodEnd string) ([]string, domain.TotalMRR, error) {
	var (
		mrr    domain.TotalMRR
		months []string
	)

	periodStartDate, _ := time.Parse(layout, periodStart)
	periodEndDate, _ := time.Parse(layout, periodEnd)

	if periodStartDate.After(periodEndDate) {
		return months, mrr, errors.New("Period start should be less than period end")
	}

	userFilePeriod := fmt.Sprintf("%s.%s-%s-%s", userID, fileID, periodStart, periodEnd)
	mrr, err := cacheRepo.GetMRR(userFilePeriod)
	if err != nil {
		return months, mrr, err
	}

	months = getMonthsBetween(periodEndDate, periodStartDate)
	if len(mrr.Total) != 0 {
		return months, mrr, nil
	}

	formedMPP, err := formMPP(storageRepo, months, userID, fileID, periodStartDate, periodEndDate)
	if err != nil {
		return months, mrr, err
	}
	mrr = convertRawMRR(calculateTotalMRR(formedMPP))
	_, err = cacheRepo.SetMRR(userFilePeriod, mrr)

	return months, mrr, err
}

func convertRawMRR(rawMRR []domain.MRR) domain.TotalMRR {
	var totalMRR domain.TotalMRR

	for _, mrr := range rawMRR {
		totalMRR.New = append(totalMRR.New, mrr.New)
		totalMRR.Old = append(totalMRR.Old, mrr.Old)
		totalMRR.Reactivation = append(totalMRR.Reactivation, mrr.Reactivation)
		totalMRR.Expansion = append(totalMRR.Expansion, mrr.Expansion)
		totalMRR.Contraction = append(totalMRR.Contraction, mrr.Contraction)
		totalMRR.Churn = append(totalMRR.Churn, mrr.Churn)
		totalMRR.Total = append(totalMRR.Total, mrr.New+mrr.Old+mrr.Reactivation+mrr.Expansion+mrr.Contraction+mrr.Churn)
	}

	return totalMRR
}

func calculateTotalMRR(mpp []domain.MPP) []domain.MRR {
	monthsCount := len(mpp[0].Months)
	totalMRR := make([]domain.MRR, monthsCount)

	for _, mppEntry := range mpp {
		clientMRR := calculateClientMRR(mppEntry)

		for i := 0; i < monthsCount; i++ {
			totalMRR[i].New += clientMRR[i].New
			totalMRR[i].Old += clientMRR[i].Old
			totalMRR[i].Reactivation += clientMRR[i].Reactivation
			totalMRR[i].Expansion += clientMRR[i].Expansion
			totalMRR[i].Contraction += clientMRR[i].Contraction
			totalMRR[i].Churn += clientMRR[i].Churn
		}
	}

	return totalMRR
}

func calculateClientMRR(mpp domain.MPP) []domain.MRR {
	monthsCount := len(mpp.Months)
	clientMRR := make([]domain.MRR, monthsCount)

	isNew := true

	for i := range mpp.Months {
		var monthMRR domain.MRR

		if mpp.Months[i] > 0 && isNew {
			monthMRR.New = mpp.Months[i]
			isNew = false
		} else if i > 0 && mpp.Months[i] == mpp.Months[i-1] {
			monthMRR.Old = mpp.Months[i]
		} else if i > 0 && mpp.Months[i-1] == 0 && mpp.Months[i] > 0 && !isNew {
			monthMRR.Reactivation = mpp.Months[i]
		} else if i > 0 && mpp.Months[i] > mpp.Months[i-1] && mpp.Months[i-1] != 0 {
			monthMRR.Expansion = mpp.Months[i] - mpp.Months[i-1]
		} else if i > 0 && mpp.Months[i] < mpp.Months[i-1] && mpp.Months[i] != 0 {
			monthMRR.Contraction = mpp.Months[i] - mpp.Months[i-1]
		} else if i > 0 && mpp.Months[i] == 0 && mpp.Months[i-1] > 0 {
			monthMRR.Churn = -mpp.Months[i-1]
		}

		clientMRR[i] = monthMRR
	}

	return clientMRR
}

func formMPP(storageRepo storagerepo.StorageRepository, months []string, userID, fileID string, periodStart, periodEnd time.Time) ([]domain.MPP, error) {
	fixedPeriodEnd := periodEnd.AddDate(0, 1, -1)

	invoices, err := storageRepo.GetInvoicesByPeriod(userID, fileID, periodStart, fixedPeriodEnd)
	if err != nil {
		return nil, err
	}

	if len(invoices) == 0 {
		return nil, errors.New("No data found for given period")
	}

	mpp := formMPPEntries(invoices, len(months), periodStart)
	fixedMPP := fixMPP(mpp)

	return fixedMPP, err
}

func fixMPP(mpp []domain.MPP) []domain.MPP {
	customerMap := make(map[uint32][]float32)

	for _, mppEntry := range mpp {
		if _, ok := customerMap[mppEntry.CustomerID]; !ok {
			customerMap[mppEntry.CustomerID] = mppEntry.Months
		} else {
			for i := range mppEntry.Months {
				customerMap[mppEntry.CustomerID][i] += mppEntry.Months[i]
			}
		}
	}

	fixedMPP := make([]domain.MPP, 0)
	for customerID, months := range customerMap {
		fixedMPP = append(fixedMPP, domain.MPP{
			CustomerID: customerID,
			Months:     months,
		})
	}

	return fixedMPP
}

func formMPPEntries(invoices []domain.Invoice, monthsCount int, periodStart time.Time) []domain.MPP {
	invoicesCount := len(invoices)
	mppEntries := make([]domain.MPP, invoicesCount)

	for i, invoice := range invoices {
		moneyPerMonth := make([]float32, monthsCount)
		paidAmount := invoice.PaidAmount
		periodLen := 1

		invoicePeriodStart, _ := time.Parse(layout, invoice.PeriodStart)
		startMonth := getMonthsDiff(invoicePeriodStart, periodStart)
		if startMonth < 0 {
			periodLen += startMonth
			startMonth = 0
		}

		if invoice.PaidPlan == "annually" {
			paidAmount /= 12
			periodLen = 12
		}

		for j := startMonth; j < monthsCount; j++ {
			if periodLen <= 0 {
				paidAmount = 0
			}
			moneyPerMonth[j] = paidAmount
			periodLen--
		}

		moneyFlow := domain.MPP{
			CustomerID: invoice.CustomerID,
			Months:     moneyPerMonth,
		}

		mppEntries[i] = moneyFlow
	}

	return mppEntries
}

func getMonthsBetween(fdate, sdate time.Time) []string {
	if fdate.Location() != sdate.Location() {
		sdate = sdate.In(fdate.Location())
	}
	if fdate.After(sdate) {
		fdate, sdate = sdate, fdate
	}

	fyear, fmonth, _ := fdate.Date()
	syear, smonth, _ := sdate.Date()

	count := int(smonth - fmonth)
	count += 12*(syear-fyear) + 1

	var months []string

	for i := 0; i < count; i++ {
		yearsDif := (int(fmonth) + i) / 12
		month := (int(fmonth) + i) % 12
		if month == 0 {
			yearsDif -= 1
			month = 12
		}
		months = append(months, fmt.Sprintf("%v.%v", month, fyear+yearsDif))
	}

	return months
}

func getMonthsDiff(fdate, sdate time.Time) int {
	if fdate.Location() != sdate.Location() {
		sdate = sdate.In(fdate.Location())
	}

	fyear, fmonth, _ := fdate.Date()
	syear, smonth, _ := sdate.Date()

	if fyear > syear {
		fyear, syear = syear, fyear
	}

	count := int(fmonth - smonth)
	count += 12 * (syear - fyear)

	return count
}
