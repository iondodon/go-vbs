package booking

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/iondodon/go-vbs/domain"
	"github.com/iondodon/go-vbs/repository"
)

//gobok:builder
type BookingRepository struct {
	queries *repository.Queries
}

func (repo *BookingRepository) Save(ctx context.Context, tx *sql.Tx, b *domain.Booking) error {
	if err := repo.queries.WithTx(tx).InsertNewBooking(ctx, repository.InsertNewBookingParams{
		Uuid:       b.UUID,
		VehicleID:  b.Vehicle.ID,
		CustomerID: b.Customer.ID,
	}); err != nil {
		return err
	}

	return nil
}

func (repo *BookingRepository) GetAll(ctx context.Context) ([]domain.Booking, error) {
	bookingsRows, err := repo.queries.SelectAllBookings(ctx)
	if err != nil {
		return nil, err
	}

	bookings := []domain.Booking{}
	for _, booking_row := range bookingsRows {
		booking := domain.Booking{}
		booking.ID = booking_row.ID.(int64)
		uuidStr := booking_row.Uuid.(string)
		parsedUUID, err := uuid.Parse(uuidStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse UUID: %w", err)
		}
		booking.UUID = parsedUUID

		booking.Vehicle = &domain.Vehicle{}
		booking.Vehicle.ID = booking_row.VehicleID

		booking.Customer = &domain.Customer{}
		booking.Customer.ID = booking_row.CustomerID

		bookings = append(bookings, booking)
	}

	return bookings, nil
}
