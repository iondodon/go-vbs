package bookingRepository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/iondodon/go-vbs/booking/business"
	"github.com/iondodon/go-vbs/domain"
	"github.com/iondodon/go-vbs/repository"
)

// Repository implements BookingRepository interface
type Repository struct {
	queries *repository.Queries
}

// Ensure Repository implements the business interface
var _ business.BookingRepository = (*Repository)(nil)

func New(queries *repository.Queries) *Repository {
	return &Repository{
		queries: queries,
	}
}

func (r *Repository) Save(ctx context.Context, tx *sql.Tx, b *domain.Booking) error {
	if err := r.queries.WithTx(tx).InsertNewBooking(ctx, repository.InsertNewBookingParams{
		Uuid:       b.UUID,
		VehicleID:  b.Vehicle.ID,
		CustomerID: b.Customer.ID,
	}); err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetAll(ctx context.Context) ([]domain.Booking, error) {
	bookingsRows, err := r.queries.SelectAllBookings(ctx)
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
