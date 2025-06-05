package bookingBusiness

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/iondodon/go-vbs/booking/bookingController"
	"github.com/iondodon/go-vbs/booking/bookingDomain"
	"github.com/iondodon/go-vbs/customer/customerDomain"
	"github.com/iondodon/go-vbs/vehicle/vehicleDomain"
)

// BookingRepository defines what the booking business logic needs from booking data access
type BookingRepository interface {
	Save(ctx context.Context, tx *sql.Tx, b *bookingDomain.Booking) error
	GetAll(ctx context.Context) ([]bookingDomain.Booking, error)
}

// BookingDateRepository defines what the booking business logic needs from booking date data access
type BookingDateRepository interface {
	FindAllInPeriodInclusive(ctx context.Context, from, to time.Time) ([]*bookingDomain.BookingDate, error)
	Save(ctx context.Context, tx *sql.Tx, bd *bookingDomain.BookingDate) error
}

// Cross-domain dependencies (defined here since booking consumes them)
type VehicleRepository interface {
	FindByUUID(ctx context.Context, vUUID uuid.UUID) (*vehicleDomain.Vehicle, error)
	VehicleHasBookedDatesOnPeriod(ctx context.Context, vUUID uuid.UUID, period bookingController.DatePeriodDTO) (bool, error)
}

type CustomerRepository interface {
	FindByUUID(ctx context.Context, cUUID uuid.UUID) (*customerDomain.Customer, error)
}
