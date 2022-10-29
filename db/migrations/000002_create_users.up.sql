CREATE TABLE "users" (
    "id" bigserial PRIMARY KEY,
    "fullname" varchar NOT NULL,
    "username" varchar NOT NULL,
    "email" varchar NOT NULL,
    "password" varchar NOT NULL,
    created_at timestamptz NOT NULL DEFAULT (now())
);

CREATE UNIQUE INDEX ON "users" ("username");
CREATE UNIQUE INDEX ON "users" ("email");