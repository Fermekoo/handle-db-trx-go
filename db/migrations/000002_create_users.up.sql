CREATE TABLE "users" (
    "id" bigserial PRIMARY KEY,
    "fullname" varchar NOT NULL,
    "username" varchar UNIQUE NOT NULL,
    "email" varchar UNIQUE NOT NULL,
    "password" varchar NOT NULL,
    "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
    created_at timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "accounts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "accounts" ADD CONSTRAINT "user_id_currency_key" UNIQUE ("user_id", "currency");