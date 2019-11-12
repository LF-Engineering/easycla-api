-- migrate:up

create table audit_events (
  id uuid PRIMARY KEY,
  event_type varchar(255) NOT NULL,
  username varchar(255) NOT NULL,
  event_time BIGINT NOT NULL DEFAULT extract(epoch from now()),
  event_data json
);

-- migrate:down
drop table events;
