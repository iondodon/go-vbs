package booking

import (
	"github.com/iondodon/go-vbs/domain"
	"github.com/iondodon/go-vbs/integration"
)

const insertNewBookingSQL = `
	INSERT INTO booking(uuid, vehicle_id, customer_id)
	VALUES (?, ?, ?)
`

type BookingRepository interface {
	Save(b *domain.Booking) error
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
