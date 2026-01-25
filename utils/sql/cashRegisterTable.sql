CREATE TABLE cash_register (
    id bigint NOT NULL,
    description character varying NOT NULL,
    customer character varying NOT NULL,
    specie character varying NOT NULL,
    input_value FLOAT,
    output_value FLOAT,
    total_balance FLOAT NOT NULL,
    date_created timestamp without time zone,
    CONSTRAINT cash_register_pkey PRIMARY KEY (id)
)