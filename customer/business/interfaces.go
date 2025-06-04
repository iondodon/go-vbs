package business

import (
	"context"

	"github.com/google/uuid"
	"github.com/iondodon/go-vbs/domain"
)

// Repository defines what the customer business logic needs from data access
type Repository interface {
	FindByUUID(ctx context.Context, cUUID uuid.UUID) (*domain.Customer, error)
}
