CREATE TABLE cars (
    car_id SERIAL PRIMARY KEY NOT NULL,
    car_name CHAR(50) NOT NULL,
    day_rate decimal NOT NULL,
    month_rate decimal NOT NULL,
    image CHAR(256) NOT NULL
);

CREATE TABLE orders (
    order_id SERIAL PRIMARY KEY NOT NULL,
    car_id int NOT NULL,
    order_date TIMESTAMPTZ DEFAULT NOW(),
    pickup_date TIMESTAMPTZ NOT NULL,
    dropoff_date TIMESTAMPTZ NOT NULL,
    pickup_location CHAR(50) NOT NULL,
    dropoff_location CHAR(50) NOT NULL
);