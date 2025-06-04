package business

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	bookingIn "github.com/iondodon/go-vbs/booking/in"
	"github.com/iondodon/go-vbs/domain"
)

// BookingRepository defines what the booking business logic needs from booking data access
type BookingRepository interface {
	Save(ctx context.Context, tx *sql.Tx, b *domain.Booking) error
	GetAll(ctx context.Context) ([]domain.Booking, error)
}

// BookingDateRepository defines what the booking business logic needs from booking date data access
type BookingDateRepository interface {
	FindAllInPeriodInclusive(ctx context.Context, from, to time.Time) ([]*domain.BookingDate, error)
	Save(ctx context.Context, tx *sql.Tx, bd *domain.BookingDate) error
}

// Cross-domain dependencies (defined here since booking consumes them)
type VehicleRepository interface {
	FindByUUID(ctx context.Context, vUUID uuid.UUID) (*domain.Vehicle, error)
	VehicleHasBookedDatesOnPeriod(ctx context.Context, vUUID uuid.UUID, period bookingIn.DatePeriodDTO) (bool, error)
}

type CustomerRepository interface {
	FindByUUID(ctx context.Context, cUUID uuid.UUID) (*domain.Customer, error)
}

type VehicleAvailabilityService interface {
	CheckForPeriod(ctx context.Context, vUUID uuid.UUID, period bookingIn.DatePeriodDTO) (bool, error)
}
