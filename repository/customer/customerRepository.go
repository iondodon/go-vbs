package customer

import (
	"context"

	uuidlib "github.com/google/uuid"
	"github.com/iondodon/go-vbs/domain"
	"github.com/iondodon/go-vbs/repository"
)

type CustomerRepository interface {
	FindByUUID(cUUID uuidlib.UUID) (*domain.Customer, error)
}

type customerRepository struct {
	queries *repository.Queries
}

func NewCustomerRepository(queries *repository.Queries) CustomerRepository {
	return &customerRepository{queries: queries}
}

func (repo *customerRepository) FindByUUID(cUUID uuidlib.UUID) (*domain.Customer, error) {
	var customer domain.Customer

	customerRow, err := repo.queries.GetCustomerByUUID(context.Background(), cUUID)
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
