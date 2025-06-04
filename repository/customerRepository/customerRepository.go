package customerRepository

import (
	"context"

	uuidlib "github.com/google/uuid"
	"github.com/iondodon/go-vbs/domain"
	"github.com/iondodon/go-vbs/repository"
)

type CustomerRepository struct {
	queries *repository.Queries
}

func (repo *CustomerRepository) FindByUUID(ctx context.Context, cUUID uuidlib.UUID) (*domain.Customer, error) {
	var customer domain.Customer

	customerRow, err := repo.queries.GetCustomerByUUID(ctx, cUUID)
	if err != nil {
		return nil, err
	}

	customer.ID = customerRow.ID.(int64)

	customerUuidStr := customerRow.Uuid.(string)
	uuid, err := uuidlib.Parse(customerUuidStr)
	if err != nil {
		return nil, err
	}
	customer.UUID = uuid

	return &customer, nil
}
