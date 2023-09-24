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

const saveNewBookingDate = `
	INSERT INTO booking_dates(time)
	VALUES (?)
`

type BookingDateRepository interface {
	FindAllInPeriodInclusive(from, to time.Time) ([]*domain.BookingDate, error)
	Save(bd *domain.BookingDate) error
}

type bookingDateRepository struct {
	db integration.DB
}

func NewBookingDateRepository(db integration.DB) BookingDateRepository {
	return &bookingDateRepository{db: db}
}

func (repo *bookingDateRepository) FindAllInPeriodInclusive(from, to time.Time) ([]*domain.BookingDate, error) {
	var bookingDates []*domain.BookingDate

	rows, err := repo.db.Query(findAllInPeriodInclusive, from, to)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		bd := domain.BookingDate{}

		err := rows.Scan(&bd.ID, &bd.Time)
		if err != nil {
			return nil, err
		}

		bookingDates = append(bookingDates, &bd)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return bookingDates, nil
}

func (repo *bookingDateRepository) Save(bd *domain.BookingDate) error {
	_, err := repo.db.Exec(saveNewBookingDate, bd.Time)
	if err != nil {
		return err
	}
	return nil
}
