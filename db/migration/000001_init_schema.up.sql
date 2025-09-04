CREATE TABLE "projects" (
  "id" integer PRIMARY KEY AUTOINCREMENT,
  "name" varchar NOT NULL,
  "description" varchar,
  "github_url" url NOT NULL,
  "repo_name" varchar NOT NULL,
  "container_port" integer,
  "host_port" integer,
  "image_id" varchar,
  "container_id" varchar,
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