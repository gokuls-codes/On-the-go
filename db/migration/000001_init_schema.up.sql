CREATE TABLE "projects" (
  "id" integer PRIMARY KEY AUTOINCREMENT,
  "name" varchar NOT NULL,
  "description" varchar,
  "githubURL" url NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "deployments" (
  "id" integer PRIMARY KEY,
  "project_id" integer NOT NULL,
  "deployed_at" timestamp,
  "status" varchar NOT NULL,
  "started_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,

  FOREIGN KEY ("project_id") REFERENCES "projects" ("id")
);

CREATE TABLE "env_vars" (
  "id" integer PRIMARY KEY,
  "project_id" integer NOT NULL,
  "key" varchar NOT NULL,
  "value" varchar NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,

  FOREIGN KEY ("project_id") REFERENCES "projects" ("id")
);