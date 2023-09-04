-- +goose Up
-- +goose StatementBegin
CREATE TABLE "users" (
    "id" serial NOT NULL,
    "email" varchar(255) NOT NULL,
    "password" bytea NOT NULL,
    "role" int DEFAULT 0,
    PRIMARY KEY ("id")
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "users";
-- +goose StatementEnd
