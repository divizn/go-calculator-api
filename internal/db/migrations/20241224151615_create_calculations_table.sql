-- +goose Up
-- +goose StatementBegin
CREATE TABLE calculations (
    id SERIAL PRIMARY KEY,
    num1 FLOAT NOT NULL,
    num2 FLOAT NOT NULL,
    operator VARCHAR(2) NOT NULL,
    result FLOAT DEFAULT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS calculations;
-- +goose StatementEnd
