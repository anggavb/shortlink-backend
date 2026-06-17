CREATE TABLE "link_clicks" (
  "id" bigserial PRIMARY KEY,
  "link_id" bigint NOT NULL,
  "ip_address" varchar,
  "user_agent" varchar,
  "clicked_at" timestamp DEFAULT (now())
);

COMMENT ON COLUMN "link_clicks"."user_agent" IS 'Mozilla/5.0, PostmanRuntime, curl';

ALTER TABLE "link_clicks" ADD FOREIGN KEY ("link_id") REFERENCES "links" ("id") DEFERRABLE INITIALLY IMMEDIATE;
