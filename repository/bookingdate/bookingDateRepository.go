package bookingdate

import (
	"context"
	"database/sql"
	"time"

	"github.com/iondodon/go-vbs/domain"
	"github.com/iondodon/go-vbs/repository"
)

type BookingDateRepository interface {
	FindAllInPeriodInclusive(ctx context.Context, from, to time.Time) ([]*domain.BookingDate, error)
	Save(ctx context.Context, tx *sql.Tx, bd *domain.BookingDate) error
}

type bookingDateRepository struct {
	queries *repository.Queries
}

func NewBookingDateRepository(queries *repository.Queries) BookingDateRepository {
	return &bookingDateRepository{queries: queries}
}

func (repo *bookingDateRepository) FindAllInPeriodInclusive(ctx context.Context, from, to time.Time) ([]*domain.BookingDate, error) {
	var bookingDates []*domain.BookingDate

	booking_date_rows, err := repo.queries.FindAllInPeriodInclusive(ctx, repository.FindAllInPeriodInclusiveParams{
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

func (repo *bookingDateRepository) Save(ctx context.Context, tx *sql.Tx, bd *domain.BookingDate) error {
	if err := repo.queries.WithTx(tx).SaveNewBookingDate(ctx, bd.Time); err != nil {
		return err
	}
	return nil
}
