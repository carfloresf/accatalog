-- Database generated with pgModeler (PostgreSQL Database Modeler).
-- pgModeler  version: 0.9.1
-- PostgreSQL version: 10.0
-- Project Site: pgmodeler.io
-- Model Author: Carlos Flores

-- object: acadmin | type: ROLE --
-- DROP ROLE IF EXISTS acadmin;
CREATE ROLE acadmin WITH LOGIN;
-- ddl-end --



-- Database creation must be done outside a multicommand file.
-- These commands were put in this file only as a convenience.
-- -- object: accat | type: DATABASE --
-- -- DROP DATABASE IF EXISTS accat;
CREATE DATABASE accat
	ENCODING = 'UTF8'
-- 	LC_COLLATE = 'es_MX'
-- 	LC_CTYPE = 'es_MX'
	OWNER = acadmin;
-- -- ddl-end --
-- 

\c accat;

-- object: public.genre | type: TYPE --
-- DROP TYPE IF EXISTS public.genre CASCADE;
CREATE TYPE public.genre AS
 ENUM ('female','male','unisex');
-- ddl-end --
ALTER TYPE public.genre OWNER TO acadmin;
-- ddl-end --

-- object: public.measure | type: TYPE --
-- DROP TYPE IF EXISTS public.measure CASCADE;
CREATE TYPE public.measure AS
 ENUM ('metros','centimetros','pulgadas');
-- ddl-end --
ALTER TYPE public.measure OWNER TO acadmin;
-- ddl-end --

-- object: public.material | type: TABLE --
-- DROP TABLE IF EXISTS public.material CASCADE;
CREATE TABLE public.material(
	material_id bigserial NOT NULL,
	description varchar,
	cost float NOT NULL,
	measure public.measure NOT NULL,
	material_type_id smallint,
	brand_id bigint NOT NULL,
	"created_at" timestamp with time zone NOT NULL DEFAULT current_timestamp,
	"active" boolean NOT NULL DEFAULT 'true',
	CONSTRAINT material_pk PRIMARY KEY (material_id)

);
-- ddl-end --
ALTER TABLE public.material OWNER TO acadmin;
-- ddl-end --

CREATE TABLE public.inventory(
	inventory_id bigserial NOT NULL,
	material_id bigserial NOT NULL,
	quantity int NOT NULL,
	minimum int NOT NULL,
	maximum int NOT NULL,
	"created_at" timestamp with time zone NOT NULL DEFAULT current_timestamp,
	CONSTRAINT inventory_pk PRIMARY KEY (inventory_id)
);
-- ddl-end --
ALTER TABLE public.inventory OWNER TO acadmin;
-- ddl-end --
ALTER TABLE public.inventory ADD CONSTRAINT material_fk FOREIGN KEY (material_id)
REFERENCES public.material (material_id) MATCH FULL
ON DELETE RESTRICT ON UPDATE CASCADE;


-- object: public.material_type | type: TABLE --
-- DROP TABLE IF EXISTS public.material_type CASCADE;
CREATE TABLE public.material_type(
	material_type_id smallserial NOT NULL,
	name varchar NOT NULL,
	CONSTRAINT material_type_pk PRIMARY KEY (material_type_id)

);
-- ddl-end --
ALTER TABLE public.material_type OWNER TO acadmin;
-- ddl-end --

-- object: public.costume_material_relation | type: TABLE --
-- DROP TABLE IF EXISTS public.costume_material_relation CASCADE;
CREATE TABLE public.costume_material_relation(
	costume_id bigint NOT NULL,
	material_id bigint NOT NULL,
	quantity float,
	"created_at" timestamp with time zone NOT NULL DEFAULT current_timestamp
);
-- ddl-end --
ALTER TABLE public.costume_material_relation OWNER TO acadmin;
alter table costume_material_relation
	add constraint costume_material_relation_pk
		unique (costume_id, material_id);
-- ddl-end --

-- object: public.costume | type: TABLE --
-- DROP TABLE IF EXISTS public.costume CASCADE;
CREATE TABLE public.costume(
	costume_id bigserial NOT NULL,
	name varchar NOT NULL,
	color varchar,
	costume_code varchar,
	genre public.genre,
	costume_category_id smallint NOT NULL,
	"created_at" timestamp with time zone NOT NULL DEFAULT current_timestamp,
	"active" boolean NOT NULL DEFAULT 'true',
	CONSTRAINT costume_pk PRIMARY KEY (costume_id)
);
-- ddl-end --
ALTER TABLE public.costume OWNER TO acadmin;
-- ddl-end --

-- object: material_fk | type: CONSTRAINT --
-- ALTER TABLE public.costume_material_relation DROP CONSTRAINT IF EXISTS material_fk CASCADE;
ALTER TABLE public.costume_material_relation ADD CONSTRAINT material_fk FOREIGN KEY (material_id)
REFERENCES public.material (material_id) MATCH FULL
ON DELETE RESTRICT ON UPDATE CASCADE;
-- ddl-end --

-- object: material_type_fk | type: CONSTRAINT --
-- ALTER TABLE public.material DROP CONSTRAINT IF EXISTS material_type_fk CASCADE;
ALTER TABLE public.material ADD CONSTRAINT material_type_fk FOREIGN KEY (material_type_id)
REFERENCES public.material_type (material_type_id) MATCH FULL
ON DELETE SET NULL ON UPDATE CASCADE;
-- ddl-end --

-- object: public.costume_cost | type: TABLE --
-- DROP TABLE IF EXISTS public.costume_cost CASCADE;
CREATE TABLE public.costume_cost(
  costume_cost_id bigserial NOT NULL,
	costume_id bigint NOT NULL,
	actual_cost float,
	calculated_cost float,
	range_cost int,
	"created_at" timestamp with time zone NOT NULL DEFAULT current_timestamp,
	"active" boolean NOT NULL DEFAULT 'true',
	CONSTRAINT costume_cost_pk PRIMARY KEY (costume_cost_id)
);
-- ddl-end --
ALTER TABLE public.costume_cost OWNER TO acadmin;

-- object: costume_fk | type: CONSTRAINT --
-- ALTER TABLE public.costume_cost DROP CONSTRAINT IF EXISTS costume_cost_fk CASCADE;
ALTER TABLE public.costume_cost ADD CONSTRAINT costume_cost_fk FOREIGN KEY (costume_id)
REFERENCES public.costume (costume_id) MATCH FULL
ON DELETE RESTRICT ON UPDATE CASCADE;
-- ddl-end --

-- object: public.accesory | type: TABLE --
-- DROP TABLE IF EXISTS public.accesory CASCADE;
CREATE TABLE public.accesory(
	accesory_id bigserial NOT NULL,
	name varchar,
	accesory_code varchar,
	accesory_category_id smallint NOT NULL,
	"created_at" timestamp with time zone NOT NULL DEFAULT current_timestamp,
	"active" boolean NOT NULL DEFAULT 'true',
	CONSTRAINT accesory_pk PRIMARY KEY (accesory_id)

);
-- ddl-end --
ALTER TABLE public.accesory OWNER TO acadmin;
-- ddl-end --

-- object: public.accesory_material_relation | type: TABLE --
-- DROP TABLE IF EXISTS public.accesory_material_relation CASCADE;
CREATE TABLE public.accesory_material_relation(
	quantity smallint,
	material_id bigint NOT NULL,
	accesory_id bigint NOT NULL,
	"created_at" timestamp with time zone NOT NULL DEFAULT current_timestamp
);
-- ddl-end --
ALTER TABLE public.accesory_material_relation OWNER TO acadmin;
-- ddl-end --

-- object: material_fk | type: CONSTRAINT --
-- ALTER TABLE public.accesory_material_relation DROP CONSTRAINT IF EXISTS material_fk CASCADE;
ALTER TABLE public.accesory_material_relation ADD CONSTRAINT material_fk FOREIGN KEY (material_id)
REFERENCES public.material (material_id) MATCH FULL
ON DELETE RESTRICT ON UPDATE CASCADE;
-- ddl-end --

-- object: accesory_fk | type: CONSTRAINT --
-- ALTER TABLE public.accesory_material_relation DROP CONSTRAINT IF EXISTS accesory_fk CASCADE;
ALTER TABLE public.accesory_material_relation ADD CONSTRAINT accesory_fk FOREIGN KEY (accesory_id)
REFERENCES public.accesory (accesory_id) MATCH FULL
ON DELETE RESTRICT ON UPDATE CASCADE;
-- ddl-end --

-- object: public.costume_category | type: TABLE --
-- DROP TABLE IF EXISTS public.costume_category CASCADE;
CREATE TABLE public.costume_category(
	costume_category_id smallserial NOT NULL,
	costume_category_name varchar NOT NULL,
	"created_at" timestamp with time zone NOT NULL DEFAULT current_timestamp,
	CONSTRAINT costume_category_pk PRIMARY KEY (costume_category_id)

);
-- ddl-end --
ALTER TABLE public.costume_category OWNER TO acadmin;
-- ddl-end --

-- object: costume_category_fk | type: CONSTRAINT --
-- ALTER TABLE public.costume DROP CONSTRAINT IF EXISTS costume_category_fk CASCADE;
ALTER TABLE public.costume ADD CONSTRAINT costume_category_fk FOREIGN KEY (costume_category_id)
REFERENCES public.costume_category (costume_category_id) MATCH FULL
ON DELETE RESTRICT ON UPDATE CASCADE;
-- ddl-end --

-- object: public.accesory_category | type: TABLE --
-- DROP TABLE IF EXISTS public.accesory_category CASCADE;
CREATE TABLE public.accesory_category(
	accesory_category_id smallserial NOT NULL,
	accesory_category_name varchar NOT NULL,
	CONSTRAINT accesory_category_pk PRIMARY KEY (accesory_category_id),
	"created_at" timestamp with time zone NOT NULL DEFAULT current_timestamp
);
-- ddl-end --
ALTER TABLE public.accesory_category OWNER TO acadmin;
-- ddl-end --

-- object: accesory_category_fk | type: CONSTRAINT --
-- ALTER TABLE public.accesory DROP CONSTRAINT IF EXISTS accesory_category_fk CASCADE;
ALTER TABLE public.accesory ADD CONSTRAINT accesory_category_fk FOREIGN KEY (accesory_category_id)
REFERENCES public.accesory_category (accesory_category_id) MATCH FULL
ON DELETE RESTRICT ON UPDATE CASCADE;
-- ddl-end --

-- object: public.brand | type: TABLE --
-- DROP TABLE IF EXISTS public.brand CASCADE;
CREATE TABLE public.brand(
	brand_id bigserial NOT NULL,
	name varchar,
	description varchar,
	"created_at" timestamp with time zone NOT NULL DEFAULT current_timestamp,
	CONSTRAINT brand_pk PRIMARY KEY (brand_id)

);
-- ddl-end --
ALTER TABLE public.brand OWNER TO acadmin;
-- ddl-end --

-- object: brand_fk | type: CONSTRAINT --
-- ALTER TABLE public.material DROP CONSTRAINT IF EXISTS brand_fk CASCADE;
ALTER TABLE public.material ADD CONSTRAINT brand_fk FOREIGN KEY (brand_id)
REFERENCES public.brand (brand_id) MATCH FULL
ON DELETE RESTRICT ON UPDATE CASCADE;
-- ddl-end --

-- object: costume_fk | type: CONSTRAINT --
-- ALTER TABLE public.costume_material_relation DROP CONSTRAINT IF EXISTS costume_fk CASCADE;
ALTER TABLE public.costume_material_relation ADD CONSTRAINT costume_fk FOREIGN KEY (costume_id)
REFERENCES public.costume (costume_id) MATCH FULL
ON DELETE RESTRICT ON UPDATE CASCADE;
-- ddl-end --


CREATE TABLE public.api_key(
	"user" varchar not null,
	api_key varchar not null,
	permission varchar,
	"created_at" timestamp with time zone NOT NULL DEFAULT current_timestamp
);

alter table api_key owner to acadmin;


CREATE TABLE public.audit_log(
	audit_id bigserial NOT NULL,
	group_id bigserial,
	"table" varchar not null,
	"column" varchar not null,
	"table_id" varchar NOT NULL,
	"change_date" timestamp with time zone NOT NULL DEFAULT current_timestamp,
	"old_value" varchar,
	"new_value" varchar,
	"action" varchar,
	CONSTRAINT audit_log_pk PRIMARY KEY (audit_id)
);

alter table audit_log owner to acadmin;

-- DATOS BÁSICOS
INSERT INTO "brand" ("brand_id", "name", "description") VALUES (DEFAULT, 'Generica', 'Comprada en todos lados');
INSERT INTO "material_type" ("material_type_id", "name") VALUES (DEFAULT, 'Tela');
INSERT INTO "costume_category" ("costume_category_id", "costume_category_name") VALUES (DEFAULT, '1980s');
INSERT INTO "api_key" ("user", "api_key","permission") VALUES ('god', '123456','a');