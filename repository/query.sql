-- name: GetVehicleByUUID :one
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
WHERE v.uuid = ?;

-- name: SelectBookingsByVehicleID :many
SELECT b.id, b.uuid, c.uuid, c.username
FROM booking b
    JOIN customer c on b.customer_id = c.id
WHERE b.vehicle_id = ?;

-- name: SelectBookingDatesByBookingID :one
SELECT bd.id, bd.time
FROM booking_date bd
WHERE bd.id = ?;

-- name: VehicleHasBookedDatesOnPeriod :one
SELECT EXISTS(
    SELECT 1 
    FROM booking b
        JOIN vehicle v on b.vehicle_id = v.id
        JOIN bookings_bookingdates bb on bb.booking_id = b.id
        JOIN booking_date bd on bb.bookingdate_id = bd.id
    WHERE v.uuid = ? AND bd.time >= ? and bd.time <= ?
);

-- name: InsertNewBooking :exec
INSERT INTO booking(uuid, vehicle_id, customer_id)
VALUES (?, ?, ?);

-- name: SelectAllBookings :many
SELECT b.id, b.uuid, b.vehicle_id, b.customer_id 
FROM booking b;

-- name: FindAllInPeriodInclusive :many
SELECT bd.id, bd.time
FROM booking_date bd
WHERE bd.time >= ? AND bd.time <= ?;

-- name: SaveNewBookingDate :exec
INSERT INTO booking_date(time)
VALUES (?);

-- name: GetCustomerByUUID :one
SELECT c.id, c.uuid, c.username
FROM customer c
WHERE c.uuid = ?;