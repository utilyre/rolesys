-- +goose Up
-- +goose StatementBegin
CREATE TABLE "sessions" (
    "id" serial NOT NULL,
    "user_id" serial NOT NULL,
    "token" char(36) UNIQUE NOT NULL,
    "expires_at" timestam,
    PRIMARY KEY ("id"),
    CONSTRAINT "fk_user" FOREIGN KEY ("user_id") REFERENCES "users"("id")
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "sessions";
-- +goose StatementEnd
