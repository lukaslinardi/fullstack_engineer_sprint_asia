CREATE TABLE "task"
(
  "id" SERIAL PRIMARY KEY,
  "task_name" VARCHAR(255) NOT NULL,
  "task_status" SMALLINT NOT NULL,
  "created_at" timestamptz DEFAULT current_timestamp,
  "updated_at" timestamptz DEFAULT current_timestamp,
);
