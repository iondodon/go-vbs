package vehicle

// import (
// 	"context"
// 	"testing"

// 	"github.com/iondodon/go-vbs/domain"
// 	vehRepo "github.com/iondodon/go-vbs/repository/vehicle"

// 	"github.com/google/uuid"
// 	"github.com/stretchr/testify/assert"
// )

// func Test_getVehicle_ByUUID(t *testing.T) {
// 	mockVehRepo := vehRepo.NewMockVehicleRepository(t)
// 	var gvuc GetVehicle = &getVehicle{vehicleRepository: mockVehRepo}

// 	uuid, err := uuid.Parse("c2df2b03-92e8-41ad-9a74-0b7b040a4cf5")
// 	if err != nil {
// 		t.Fatalf("unexpected error: %s", err.Error())
// 	}

// 	expectedVeh := &domain.Vehicle{}

// 	t.Run("the repository is called with the correct UUID", func(t *testing.T) {
// 		mockVehRepo.EXPECT().FindByUUID(context.Background(), uuid).Return(expectedVeh, nil).Times(1)

// 		veh, err := gvuc.ByUUID(context.Background(), uuid)

// 		assert.Nil(t, err)
// 		assert.Equal(t, expectedVeh, veh)
// 	})

// }
