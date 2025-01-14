package booking

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/iondodon/go-vbs/domain"
	"github.com/iondodon/go-vbs/repository"
)

type BookingRepository interface {
	Save(context.Context, *sql.Tx, *domain.Booking) error
	GetAll(context.Context) ([]domain.Booking, error)
}

type bookingRepository struct {
	queries *repository.Queries
}

func NewBookingRepository(queries *repository.Queries) BookingRepository {
	return &bookingRepository{queries: queries}
}

func (repo *bookingRepository) Save(ctx context.Context, tx *sql.Tx, b *domain.Booking) error {
	if err := repo.queries.WithTx(tx).InsertNewBooking(ctx, repository.InsertNewBookingParams{
		Uuid:       b.UUID,
		VehicleID:  b.Vehicle.ID,
		CustomerID: b.Customer.ID,
	}); err != nil {
		return err
	}

	return nil
}

func (repo *bookingRepository) GetAll(ctx context.Context) ([]domain.Booking, error) {
	bookingsRows, err := repo.queries.SelectAllBookings(ctx)
	if err != nil {
		return nil, err
	}

	bookings := []domain.Booking{}
	for _, booking_row := range bookingsRows {
		booking := domain.Booking{}
		booking.ID = booking_row.ID.(int64)
		booking.UUID = booking_row.Uuid.(uuid.UUID)

		booking.Vehicle = &domain.Vehicle{}
		booking.Vehicle.ID = booking_row.VehicleID

		booking.Customer = &domain.Customer{}
		booking.Customer.ID = booking_row.CustomerID

		bookings = append(bookings, booking)
	}

	return bookings, nil
}
