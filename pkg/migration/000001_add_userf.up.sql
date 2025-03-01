CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "telegram_id" BIGINT NOT NULL,
  "username" varchar NOT NULL,
  "score" INTEGER NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);