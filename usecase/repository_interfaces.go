// Package usecase defines the business logic layer interfaces.
// This file contains repository interfaces that define what the business logic
// needs from the infrastructure layer. Following the Dependency Inversion Principle,
// the business logic defines these contracts and the infrastructure implements them.
package usecase

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/iondodon/go-vbs/domain"
	"github.com/iondodon/go-vbs/dto"
)

// VehicleRepositoryInterface defines what the business logic needs from vehicle data access
type VehicleRepositoryInterface interface {
	FindByUUID(ctx context.Context, vUUID uuid.UUID) (*domain.Vehicle, error)
	VehicleHasBookedDatesOnPeriod(ctx context.Context, vUUID uuid.UUID, period dto.DatePeriodDTO) (bool, error)
}

// BookingRepositoryInterface defines what the business logic needs from booking data access
type BookingRepositoryInterface interface {
	Save(ctx context.Context, tx *sql.Tx, b *domain.Booking) error
	GetAll(ctx context.Context) ([]domain.Booking, error)
}

// CustomerRepositoryInterface defines what the business logic needs from customer data access
type CustomerRepositoryInterface interface {
	FindByUUID(ctx context.Context, cUUID uuid.UUID) (*domain.Customer, error)
}

// BookingDateRepositoryInterface defines what the business logic needs from booking date data access
type BookingDateRepositoryInterface interface {
	FindAllInPeriodInclusive(ctx context.Context, from, to time.Time) ([]*domain.BookingDate, error)
	Save(ctx context.Context, tx *sql.Tx, bd *domain.BookingDate) error
}
