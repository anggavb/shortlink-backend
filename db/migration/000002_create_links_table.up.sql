CREATE TABLE "links" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "original_url" varchar NOT NULL UNIQUE,
  "slug" varchar UNIQUE,
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp,
  "deleted_at" timestamp
);

ALTER TABLE "links" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") DEFERRABLE INITIALLY IMMEDIATE;