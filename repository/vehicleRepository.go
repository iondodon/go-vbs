package repository

import (
	"database/sql"
	"errors"
	"fmt"
	uuidLib "github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"go-vbs/domain"
	"log"
)

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
	_, err = db.Exec(`
		CREATE TABLE vehicle (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			uuid CHAR(36) UNIQUE NOT NULL,
			registrationNumber VARCHAR(10) UNIQUE NOT NULL,
			make VARCHAR(20) NOT NULL,
			model VARCHAR(20) NOT NULL,
			fuelType VARCHAR(10) NOT NULL,
			vehCatID BIGINT NOT NULL,
			vehCatType VARCHAR(10) NOT NULL
		)
	`)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	// Insert mock data
	_, err = db.Exec(`
		INSERT INTO vehicle (id, uuid, registrationNumber, make, model, fuelType, vehCatID, vehCatType)
		VALUES (1, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'ABC-123', 'Make', 'Y', "diesel", 123, "Van");
	`)
	if err != nil {
		log.Fatalf("Failed to insert row: %v", err)
	}

	return &vehicleRepository{
		db: db,
	}
}

func (vrp *vehicleRepository) FindByUUID(vUUID uuidLib.UUID) (*domain.Vehicle, error) {
	var veh domain.Vehicle
	var vehUUID string
	err := vrp.db.QueryRow("SELECT * FROM vehicle WHERE uuid = ?", vUUID.String()).Scan(
		&veh.ID,
		&vehUUID,
		&veh.RegistrationNumber,
		&veh.Make,
		&veh.Model,
		&veh.FuelType,
		&veh.VehicleCategory.ID,
		&veh.VehicleCategory.VehicleType,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("No rows returned")
		} else {
			log.Fatalf("Failed to query table: %v", err)
		}
		return nil, fmt.Errorf("error")
	}

	veh.UUID, err = uuidLib.Parse(vehUUID)
	if err != nil {
		return nil, err
	}

	return &veh, nil
}
