SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: cla; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA cla;


--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


--
-- Name: pgcrypto; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pgcrypto WITH SCHEMA public;


--
-- Name: EXTENSION pgcrypto; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION pgcrypto IS 'cryptographic functions';


SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: cla_group_project_managers; Type: TABLE; Schema: cla; Owner: -
--

CREATE TABLE cla.cla_group_project_managers (
    cla_group_id uuid NOT NULL,
    project_manager_id uuid NOT NULL
);


--
-- Name: cla_groups; Type: TABLE; Schema: cla; Owner: -
--

CREATE TABLE cla.cla_groups (
    id uuid DEFAULT public.gen_random_uuid() NOT NULL,
    cla_group_name character varying(255) NOT NULL,
    project_id character varying(255) NOT NULL,
    created_at bigint DEFAULT date_part('epoch'::text, now()) NOT NULL,
    updated_at bigint DEFAULT date_part('epoch'::text, now()) NOT NULL,
    ccla_enabled boolean,
    icla_enabled boolean
);


--
-- Name: cla_templates; Type: TABLE; Schema: cla; Owner: -
--

CREATE TABLE cla.cla_templates (
    id uuid DEFAULT public.gen_random_uuid() NOT NULL,
    name character varying(255) NOT NULL,
    description text NOT NULL,
    created_at bigint DEFAULT date_part('epoch'::text, now()) NOT NULL,
    updated_at bigint DEFAULT date_part('epoch'::text, now()) NOT NULL,
    version integer NOT NULL,
    icla_html_body bytea,
    ccla_html_body bytea,
    meta_fields jsonb,
    icla_fields jsonb,
    ccla_fields jsonb
);


--
-- Name: events; Type: TABLE; Schema: cla; Owner: -
--

CREATE TABLE cla.events (
    id uuid DEFAULT public.gen_random_uuid() NOT NULL,
    event_type character varying(255) NOT NULL,
    user_id character varying(255) NOT NULL,
    project_id character varying(255) DEFAULT NULL::character varying,
    company_id character varying(255) DEFAULT NULL::character varying,
    event_time bigint DEFAULT date_part('epoch'::text, now()) NOT NULL,
    event_data jsonb
);


--
-- Name: repositories; Type: TABLE; Schema: cla; Owner: -
--

CREATE TABLE cla.repositories (
    id uuid DEFAULT public.gen_random_uuid() NOT NULL,
    repository_type character varying(255) NOT NULL,
    name character varying(255) NOT NULL,
    organization_name character varying(255) NOT NULL,
    url character varying(4096) NOT NULL,
    enabled boolean DEFAULT false NOT NULL,
    project_id character varying(255) NOT NULL,
    cla_group_id uuid NOT NULL,
    created_at bigint DEFAULT date_part('epoch'::text, now()) NOT NULL
);


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.schema_migrations (
    version character varying(255) NOT NULL
);


--
-- Name: cla_group_project_managers cla_group_project_managers_pkey; Type: CONSTRAINT; Schema: cla; Owner: -
--

ALTER TABLE ONLY cla.cla_group_project_managers
    ADD CONSTRAINT cla_group_project_managers_pkey PRIMARY KEY (cla_group_id, project_manager_id);


--
-- Name: cla_groups cla_groups_pkey; Type: CONSTRAINT; Schema: cla; Owner: -
--

ALTER TABLE ONLY cla.cla_groups
    ADD CONSTRAINT cla_groups_pkey PRIMARY KEY (id);


--
-- Name: cla_groups cla_groups_project_id_cla_group_name_key; Type: CONSTRAINT; Schema: cla; Owner: -
--

ALTER TABLE ONLY cla.cla_groups
    ADD CONSTRAINT cla_groups_project_id_cla_group_name_key UNIQUE (project_id, cla_group_name);


--
-- Name: cla_templates cla_templates_pkey; Type: CONSTRAINT; Schema: cla; Owner: -
--

ALTER TABLE ONLY cla.cla_templates
    ADD CONSTRAINT cla_templates_pkey PRIMARY KEY (id);


--
-- Name: events events_pkey; Type: CONSTRAINT; Schema: cla; Owner: -
--

ALTER TABLE ONLY cla.events
    ADD CONSTRAINT events_pkey PRIMARY KEY (id);


--
-- Name: repositories repositories_pkey; Type: CONSTRAINT; Schema: cla; Owner: -
--

ALTER TABLE ONLY cla.repositories
    ADD CONSTRAINT repositories_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: cla_group_project_managers cla_group_project_managers_cla_group_id_fkey; Type: FK CONSTRAINT; Schema: cla; Owner: -
--

ALTER TABLE ONLY cla.cla_group_project_managers
    ADD CONSTRAINT cla_group_project_managers_cla_group_id_fkey FOREIGN KEY (cla_group_id) REFERENCES cla.cla_groups(id);


--
-- Name: repositories repositories_cla_group_id_fkey; Type: FK CONSTRAINT; Schema: cla; Owner: -
--

ALTER TABLE ONLY cla.repositories
    ADD CONSTRAINT repositories_cla_group_id_fkey FOREIGN KEY (cla_group_id) REFERENCES cla.cla_groups(id);


--
-- PostgreSQL database dump complete
--


--
-- Dbmate schema migrations
--

INSERT INTO public.schema_migrations (version) VALUES
    ('20191112093719'),
    ('20191118083102'),
    ('20191209091738'),
    ('20191212125542');
