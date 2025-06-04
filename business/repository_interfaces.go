// Package business defines the business logic layer interfaces.
// This file contains repository interfaces that define what the business logic
// needs from the infrastructure layer. Following the Dependency Inversion Principle,
// the business logic defines these contracts and the infrastructure implements them.
package business

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/iondodon/go-vbs/domain"
	"github.com/iondodon/go-vbs/dto"
)

// VehicleRepository defines what the business logic needs from vehicle data access
type VehicleRepository interface {
	FindByUUID(ctx context.Context, vUUID uuid.UUID) (*domain.Vehicle, error)
	VehicleHasBookedDatesOnPeriod(ctx context.Context, vUUID uuid.UUID, period dto.DatePeriodDTO) (bool, error)
}

// BookingRepository defines what the business logic needs from booking data access
type BookingRepository interface {
	Save(ctx context.Context, tx *sql.Tx, b *domain.Booking) error
	GetAll(ctx context.Context) ([]domain.Booking, error)
}

// CustomerRepository defines what the business logic needs from customer data access
type CustomerRepository interface {
	FindByUUID(ctx context.Context, cUUID uuid.UUID) (*domain.Customer, error)
}

// BookingDateRepository defines what the business logic needs from booking date data access
type BookingDateRepository interface {
	FindAllInPeriodInclusive(ctx context.Context, from, to time.Time) ([]*domain.BookingDate, error)
	Save(ctx context.Context, tx *sql.Tx, bd *domain.BookingDate) error
}
