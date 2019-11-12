-- migrate:up

create table audit_log (
  id uuid PRIMARY KEY,
  event varchar(255) NOT NULL,
  username varchar(255) NOT NULL,
  event_time timestamptz NOT NULL DEFAULT now(),
  event_data json
);

-- migrate:down
drop table audit_log;
