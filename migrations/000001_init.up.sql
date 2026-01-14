CREATE TABLE IF NOT EXISTS cargo_type (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    process_cost DECIMAL(10, 2) NOT NULL
);

CREATE TABLE IF NOT EXISTS vessel (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    vessel_type VARCHAR(200) NOT NULL,
    max_load DECIMAL(10, 2) NOT NULL CHECK (max_load > 0)
);


CREATE TABLE IF NOT EXISTS operation (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS cargo (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    type_id INTEGER NOT NULL REFERENCES cargo_type(id),
    weight DECIMAL(10, 2) NOT NULL CHECK (weight > 0),
    volume DECIMAL(10, 2) NOT NULL CHECK (volume > 0),
    vessel_id INTEGER NOT NULL REFERENCES vessel(id)
);

CREATE TABLE IF NOT EXISTS storage_loc (
    id SERIAL PRIMARY KEY,
    cargo_type_id INTEGER NOT NULL REFERENCES cargo_type(id), 
    max_weight DECIMAL(10, 2) NOT NULL,
    max_volume DECIMAL(10, 2) NOT NULL,
    cargo_id INTEGER REFERENCES cargo(id), 
    date_of_placement TIMESTAMP
);

CREATE TABLE IF NOT EXISTS operation_cargo (
    operation_id INTEGER REFERENCES operation(id) NOT NULL,
    cargo_id INTEGER REFERENCES cargo(id) NOT NULL,
    PRIMARY KEY (operation_id, cargo_id)
);