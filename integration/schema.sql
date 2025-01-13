CREATE TABLE vehicle_category (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    category VARCHAR(10) NOT NULL,
    price_per_day REAL NOT NULL
);

CREATE TABLE booking_date (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    time TIMESTAMP NOT NULL
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

CREATE TABLE bookings_bookingdates (
    booking_id BIGINT NOT NULL,
    bookingdate_id BIGINT NOT NULL,
    FOREIGN KEY (booking_id) REFERENCES booking(id),
    FOREIGN KEY (bookingdate_id) REFERENCES booking_date(id)
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
