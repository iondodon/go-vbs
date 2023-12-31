package vehicle

import (
	"testing"

	"github.com/iondodon/go-vbs/domain"
	vehRepo "github.com/iondodon/go-vbs/repository/vehicle"

	uuidLib "github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_getVehicle_ByUUID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVehRepo := vehRepo.NewMockVehicleRepository(ctrl)
	var gvuc GetVehicleUseCase = &getVehicleUseCase{vehicleRepository: mockVehRepo}

	uuid, err := uuidLib.Parse("c2df2b03-92e8-41ad-9a74-0b7b040a4cf5")
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}

	expectedVeh := &domain.Vehicle{}

	t.Run("the repository is called with the correct UUID", func(t *testing.T) {
		mockVehRepo.EXPECT().FindByUUID(uuid).Return(expectedVeh, nil).Times(1)

		veh, err := gvuc.ByUUID(uuid)

		assert.Nil(t, err)
		assert.Equal(t, expectedVeh, veh)
	})

}
