package bookingdate

import (
	"time"

	"github.com/iondodon/go-vbs/domain"
	"github.com/iondodon/go-vbs/integration"
)

const findAllInPeriodInclusive = `
	SELECT bd.id, bd.time
	FROM booking_date bd
	WHERE bd.time >= ? AND bd.time <= ?
`

type BookingDateRepository interface {
	FindAllInPeriodInclusive(from, to time.Time) ([]domain.BookingDate, error)
}

type bookingDateRepository struct {
	db integration.DB
}

func NewBookingDateRepository(db integration.DB) BookingDateRepository {
	return &bookingDateRepository{db: db}
}

func (repo *bookingDateRepository) FindAllInPeriodInclusive(from, to time.Time) ([]domain.BookingDate, error) {
	var bookingDates []domain.BookingDate

	rows, err := repo.db.Query(findAllInPeriodInclusive, from, to)
	if err != nil {
		defer rows.Close()
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		bd := domain.BookingDate{}

		err := rows.Scan(&bd.ID, &bd.Time)
		if err != nil {
			return nil, err
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return bookingDates, nil
}
