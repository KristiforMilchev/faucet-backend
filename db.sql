BEGIN;


CREATE TABLE IF NOT EXISTS public.accounts
(
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    address text COLLATE pg_catalog."default" NOT NULL,
    last_lease timestamp without time zone NOT NULL,
    total_lease double precision,
    CONSTRAINT accounts_pkey PRIMARY KEY (id)
);
END;
