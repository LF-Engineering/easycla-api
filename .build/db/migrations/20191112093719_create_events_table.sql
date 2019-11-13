-- migrate:up

create table events (
  id uuid PRIMARY KEY,
  event_type varchar(255) NOT NULL,
  user_id varchar(255) NOT NULL,
  project_id varchar(255) DEFAULT NULL,
  company_id varchar(255) DEFAULT NULL,
  event_time BIGINT NOT NULL DEFAULT extract(epoch from now()),
  event_data jsonb
);

-- migrate:down
drop table events;
