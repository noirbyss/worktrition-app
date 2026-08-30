CREATE TABLE user_weight_measurements (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    measured_on DATE NOT NULL,
    weight_kg NUMERIC(5, 2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (user_id, measured_on),
    CONSTRAINT user_weight_measurements_weight_valid CHECK (weight_kg BETWEEN 25 AND 400)
);

CREATE INDEX user_weight_measurements_measured_on_idx
ON user_weight_measurements (user_id, measured_on DESC);

CREATE TRIGGER user_weight_measurements_set_updated_at
BEFORE UPDATE ON user_weight_measurements
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();
