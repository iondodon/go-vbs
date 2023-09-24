package booking

import (
	"github.com/iondodon/go-vbs/domain"
	"github.com/iondodon/go-vbs/integration"
)

const insertNewBookingSQL = `
	INSERT INTO booking(uuid, vehicle_id, customer_id)
	VALUES (?, ?, ?)
`

const selectAllBookings = `
	SELECT b.id, b.uuid, b.vehicle_id, b.customer_id 
	FROM booking b
`

type BookingRepository interface {
	Save(b *domain.Booking) error
	GetAll() ([]domain.Booking, error)
}

type bookingRepository struct {
	db integration.DB
}

func NewBookingRepository(db integration.DB) BookingRepository {
	return &bookingRepository{db: db}
}

func (repo *bookingRepository) Save(b *domain.Booking) error {
	_, err := repo.db.Exec(insertNewBookingSQL, b.UUID, b.Vehicle.ID, b.Customer.ID)
	if err != nil {
		return err
	}
	return nil
}

func (repo *bookingRepository) GetAll() ([]domain.Booking, error) {
	rows, err := repo.db.Query(selectAllBookings)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bookings := []domain.Booking{}
	for rows.Next() {
		booking := domain.Booking{}
		booking.Vehicle = &domain.Vehicle{}
		booking.Customer = &domain.Customer{}

		err := rows.Scan(&booking.ID, &booking.UUID, &booking.Vehicle.ID, &booking.Customer.ID)
		if err != nil {
			return nil, err
		}

		bookings = append(bookings, booking)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return bookings, nil
}
