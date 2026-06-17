CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "email" varchar UNIQUE NOT NULL,
  "password" varchar NOT NULL,
  "first_name" varchar,
  "last_name" varchar,
  "role" varchar,
  "workplace" varchar,
  "is_member" bool DEFAULT false,
  "is_receive_email" bool DEFAULT true,
  "photo" varchar,
  "verified_at" timestamp,
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp
);