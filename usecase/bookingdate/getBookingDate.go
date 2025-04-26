package bookingdate

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/iondodon/go-vbs/domain"
	"github.com/iondodon/go-vbs/dto"
	bdRepoPkg "github.com/iondodon/go-vbs/repository/bookingdate"
)

//gobok:builder
type GetBookingDates struct {
	bdRepo bdRepoPkg.BookingDateRepository
}

func (uc *GetBookingDates) ForPeriod(ctx context.Context, tx *sql.Tx, customerUUID, vehicleUUID uuid.UUID, period dto.DatePeriodDTO) ([]*domain.BookingDate, error) {
	persistedBookingDates, err := uc.bdRepo.FindAllInPeriodInclusive(ctx, period.FromDate, period.ToDate)
	if err != nil {
		return nil, err
	}

	var dates []time.Time
	for _, bd := range persistedBookingDates {
		bd.Time = removeTimePart(bd.Time)
		dates = append(dates, bd.Time)
	}

	fromDate := removeTimePart(period.FromDate)
	toDate := removeTimePart(period.ToDate)
	for date := fromDate; date.Before(toDate) || date.Equal(toDate); date = date.Add(time.Hour + 24) {
		if !contains(dates, date) {
			bd := domain.BookingDate{Time: date}
			uc.bdRepo.Save(ctx, tx, &bd)
			persistedBookingDates = append(persistedBookingDates, &bd)
		}
	}

	return persistedBookingDates, nil
}

func removeTimePart(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

func contains(dates []time.Time, targetDate time.Time) bool {
	for _, d := range dates {
		if d.Equal(targetDate) {
			return true
		}
	}
	return false
}
