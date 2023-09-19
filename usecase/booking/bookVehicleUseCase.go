package booking

import (
	"fmt"
	"log"

	"github.com/iondodon/go-vbs/dto"
	custRepo "github.com/iondodon/go-vbs/repository/customer"
	vehRepo "github.com/iondodon/go-vbs/repository/vehicle"
	bdUCs "github.com/iondodon/go-vbs/usecase/bookingdate"
	vehUCs "github.com/iondodon/go-vbs/usecase/vehicle"

	"github.com/google/uuid"
)

const alreadyHired = "vehicle with UUID %s is already taken for at leas one day of this period"

type BookVehicleUseCase interface {
	ForPeriod(customerUID, vehicleUUID uuid.UUID, period dto.DatePeriodDTO) error
}

type bookVehicleUseCase struct {
	infoLog, errorLog      *log.Logger
	vehRepo                vehRepo.VehicleRepository
	custRepo               custRepo.CustomerRepository
	isAvailableForHireUS   vehUCs.IsAvailableForHireUseCase
	getBookingDatesUseCase bdUCs.GetBookingDatesUseCase
}

func NewBookVehicleUseCase(
	infoLog, errorLog *log.Logger,
	vrp vehRepo.VehicleRepository,
	crp custRepo.CustomerRepository,
	isAvailableForHireUS vehUCs.IsAvailableForHireUseCase,
	getBookingDatesUseCase bdUCs.GetBookingDatesUseCase,
) BookVehicleUseCase {
	return &bookVehicleUseCase{
		infoLog:                infoLog,
		errorLog:               errorLog,
		vehRepo:                vrp,
		custRepo:               crp,
		isAvailableForHireUS:   isAvailableForHireUS,
		getBookingDatesUseCase: getBookingDatesUseCase,
	}
}

func (uc *bookVehicleUseCase) ForPeriod(customerUID, vehicleUUID uuid.UUID, period dto.DatePeriodDTO) error {
	uc.infoLog.Printf("Booking vehicle with UUID %s starting from %s to %s \n", vehicleUUID, period.FromDate, period.ToDate)

	isAvailable, err := uc.isAvailableForHireUS.CheckForPeriod(vehicleUUID, period)
	if err != nil {
		return err
	}
	if !isAvailable {
		return fmt.Errorf(alreadyHired, vehicleUUID)
	}

	_, err = uc.vehRepo.FindByUUID(vehicleUUID)
	if err != nil {
		return err
	}

	_, err = uc.custRepo.FindByUUID(customerUID)
	if err != nil {
		return err
	}

	uc.infoLog.Printf("Booking vehicle with UUID %s \n", vehicleUUID)

	_, err = uc.getBookingDatesUseCase.ForPeriod(customerUID, vehicleUUID, period)
	if err != nil {
		return err
	}

	// create new Booking and save it

	return nil
}
