package vehicle

import (
	"context"

	"github.com/iondodon/go-vbs/domain"
	"github.com/iondodon/go-vbs/dto"
	"github.com/iondodon/go-vbs/repository"

	uuidLib "github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type VehicleRepository interface {
	FindByUUID(ctx context.Context, vUUID uuidLib.UUID) (*domain.Vehicle, error)
	VehicleHasBookedDatesOnPeriod(ctx context.Context, vUUID uuidLib.UUID, period dto.DatePeriodDTO) (bool, error)
}

type vehicleRepository struct {
	queries *repository.Queries
}

func NewVehicleRepository(queries *repository.Queries) VehicleRepository {
	return &vehicleRepository{queries: queries}
}

func (repo *vehicleRepository) FindByUUID(ctx context.Context, vUUID uuidLib.UUID) (*domain.Vehicle, error) {
	vehicleRow, err := repo.queries.GetVehicleByUUID(ctx, vUUID)
	if err != nil {
		return nil, err
	}

	var vehicle domain.Vehicle
	vehicle.ID = vehicleRow.ID.(int64)
	vehicle.UUID = vehicleRow.Uuid.(uuidLib.UUID)
	vehicle.RegistrationNumber = vehicleRow.RegistrationNumber
	vehicle.Make = vehicleRow.Make
	vehicle.Model = vehicleRow.Model
	vehicle.FuelType = domain.FuelType(vehicleRow.FuelType)

	var vehCat domain.VehicleCategory
	vehCat.ID = vehicleRow.ID_2.(int64)
	vehCat.PricePerDay = float32(vehicleRow.PricePerDay)
	vehCat.VehicleType = domain.VehicleType(vehicleRow.Category)
	vehicle.VehicleCategory = &vehCat

	vehicleBookingsRows, err := repo.queries.SelectBookingsByVehicleID(ctx, vehicle.ID)
	if err != nil {
		return nil, err
	}

	vehicle.Bookings = []*domain.Booking{}
	for _, vehicleBookingRow := range vehicleBookingsRows {
		var booking domain.Booking
		var customer domain.Customer

		booking.ID = vehicleBookingRow.ID.(int64)
		booking.UUID = vehicleBookingRow.Uuid.(uuidLib.UUID)
		booking.Vehicle = &vehicle
		booking.Customer = &customer

		vehicle.Bookings = append(vehicle.Bookings, &booking)
	}

	return &vehicle, nil
}

func (repo *vehicleRepository) VehicleHasBookedDatesOnPeriod(ctx context.Context, vUUID uuidLib.UUID, period dto.DatePeriodDTO) (bool, error) {
	res, err := repo.queries.VehicleHasBookedDatesOnPeriod(ctx, repository.VehicleHasBookedDatesOnPeriodParams{
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
