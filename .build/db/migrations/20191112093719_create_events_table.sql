-- migrate:up
CREATE SCHEMA cla;
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

create table cla.events (
  id uuid NOT NULL DEFAULT gen_random_uuid(),
  event_type varchar(255) NOT NULL,
  user_id varchar(255) NOT NULL,
  project_id varchar(255) DEFAULT NULL,
  company_id varchar(255) DEFAULT NULL,
  event_time BIGINT NOT NULL DEFAULT extract(epoch from now()),
  event_data jsonb,
  PRIMARY KEY(id)
);

-- migrate:down
drop table cla.events;
drop schema cla;
