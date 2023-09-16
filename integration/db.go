package integration

import (
	"database/sql"
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

type DB interface {
	Close() error
	QueryRow(query string, args ...any) *sql.Row
	Query(query string, args ...any) (*sql.Rows, error)
}

func NewInMemDBConn() (DB, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// DDL
	if _, err = db.Exec(ddl); err != nil {
		return nil, err
	}

	// Insert mock data
	if _, err = db.Exec(insertMockData); err != nil {
		return nil, err
	}

	return db, nil
}
