package business

import (
	"context"

	"github.com/google/uuid"
	"github.com/iondodon/go-vbs/customer/domain"
)

// CustomerRepository defines what the customer business logic needs from data access
type CustomerRepository interface {
	FindByUUID(ctx context.Context, cUUID uuid.UUID) (*domain.Customer, error)
}
