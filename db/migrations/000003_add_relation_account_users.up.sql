ALTER TABLE "accounts" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");
ALTER TABLE "accounts" ADD CONSTRAINT "owner_currency_key" UNIQUE ("owner", "currency")