package bookVehicle

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/iondodon/go-vbs/business"
	"github.com/iondodon/go-vbs/business/bookingdate/getBookingDate"
	"github.com/iondodon/go-vbs/business/vehicle/isVehicleAvailable"
	"github.com/iondodon/go-vbs/domain"
	"github.com/iondodon/go-vbs/dto"

	"github.com/google/uuid"
)

const alreadyHired = "vehicle with UUID %s is already taken for at leas one day of this period"

type Service struct {
	infoLog            *log.Logger
	errorLog           *log.Logger
	vehRepo            business.VehicleRepository
	custRepo           business.CustomerRepository
	bookingRepo        business.BookingRepository
	isAvailableForHire *isVehicleAvailable.Service
	getBookingDates    *getBookingDate.Service
}

func New(
	infoLog *log.Logger,
	errorLog *log.Logger,
	vehRepo business.VehicleRepository,
	custRepo business.CustomerRepository,
	bookingRepo business.BookingRepository,
	isAvailableForHire *isVehicleAvailable.Service,
	getBookingDates *getBookingDate.Service,
) *Service {
	return &Service{
		infoLog:            infoLog,
		errorLog:           errorLog,
		vehRepo:            vehRepo,
		custRepo:           custRepo,
		bookingRepo:        bookingRepo,
		isAvailableForHire: isAvailableForHire,
		getBookingDates:    getBookingDates,
	}
}

func (uc *Service) ForPeriod(ctx context.Context, tx *sql.Tx, customerUID, vehicleUUID uuid.UUID, period dto.DatePeriodDTO) error {
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
