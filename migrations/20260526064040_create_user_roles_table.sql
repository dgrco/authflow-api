-- +goose Up
CREATE TYPE user_role AS ENUM ('admin', 'manager', 'employee');

CREATE TABLE user_roles (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id),
  business_id UUID NOT NULL REFERENCES businesses(id),
  location_id UUID REFERENCES locations(id),
  role user_role NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT unique_user_business UNIQUE (user_id, business_id)
);

-- +goose Down
DROP TABLE user_roles;
DROP TYPE user_role;
