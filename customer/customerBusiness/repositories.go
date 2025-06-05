package customerBusiness

import (
	"context"

	"github.com/google/uuid"
	"github.com/iondodon/go-vbs/customer/customerDomain"
)

// CustomerRepository defines what the customer business logic needs from data access
type CustomerRepository interface {
	FindByUUID(ctx context.Context, cUUID uuid.UUID) (*customerDomain.Customer, error)
}
