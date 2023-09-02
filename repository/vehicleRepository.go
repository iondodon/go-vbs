package repository

import (
	"database/sql"
	"go-vbs/domain"
	"log"

	uuidLib "github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

const insertMockData = `
	INSERT INTO vehicle_category(id, category, price_per_day)
	VALUES (1, 'Van', 12.321);

	INSERT INTO vehicle (id, uuid, registration_number, make, model, fuel_type, category_id)
	VALUES (1, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'ABC-123', 'Make', 'Y', "diesel", 1);

	INSERT INTO customer (id, uuid, username)
	VALUES (1, 'eba846c2-1d57-4f5d-b17e-fa9f922ac093', 'username123');

	INSERT INTO booking(id, uuid, vehicle_id, customer_id)
	VALUES (1, 'de399bc0-a622-4449-b264-5783562c38fa', 1, 1);

	INSERT INTO booking_date (id, time, booking_id)
	VALUES (1, current_date, 1);
`

const ddl = `
	CREATE TABLE vehicle_category (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		category VARCHAR(10) NOT NULL,
		price_per_day REAL NOT NULL
	);

	CREATE TABLE booking_date (
	    id BIGINT AUTO_INCREMENT PRIMARY KEY,
	    time TIMESTAMP NOT NULL,
	    booking_id BIGINT NOT NULL,
	    FOREIGN KEY (booking_id) REFERENCES booking(id)
	);

	CREATE TABLE customer (
	    id BIGINT AUTO_INCREMENT PRIMARY KEY,
	    uuid CHAR(36) NOT NULL,
	    username VARCHAR(20) NOT NULL
	);

	CREATE TABLE booking (
	    id BIGINT AUTO_INCREMENT PRIMARY KEY,
	    uuid CHAR(36) UNIQUE NOT NULL,
	    vehicle_id BIGINT NOT NULL,
	    customer_id BIGINT NOT NULL,
	    FOREIGN KEY (customer_id) REFERENCES customer(id),
	    FOREIGN KEY (vehicle_id) REFERENCES vehicle(id)
	);

	CREATE TABLE vehicle (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		uuid CHAR(36) UNIQUE NOT NULL,
		registration_number VARCHAR(10) UNIQUE NOT NULL,
		make VARCHAR(20) NOT NULL,
		model VARCHAR(20) NOT NULL,
		fuel_type VARCHAR(10) NOT NULL,
		category_id BIGINT NOT NULL,
		FOREIGN KEY (category_id) REFERENCES vehicle_category(id)
	);
`

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

const selectBookingsByVehicleId = `
	SELECT id, uuid
	FROM booking
	WHERE vehicle_id = ?
`

type VehicleRepository interface {
	FindByUUID(vid uuidLib.UUID) (*domain.Vehicle, error)
}

type vehicleRepository struct {
	db *sql.DB
}

func NewVehicleRepository() VehicleRepository {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	//defer func(db *sql.DB) {
	//	err := db.Close()
	//	if err != nil {
	//
	//	}
	//}(db)

	// DDL
	_, err = db.Exec(ddl)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	// Insert mock data
	_, err = db.Exec(insertMockData)
	if err != nil {
		log.Fatalf("Failed to insert row: %v", err)
	}

	return &vehicleRepository{
		db: db,
	}
}

func (vrp *vehicleRepository) FindByUUID(vUUID uuidLib.UUID) (*domain.Vehicle, error) {
	var vehicle domain.Vehicle
	var vehCat domain.VehicleCategory
	var vehUUID string
	err := vrp.db.QueryRow(getVehicleByUUID, vUUID.String()).Scan(
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

	rows, err := vrp.db.Query(selectBookingsByVehicleId, vehicle.ID)
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
		err := rows.Scan(&booking.ID, &booking.UUID)
		if err != nil {
			return nil, err
		}
		vehicle.Bookings = append(vehicle.Bookings, &booking)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &vehicle, nil
}
