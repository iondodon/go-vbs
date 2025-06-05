package bookingDateRepository

import (
	"context"
	"database/sql"
	"time"

	"github.com/iondodon/go-vbs/booking/bookingBusiness"
	"github.com/iondodon/go-vbs/booking/bookingDomain"
	"github.com/iondodon/go-vbs/repository"

	_ "github.com/mattn/go-sqlite3"
)

// Repository implements BookingDateRepository interface
type Repository struct {
	queries *repository.Queries
}

// Ensure Repository implements the business interface
var _ bookingBusiness.BookingDateRepository = (*Repository)(nil)

func New(queries *repository.Queries) *Repository {
	return &Repository{
		queries: queries,
	}
}

func (r *Repository) FindAllInPeriodInclusive(ctx context.Context, from, to time.Time) ([]*bookingDomain.BookingDate, error) {
	var bookingDates []*bookingDomain.BookingDate

	booking_date_rows, err := r.queries.FindAllInPeriodInclusive(ctx, repository.FindAllInPeriodInclusiveParams{
		Time:   from,
		Time_2: to,
	})
	if err != nil {
		return nil, err
	}

	for _, booking_date_row := range booking_date_rows {
		bd := bookingDomain.BookingDate{}

		bd.ID = booking_date_row.ID.(int64)
		bd.Time = booking_date_row.Time

		bookingDates = append(bookingDates, &bd)
	}

	return bookingDates, nil
}

func (r *Repository) Save(ctx context.Context, tx *sql.Tx, bd *bookingDomain.BookingDate) error {
	if err := r.queries.WithTx(tx).SaveNewBookingDate(ctx, bd.Time); err != nil {
		return err
	}
	return nil
}
