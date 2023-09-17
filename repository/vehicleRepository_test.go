package repository

import (
	"testing"

	"github.com/iondodon/go-vbs/domain"

	"github.com/DATA-DOG/go-sqlmock"
	uuidLib "github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_vehicleRepository_FindByUUID(t *testing.T) {
	uuid, err := uuidLib.Parse("d06a5744-ce7d-4aa7-ba47-076cae095bb1")
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}

	mv := &domain.Vehicle{
		ID:                 123,
		UUID:               uuid,
		RegistrationNumber: "reg number",
		Make:               "Tesla",
		Model:              "X",
		FuelType:           "DIESEL",
		VehicleCategory: &domain.VehicleCategory{
			ID:          321,
			VehicleType: "SMALL_CAR",
			PricePerDay: 123.321,
		},
		Bookings: nil,
	}

	t.Run("the correct query is executed", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err.Error())
		}
		defer db.Close()

		var vrp VehicleRepository = &vehicleRepository{db: db}

		rows := sqlmock.
			NewRows([]string{
				"id", "uuid", "registration_number", "make", "model", "fuel_type",
				"id", "category", "price_per_day",
			}).
			AddRow(
				mv.ID, mv.UUID, mv.RegistrationNumber, mv.Make, mv.Model, mv.FuelType,
				mv.VehicleCategory.ID, mv.VehicleCategory.VehicleType, mv.VehicleCategory.PricePerDay,
			)

		mock.ExpectQuery(getVehicleByUUID).WithArgs(mv.UUID.String()).WillReturnRows(rows)
		mock.ExpectQuery(selectBookingsByVehicleID).WithArgs(mv.ID).WillReturnRows(sqlmock.NewRows([]string{}))

		v, err := vrp.FindByUUID(mv.UUID)

		assert.NoError(t, err)
		assert.NotNil(t, v)

		err = mock.ExpectationsWereMet()
		if err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

}
