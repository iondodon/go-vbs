package booking

import (
	"context"
	"database/sql"
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

//gobok:constructor
//ctxboot:component
type BookVehicle struct {
	infoLog            *log.Logger                   `ctxboot:"inject"`
	errorLog           *log.Logger                   `ctxboot:"inject"`
	vehRepo            vehRepo.VehicleRepository     `ctxboot:"inject"`
	custRepo           custRepo.CustomerRepository   `ctxboot:"inject"`
	bookingRepo        bookingRepo.BookingRepository `ctxboot:"inject"`
	isAvailableForHire vehUCs.IsAvailableForHire     `ctxboot:"inject"`
	getBookingDates    bdUCs.GetBookingDates         `ctxboot:"inject"`
}

func (uc *BookVehicle) ForPeriod(ctx context.Context, tx *sql.Tx, customerUID, vehicleUUID uuid.UUID, period dto.DatePeriodDTO) error {
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
