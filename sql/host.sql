-- Table: public.Host

-- DROP TABLE IF EXISTS public."Host";

CREATE TABLE IF NOT EXISTS public.host
(
    name character varying(100) COLLATE pg_catalog."default",
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    description character varying(200) COLLATE pg_catalog."default",
    ip character varying(50) COLLATE pg_catalog."default",
    os character varying(100) COLLATE pg_catalog."default",
    status character varying(20) COLLATE pg_catalog."default",
    template character varying(100) COLLATE pg_catalog."default",
    tags character varying(100)[] COLLATE pg_catalog."default",
    num_cpu integer,
    mem_size integer,
    notes character varying(1000) COLLATE pg_catalog."default",
    CONSTRAINT "Host_pkey" PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.host
    OWNER to brianr;