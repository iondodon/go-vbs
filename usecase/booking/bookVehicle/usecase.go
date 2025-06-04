package bookVehicle

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/iondodon/go-vbs/domain"
	"github.com/iondodon/go-vbs/dto"
	"github.com/iondodon/go-vbs/usecase"
	"github.com/iondodon/go-vbs/usecase/bookingdate/getBookingDate"
	"github.com/iondodon/go-vbs/usecase/vehicle/isVehicleAvailable"

	"github.com/google/uuid"
)

const alreadyHired = "vehicle with UUID %s is already taken for at leas one day of this period"

type Interface interface {
	ForPeriod(ctx context.Context, tx *sql.Tx, customerUID, vehicleUUID uuid.UUID, period dto.DatePeriodDTO) error
}

type UseCase struct {
	infoLog            *log.Logger
	errorLog           *log.Logger
	vehRepo            usecase.VehicleRepositoryInterface
	custRepo           usecase.CustomerRepositoryInterface
	bookingRepo        usecase.BookingRepositoryInterface
	isAvailableForHire *isVehicleAvailable.UseCase
	getBookingDates    *getBookingDate.UseCase
}

func New(
	infoLog *log.Logger,
	errorLog *log.Logger,
	vehRepo usecase.VehicleRepositoryInterface,
	custRepo usecase.CustomerRepositoryInterface,
	bookingRepo usecase.BookingRepositoryInterface,
	isAvailableForHire *isVehicleAvailable.UseCase,
	getBookingDates *getBookingDate.UseCase,
) Interface {
	return &UseCase{
		infoLog:            infoLog,
		errorLog:           errorLog,
		vehRepo:            vehRepo,
		custRepo:           custRepo,
		bookingRepo:        bookingRepo,
		isAvailableForHire: isAvailableForHire,
		getBookingDates:    getBookingDates,
	}
}

func (uc *UseCase) ForPeriod(ctx context.Context, tx *sql.Tx, customerUID, vehicleUUID uuid.UUID, period dto.DatePeriodDTO) error {
	isAvailable, err := uc.isAvailableForHire.CheckForPeriod(ctx, vehicleUUID, period)
	if err != nil {
		return err
	}
	if !isAvailable {
		return fmt.Errorf(alreadyHired, vehicleUUID)
	}

	uc.infoLog.Printf("Booking vehicle with UUID %s starting from %s to %s \n", vehicleUUID, period.FromDate, period.ToDate)

	bDates, err := uc.getBookingDates.ForPeriod(ctx, tx, customerUID, vehicleUUID, period)
	if err != nil {
		return err
	}

	veh, err := uc.vehRepo.FindByUUID(ctx, vehicleUUID)
	if err != nil {
		return err
	}

	cust, err := uc.custRepo.FindByUUID(ctx, customerUID)
	if err != nil {
		return err
	}

	booking := domain.Booking{
		UUID:         uuid.New(),
		BookingDates: bDates,
		Vehicle:      veh,
		Customer:     cust,
	}

	err = uc.bookingRepo.Save(ctx, tx, &booking)
	if err != nil {
		return err
	}

	return nil
}
