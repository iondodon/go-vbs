package services

import (
	"context"

	"github.com/google/uuid"
	customerDomain "github.com/iondodon/go-vbs/internal/customer/domain"
)

// CustomerRepository defines what the customer services layer needs from data access
type CustomerRepository interface {
	FindByUUID(ctx context.Context, cUUID uuid.UUID) (*customerDomain.Customer, error)
}
