-- migrate:up

create table cla.cla_templates (
  id uuid NOT NULL DEFAULT gen_random_uuid(),
  name varchar(255) NOT NULL,
  description text NOT NULL,
  created_at BIGINT NOT NULL DEFAULT extract(epoch from now()),
  updated_at BIGINT NOT NULL DEFAULT extract(epoch from now()),
  version integer NOT NULL,
  icla_html_body bytea,
  ccla_html_body bytea,
  meta_fields jsonb,
  icla_fields jsonb,
  ccla_fields jsonb,
  PRIMARY KEY(id)
);

-- migrate:down
drop table cla.cla_templates;
