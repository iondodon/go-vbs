package booking

import (
	"fmt"
	"log"

	"github.com/iondodon/go-vbs/domain"
	"github.com/iondodon/go-vbs/dto"
	bookingRepo "github.com/iondodon/go-vbs/repository/booking"
	custRepo "github.com/iondodon/go-vbs/repository/customer"
	vehRepo "github.com/iondodon/go-vbs/repository/vehicle"
	bdUCs "github.com/iondodon/go-vbs/usecase/bookingdate"
	vehUCs "github.com/iondodon/go-vbs/usecase/vehicle"

	"github.com/google/uuid"
)

const alreadyHired = "vehicle with UUID %s is already taken for at leas one day of this period"

type BookVehicle interface {
	ForPeriod(customerUUID, vehicleUUID uuid.UUID, period dto.DatePeriodDTO) error
}

type bookVehicle struct {
	infoLog, errorLog  *log.Logger
	vehRepo            vehRepo.VehicleRepository
	custRepo           custRepo.CustomerRepository
	bookingRepo        bookingRepo.BookingRepository
	isAvailableForHire vehUCs.IsAvailableForHire
	getBookingDates    bdUCs.GetBookingDates
}

func NewBookVehicle(
	infoLog, errorLog *log.Logger,
	vrp vehRepo.VehicleRepository,
	crp custRepo.CustomerRepository,
	brp bookingRepo.BookingRepository,
	isAvailableForHireUS vehUCs.IsAvailableForHire,
	getBookingDatesUseCase bdUCs.GetBookingDates,
) BookVehicle {
	return &bookVehicle{
		infoLog:            infoLog,
		errorLog:           errorLog,
		vehRepo:            vrp,
		custRepo:           crp,
		bookingRepo:        brp,
		isAvailableForHire: isAvailableForHireUS,
		getBookingDates:    getBookingDatesUseCase,
	}
}

func (uc *bookVehicle) ForPeriod(customerUID, vehicleUUID uuid.UUID, period dto.DatePeriodDTO) error {
	isAvailable, err := uc.isAvailableForHire.CheckForPeriod(vehicleUUID, period)
	if err != nil {
		return err
	}
	if !isAvailable {
		return fmt.Errorf(alreadyHired, vehicleUUID)
	}

	uc.infoLog.Printf("Booking vehicle with UUID %s starting from %s to %s \n", vehicleUUID, period.FromDate, period.ToDate)

	bDates, err := uc.getBookingDates.ForPeriod(customerUID, vehicleUUID, period)
	if err != nil {
		return err
	}

	veh, err := uc.vehRepo.FindByUUID(vehicleUUID)
	if err != nil {
		return err
	}

	cust, err := uc.custRepo.FindByUUID(customerUID)
	if err != nil {
		return err
	}

	booking := domain.Booking{
		UUID:         uuid.New(),
		BookingDates: bDates,
		Vehicle:      veh,
		Customer:     cust,
	}

	err = uc.bookingRepo.Save(&booking)
	if err != nil {
		return err
	}

	return nil
}
