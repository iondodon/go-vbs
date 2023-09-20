package vehicle

import (
	"database/sql"

	"github.com/iondodon/go-vbs/domain"
	"github.com/iondodon/go-vbs/dto"
	"github.com/iondodon/go-vbs/integration"

	uuidLib "github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

const getVehicleByUUID = `
	SELECT 
	    v.id, 
	    v.uuid, 
	    v.registration_number, 
	    v.make, 
	    v.model, 
	    v.fuel_type, 
	    vc.id, 
	    vc.category, 
	    vc.price_per_day 
	FROM vehicle v
		JOIN vehicle_category vc on v.category_id = vc.id
	WHERE v.uuid = ?
`

const selectBookingsByVehicleID = `
	SELECT b.id, b.uuid, c.uuid, c.username
	FROM booking b
		JOIN customer c on b.customer_id = c.id
	WHERE b.vehicle_id = ?
`

const selectBookingDatesByBookingID = `
	SELECT bd.id, bd.time
	FROM booking_date bd
	WHERE bd.id = ?
`

const vehicleHasBookedDatesOnPeriod = `
	SELECT EXISTS(
		SELECT 1 
		FROM booking b
			JOIN vehicle v on b.vehicle_id = v.id
			JOIN bookings_bookingdates bb on bb.booking_id = b.id
			JOIN booking_date bd on bb.bookingdate_id = bd.id
		WHERE v.uuid = ? AND bd.time >= ? and bd.time <= ?
	)
`

type VehicleRepository interface {
	FindByUUID(vUUID uuidLib.UUID) (*domain.Vehicle, error)
	VehicleHasBookedDatesOnPeriod(vUUID uuidLib.UUID, period dto.DatePeriodDTO) (bool, error)
}

type vehicleRepository struct {
	db integration.DB
}

func NewVehicleRepository(db integration.DB) VehicleRepository {
	return &vehicleRepository{db: db}
}

func (repo *vehicleRepository) FindByUUID(vUUID uuidLib.UUID) (*domain.Vehicle, error) {
	var vehicle domain.Vehicle
	var vehCat domain.VehicleCategory
	var vehUUID string
	err := repo.db.QueryRow(getVehicleByUUID, vUUID.String()).Scan(
		&vehicle.ID,
		&vehUUID,
		&vehicle.RegistrationNumber,
		&vehicle.Make,
		&vehicle.Model,
		&vehicle.FuelType,
		&vehCat.ID,
		&vehCat.VehicleType,
		&vehCat.PricePerDay,
	)
	if err != nil {
		return nil, err
	}
	vehicle.VehicleCategory = &vehCat

	vehicle.UUID, err = uuidLib.Parse(vehUUID)
	if err != nil {
		return nil, err
	}

	rows, err := repo.db.Query(selectBookingsByVehicleID, vehicle.ID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	for rows.Next() {
		var booking domain.Booking
		var customer domain.Customer
		err := rows.Scan(&booking.ID, &booking.UUID, &customer.UUID, &customer.Username)
		if err != nil {
			return nil, err
		}
		booking.Vehicle = &vehicle
		booking.Customer = &customer
		vehicle.Bookings = append(vehicle.Bookings, &booking)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &vehicle, nil
}

func (repo *vehicleRepository) VehicleHasBookedDatesOnPeriod(vUUID uuidLib.UUID, period dto.DatePeriodDTO) (bool, error) {
	var exists bool

	err := repo.db.QueryRow(vehicleHasBookedDatesOnPeriod, vUUID, period.FromDate, period.ToDate).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
