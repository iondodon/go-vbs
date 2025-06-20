package bookVehicleService

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/iondodon/go-vbs/booking/bookingBusiness"
	"github.com/iondodon/go-vbs/booking/bookingController"
	"github.com/iondodon/go-vbs/booking/bookingDomain"
)

const alreadyHired = "vehicle with UUID %s is already taken for at leas one day of this period"

// Service handles vehicle booking
type Service struct {
	vehicleRepo         bookingBusiness.VehicleRepository
	customerRepo        bookingBusiness.CustomerRepository
	bookingRepo         bookingBusiness.BookingRepository
	bookingDateRepo     bookingBusiness.BookingDateRepository
	availabilityService bookingBusiness.VehicleAvailabilityService
}

func New(
	vehicleRepo bookingBusiness.VehicleRepository,
	customerRepo bookingBusiness.CustomerRepository,
	bookingRepo bookingBusiness.BookingRepository,
	bookingDateRepo bookingBusiness.BookingDateRepository,
	availabilityService bookingBusiness.VehicleAvailabilityService,
) *Service {
	return &Service{
		vehicleRepo:         vehicleRepo,
		customerRepo:        customerRepo,
		bookingRepo:         bookingRepo,
		bookingDateRepo:     bookingDateRepo,
		availabilityService: availabilityService,
	}
}

func (s *Service) ForPeriod(ctx context.Context, tx *sql.Tx, customerUID, vehicleUUID uuid.UUID, period bookingController.DatePeriodDTO) error {
	isAvailable, err := s.availabilityService.CheckForPeriod(ctx, vehicleUUID, period)
	if err != nil {
		return err
	}
	if !isAvailable {
		return fmt.Errorf(alreadyHired, vehicleUUID)
	}

	slog.Info("Booking vehicle", "vehicleUUID", vehicleUUID, "fromDate", period.FromDate, "toDate", period.ToDate)

	bDates, err := s.getBookingDatesForPeriod(ctx, tx, customerUID, vehicleUUID, period)
	if err != nil {
		return err
	}

	veh, err := s.vehicleRepo.FindByUUID(ctx, vehicleUUID)
	if err != nil {
		return err
	}

	cust, err := s.customerRepo.FindByUUID(ctx, customerUID)
	if err != nil {
		return err
	}

	booking := bookingDomain.Booking{
		UUID:         uuid.New(),
		BookingDates: bDates,
		Vehicle:      veh,
		Customer:     cust,
	}

	err = s.bookingRepo.Save(ctx, tx, &booking)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) getBookingDatesForPeriod(ctx context.Context, tx *sql.Tx, customerUUID, vehicleUUID uuid.UUID, period bookingController.DatePeriodDTO) ([]*bookingDomain.BookingDate, error) {
	persistedBookingDates, err := s.bookingDateRepo.FindAllInPeriodInclusive(ctx, period.FromDate, period.ToDate)
	if err != nil {
		return nil, err
	}

	var dates []time.Time
	for _, bd := range persistedBookingDates {
		bd.Time = removeTimePart(bd.Time)
		dates = append(dates, bd.Time)
	}

	fromDate := removeTimePart(period.FromDate)
	toDate := removeTimePart(period.ToDate)
	for date := fromDate; date.Before(toDate) || date.Equal(toDate); date = date.Add(time.Hour + 24) {
		if !contains(dates, date) {
			bd := bookingDomain.BookingDate{Time: date}
			s.bookingDateRepo.Save(ctx, tx, &bd)
			persistedBookingDates = append(persistedBookingDates, &bd)
		}
	}

	return persistedBookingDates, nil
}

// Helper functions
func removeTimePart(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

func contains(dates []time.Time, targetDate time.Time) bool {
	for _, d := range dates {
		if d.Equal(targetDate) {
			return true
		}
	}
	return false
}
