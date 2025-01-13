package bookingdate

import (
	"context"
	"time"

	"github.com/iondodon/go-vbs/domain"
	"github.com/iondodon/go-vbs/repository"
)

type BookingDateRepository interface {
	FindAllInPeriodInclusive(from, to time.Time) ([]*domain.BookingDate, error)
	Save(bd *domain.BookingDate) error
}

type bookingDateRepository struct {
	queries *repository.Queries
}

func NewBookingDateRepository(queries *repository.Queries) BookingDateRepository {
	return &bookingDateRepository{queries: queries}
}

func (repo *bookingDateRepository) FindAllInPeriodInclusive(from, to time.Time) ([]*domain.BookingDate, error) {
	var bookingDates []*domain.BookingDate

	booking_date_rows, err := repo.queries.FindAllInPeriodInclusive(context.Background(), repository.FindAllInPeriodInclusiveParams{
		Time:   from,
		Time_2: to,
	})
	if err != nil {
		return nil, err
	}

	for _, booking_date_row := range booking_date_rows {
		bd := domain.BookingDate{}

		bd.ID = booking_date_row.ID.(int64)
		bd.Time = booking_date_row.Time

		bookingDates = append(bookingDates, &bd)
	}

	return bookingDates, nil
}

func (repo *bookingDateRepository) Save(bd *domain.BookingDate) error {
	err := repo.queries.SaveNewBookingDate(context.Background(), bd.Time)
	if err != nil {
		return err
	}
	return nil
}
