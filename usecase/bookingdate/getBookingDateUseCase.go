package bookingdate

import (
	"github.com/google/uuid"
	"github.com/iondodon/go-vbs/dto"
	bdRepoPkg "github.com/iondodon/go-vbs/repository/bookingdate"
)

type GetBookingDatesUseCase interface {
	ForPeriod(customerUUID, vehicleUUID uuid.UUID, period dto.DatePeriodDTO) error
}

type getBookingDatesUseCase struct {
	bdRepo bdRepoPkg.BookingDateRepository
}

func NewGetBookingDatesUseCase(bdRepo bdRepoPkg.BookingDateRepository) GetBookingDatesUseCase {
	return &getBookingDatesUseCase{bdRepo: bdRepo}
}

func (uc *getBookingDatesUseCase) ForPeriod(customerUUID, vehicleUUID uuid.UUID, period dto.DatePeriodDTO) error {
	// TODO
	return nil
}
