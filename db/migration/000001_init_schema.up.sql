CREATE TABLE "projects" (
  "id" integer PRIMARY KEY,
  "name" varchar,
  "description" varchar,
  "githubURL" url,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "deployments" (
  "id" integer PRIMARY KEY,
  "project_id" integer NOT NULL,
  "deployed_at" timestamp,

  FOREIGN KEY ("project_id") REFERENCES "projects" ("id")
);