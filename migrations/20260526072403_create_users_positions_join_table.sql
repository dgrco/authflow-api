-- +goose Up
CREATE TABLE employee_positions (
  user_id UUID NOT NULL REFERENCES users(id),
  position_id UUID NOT NULL REFERENCES positions(id),
  PRIMARY KEY (user_id, position_id)
);

-- +goose Down
DROP TABLE employee_positions;
