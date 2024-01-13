-- +goose Up
-- +goose StatementBegin
CREATE TABLE "users" (
     "id" UUID PRIMARY KEY,
     "password" VARCHAR(256) NOT NULL,
     "first_name" VARCHAR(256),
     "second_name" VARCHAR(256),
     "birthdate" DATE,
     "biography" TEXT,
     "city" VARCHAR(512)
);

CREATE TABLE "tokens" (
    "token" UUID PRIMARY KEY,
    "user_id" UUID REFERENCES "users" ("id") ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "tokens";
DROP TABLE "users";
-- +goose StatementEnd
