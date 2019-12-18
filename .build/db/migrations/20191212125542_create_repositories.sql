-- migrate:up
create table cla.repositories (
  id uuid NOT NULL DEFAULT gen_random_uuid(),
  repository_type varchar(255) NOT NULL,
  name varchar(255) NOT NULL,
  organization_name varchar(255) NOT NULL,
  external_id varchar(4096) NOT NULL,
  url varchar(4096) NOT NULL,
  enabled boolean NOT NULL default TRUE,
  project_id varchar(255) NOT NULL,
  cla_group_id uuid NOT NULL REFERENCES cla.cla_groups(id),
  created_at BIGINT NOT NULL DEFAULT extract(epoch from now()),
  PRIMARY KEY(id)
);

-- migrate:down
drop table cla.repositories;
