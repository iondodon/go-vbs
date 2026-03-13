-- +goose Up
-- +goose StatementBegin
INSERT INTO vehicle_category(id, category, price_per_day)
VALUES (1, 'Van', 12.321);

INSERT INTO vehicle (id, uuid, registration_number, make, model, fuel_type, category_id)
VALUES (1, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'ABC-123', 'Make', 'Y', "DIESEL", 1);

INSERT INTO customer (id, uuid, username)
VALUES (1, 'eba846c2-1d57-4f5d-b17e-fa9f922ac093', 'username123');

INSERT INTO booking(id, uuid, vehicle_id, customer_id)
VALUES (1, 'de399bc0-a622-4449-b264-5783562c38fa', 1, 1);

INSERT INTO booking_date (id, time)
VALUES (1, current_date);

INSERT INTO bookings_bookingdates (booking_id, bookingdate_id)
VALUES (1, 1);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- +goose StatementEnd
