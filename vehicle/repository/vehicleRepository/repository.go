package vehicleRepository

import (
	"context"
	"fmt"

	uuidLib "github.com/google/uuid"
	"github.com/iondodon/go-vbs/booking/controller"
	"github.com/iondodon/go-vbs/repository"
	"github.com/iondodon/go-vbs/vehicle/business"
	"github.com/iondodon/go-vbs/vehicle/domain"

	_ "github.com/mattn/go-sqlite3"
)

type Repository struct {
	queries *repository.Queries
}

// Ensure Repository implements the business interface
var _ business.VehicleRepository = (*Repository)(nil)

func New(queries *repository.Queries) *Repository {
	return &Repository{
		queries: queries,
	}
}

func (r *Repository) FindByUUID(ctx context.Context, vUUID uuidLib.UUID) (*domain.Vehicle, error) {
	vehicleRow, err := r.queries.GetVehicleByUUID(ctx, vUUID)
	if err != nil {
		return nil, err
	}

	var vehicle domain.Vehicle
	vehicle.ID = vehicleRow.ID.(int64)
	uuidStr := vehicleRow.Uuid.(string)
	parsedUUID, err := uuidLib.Parse(uuidStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse UUID: %w", err)
	}
	vehicle.UUID = parsedUUID
	vehicle.RegistrationNumber = vehicleRow.RegistrationNumber
	vehicle.Make = vehicleRow.Make
	vehicle.Model = vehicleRow.Model
	vehicle.FuelType = domain.FuelType(vehicleRow.FuelType)

	var vehCat domain.VehicleCategory
	vehCat.ID = vehicleRow.ID_2.(int64)
	vehCat.PricePerDay = float32(vehicleRow.PricePerDay)
	vehCat.VehicleType = domain.VehicleType(vehicleRow.Category)
	vehicle.VehicleCategory = &vehCat

	return &vehicle, nil
}

func (r *Repository) VehicleHasBookedDatesOnPeriod(ctx context.Context, vUUID uuidLib.UUID, period controller.DatePeriodDTO) (bool, error) {
	res, err := r.queries.VehicleHasBookedDatesOnPeriod(ctx, repository.VehicleHasBookedDatesOnPeriodParams{
		Uuid:   vUUID,
		Time:   period.FromDate,
		Time_2: period.ToDate,
	})
	if err != nil {
		return false, err
	}

	if res == 1 {
		return true, nil
	}

	return false, nil
}
