-- migrate:up

create table cla.cla_groups (
  id uuid NOT NULL DEFAULT gen_random_uuid(),
  cla_group_name varchar(255) NOT NULL,
  foundation_id varchar(255) NOT NULL,
  created_at BIGINT NOT NULL DEFAULT extract(epoch from now()),
  updated_at BIGINT NOT NULL DEFAULT extract(epoch from now()),
  ccla_enabled boolean,
  icla_enabled boolean,
  UNIQUE (foundation_id, cla_group_name),
  PRIMARY KEY(id)
);

create table cla.cla_group_project_managers (
  cla_group_id uuid NOT NULL REFERENCES cla.cla_groups(id), 
  project_manager_id uuid NOT NULL,
  PRIMARY KEY(cla_group_id,project_manager_id)
)
-- migrate:down
drop table cla.cla_group_project_managers;
drop table cla.cla_groups;

