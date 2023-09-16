package repository

import (
	"database/sql"
	uuidLib "github.com/google/uuid"
	"go-vbs/integration/mocks"
	"go.uber.org/mock/gomock"
	"testing"
)

func Test_vehicleRepository_FindByUUID(t *testing.T) {
	t.Run("name", func(t *testing.T) {
		mockDB := mocks.NewMockDB(gomock.NewController(t))

		var vrp VehicleRepository = &vehicleRepository{
			db: mockDB,
		}

		uuid, _ := uuidLib.Parse("d06a5744-ce7d-4aa7-ba47-076cae095bb1")

		mockDB.EXPECT().
			QueryRow(getVehicleByUUID, uuid.String()).
			Return(&sql.Row{}).
			Times(1)

		_, _ = vrp.FindByUUID(uuid)
	})

}
